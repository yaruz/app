package account

import "github.com/pkg/errors"

type AccountSettings struct {
	LangID uint `json:"langId"`
}

func NewSettings() *AccountSettings {
	return &AccountSettings{}
}

func (e *AccountSettings) Validate() error {
	if e.LangID == 0 {
		return errors.New("AccountSettings validation error: LangID = 0")
	}
	return nil
}
