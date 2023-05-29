package tg

const (
	defaultTgSessionStringLength = 1000
)

type Session struct {
	IsAuthorized  bool
	Session       []byte
	ID            string
	PhoneCodeHash string
	Phone         string
}

func NewSession(ID string) *Session {
	return &Session{
		ID:      ID,
		Session: make([]byte, 0, defaultTgSessionStringLength),
	}
}
