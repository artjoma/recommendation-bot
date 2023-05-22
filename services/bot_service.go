package services

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
	"strconv"
	"strings"
)

type TelegramBotService struct {
	appCtx *AppCtx
}

func NewTelegramBot(appCtx *AppCtx) *TelegramBotService {
	return &TelegramBotService{appCtx: appCtx}
}

/*
Check example: https://api.telegram.org/bot5721906598:AAGix2-ryTeaOsy4Jtbq5zJ3iCH1Zb9wXrA/getMe
*/
func (service *TelegramBotService) StartBot() {
	log.Info().Msg("Start bot service: " + service.appCtx.GetAppCfg().TokenId)
	log.Info().Msgf("Admins count: %d \n", len(service.appCtx.GetAppCfg().Admins))
	log.Info().Msgf("Templates count: %d \n", len(service.appCtx.templateService.templates.Templates))
	log.Info().Msgf("Triggers count: %d \n", len(service.appCtx.triggerService.triggers.Triggers))
	bot, err := tgbotapi.NewBotAPI(service.appCtx.GetAppCfg().TokenId)
	if err != nil {
		panic("Failed to create Telegram bot. Err:" + err.Error())
	}

	bot.Debug = false
	log.Info().Msg("Authorized on account: " + bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 600
	updates := bot.GetUpdatesChan(u)
	adminSessions := NewAdminSessions()

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.Chat.Type == "private" {
			// admin commands
			// check if hash admin role
			if !service.isUserAdmin(update.Message.From.ID) {
				continue
			}
			msgTxt := strings.TrimSpace(update.Message.Text)
			if len(msgTxt) == 0 {
				continue
			}
			accountId := update.Message.From.ID
			//log.Printf("[%d/%s] %s", accountId, update.Message.From.FirstName, update.Message.Text)
			cmd := update.Message.Command()
			log.Info().Msgf("[%d] %s", accountId, cmd)

			switch cmd {
			case AddTemplateCmd:
				adminSessions.prepareSession(accountId, AddTemplateCmd)
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Ok. Send me a new template name"))
			case RemoveTemplateCmd:
				adminSessions.prepareSession(accountId, RemoveTemplateCmd)
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Ok. Send me a template name"))
			case ListTemplatesCmd:
				templates := service.appCtx.templateService.templates.Templates
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Templates count: "+strconv.Itoa(len(templates))))
				for id, content := range templates {
					var buf strings.Builder
					fmt.Fprintf(&buf, "Template name: %s\n %s\n\n\n", id, content)
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, buf.String()))
				}
			case AddTriggerCmd:
				adminSessions.prepareSession(accountId, AddTriggerCmd)
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Ok. Send me a new trigger text"))
			case RemoveTriggerCmd:
				adminSessions.prepareSession(accountId, RemoveTriggerCmd)
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Ok. Send me a trigger id"))
			case ListTiggersCmd:
				triggers := service.appCtx.triggerService.triggers.Triggers
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Triggers count: "+strconv.Itoa(len(triggers))))
				var buf strings.Builder
				for id, trigger := range triggers {
					fmt.Fprintf(&buf, "| %20s | %s\n ", id, trigger)
				}
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, buf.String()))
			default:
				// is wizard active
				adminSession := adminSessions.getAdminSession(accountId)
				if adminSession == nil {
					// wizard not active
					continue
				}
				wizard := adminSession.wizardSteps
				switch adminSession.wizardName {
				case AddTemplateCmd:
					if wizard.IsEmpty() {
						// add template name(id)
						wizard.Push(msgTxt)
						bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Good. Send me template content"))
					} else {
						// get template id from prev. step
						templateId, _ := wizard.Pop()
						service.appCtx.templateService.AddTemplate(templateId, msgTxt)
						bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Success. Template with name: "+templateId+" was created"))
						adminSessions.destroySession(accountId)
					}
				case RemoveTemplateCmd:
					templateId := msgTxt
					service.appCtx.templateService.DeleteTemplate(templateId)
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Good. Template with name: "+templateId+" was removed"))
					adminSessions.destroySession(accountId)
				case AddTriggerCmd:
					if wizard.IsEmpty() {
						// add trigger text
						wizard.Push(msgTxt)
						bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Good. Send me a template name. "+
							"This trigger will be associated with this template"))
					} else {
						triggerText, _ := wizard.Pop()
						templateId := msgTxt
						triggerId := service.appCtx.triggerService.AddTrigger(templateId, triggerText)
						bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Success. Trigger with id: "+triggerId+" was created"))
						adminSessions.destroySession(accountId)
					}
				case RemoveTriggerCmd:
					triggerId := msgTxt
					service.appCtx.triggerService.DeleteTrigger(triggerId)
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Good. Trigger with id: "+triggerId+" was removed"))
					adminSessions.destroySession(accountId)
				}

			}
		} else {
			// user publish message to group
			if templateId := service.trySearch(update.Message.Text); templateId != "" {
				template := service.appCtx.templateService.templates.Templates[templateId]
				if template == "" {
					log.Error().Msgf("Template not found by id:%s" + templateId)
					continue
				}
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, template)
				msg.ReplyToMessageID = update.Message.MessageID
				bot.Send(msg)

			}
		}
	}
}

// TODO try to use full lexical search!
// trySearch return template id or empty string if not found
func (service *TelegramBotService) trySearch(msg string) string {
	normalizedMsg := strings.ToLower(msg)
	for triggerId, val := range service.appCtx.triggerService.triggers.Triggers {
		if strings.Contains(normalizedMsg, strings.ToLower(val)) {
			// for example TV:f5S
			return strings.Split(triggerId, ":")[0]
		}
	}

	return ""
}

func (service *TelegramBotService) isUserAdmin(userId int64) bool {
	for _, id := range service.appCtx.GetAppCfg().Admins {
		if id == userId {
			return true
		}
	}

	return false
}
