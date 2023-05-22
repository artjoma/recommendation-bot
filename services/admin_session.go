package services

type AdminSessions struct {
	sessions map[int64]*AdminSession // store admin state
}

func NewAdminSessions() *AdminSessions {
	sessions := make(map[int64]*AdminSession)
	return &AdminSessions{
		sessions: sessions,
	}
}

// resetSession remove all items from state machine
func (sessions *AdminSessions) prepareSession(accountId int64, command string) {
	sessions.sessions[accountId] = &AdminSession{
		wizardName:  command,
		wizardSteps: &WizardSteps{},
	}
}

// destroySession free mem resources
func (sessions *AdminSessions) destroySession(accountId int64) {
	delete(sessions.sessions, accountId)
}

func (sessions *AdminSessions) getAdminSession(accountId int64) *AdminSession {
	return sessions.sessions[accountId]
}

// wizardName string
type AdminSession struct {
	wizardName  string
	wizardSteps *WizardSteps
}

type WizardSteps []string

// IsEmpty: check if stack is empty
func (s *WizardSteps) IsEmpty() bool {
	return len(*s) == 0
}

func (s *WizardSteps) Len() int {
	return len(*s)
}

// Push a new value onto the stack
func (s *WizardSteps) Push(str string) {
	*s = append(*s, str) // Simply append the new value to the end of the stack
}

// Remove and return top element of stack. Return false if stack is empty.
func (s *WizardSteps) Pop() (string, bool) {
	if s.IsEmpty() {
		return "", false
	} else {
		index := len(*s) - 1
		element := (*s)[index]
		*s = (*s)[:index]
		return element, true
	}
}
