## Recommendation Telegram bot
<hr>
This bot provide recommendation message logic when user write specific phrase(trigger) at public chat. Bot try recognize a phrase for each message at chat from any user
and send to hime recommendation prepared message.

# Configuration and setup:
Create a `config.yaml` file at the same folder where application binary
```yaml 
# Bot token Id
tokenId: "<YOURS_TELEGRAM_BOT_TOKEN_ID_HERE>"

# Admins list
admins: [<TELEGRAM_ACCOUNT_ID_1>, <TELEGRAM_ACCOUNT_ID_2>, <TELEGRAM_ACCOUNT_ID_N>]
```
Create a `templates.json` file at the same folder where application binary
```json
{
  "modon": "2023-05-22T19:45:11.115555543Z",
  "templates": {
    "init": "test line\n- test1"
  }
}
```
Create a `triggers.json` file at the same folder where application binary
```json
{
  "modon": "2023-05-22T20:36:16.560233153Z",
  "triggers": {
    "init:WDU": "some best",
    "init:19b": "tell AAAA"
  }
}
```
## Features
Supported bot commands:
```go
const (
    AddTemplateCmd    string = "addtemplate"
    RemoveTemplateCmd        = "rmtemplate"
    ListTemplatesCmd         = "lstemplates" // Print all templates
    AddTriggerCmd            = "addtrigger"
    RemoveTriggerCmd         = "rmtrigger"
    ListTiggersCmd           = "lstriggers" // Print all triggers 
)

```

For stopping bot use `kill <PID>` or `CTRL+C` at console

After `/addXXX` or `/rmXXX` commands bot automatically flush changes to file: `template.json` or `trigger.json`

