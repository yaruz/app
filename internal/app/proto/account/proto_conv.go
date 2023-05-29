package account

import (
	"github.com/yaruz/app/internal/domain/account"
)

func AccountSettingsProto2AccountSettings(accountSettingsProto *AccountSettings) (accountSettings *account.AccountSettings) {
	if accountSettingsProto == nil {
		return nil
	}

	accountSettings = &account.AccountSettings{
		LangID: uint(accountSettingsProto.LangID),
	}
	return accountSettings
}

func AccountSettings2AccountSettingsProto(accountSettings *account.AccountSettings) (accountSettingsProto *AccountSettings) {
	if accountSettings == nil {
		return nil
	}

	accountSettingsProto = &AccountSettings{
		LangID: uint64(accountSettings.LangID),
	}
	return accountSettingsProto
}
