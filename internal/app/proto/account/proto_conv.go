package account

import (
	"github.com/casdoor/casdoor-go-sdk/auth"
	"github.com/yaruz/app/internal/domain/account"
)

func AccountProto2Account(accountProto *Account) (account *auth.User, err error) {
	if accountProto == nil {
		return nil, nil
	}

	account = &auth.User{
		Owner:             accountProto.Owner,
		Name:              accountProto.Name,
		CreatedTime:       accountProto.CreatedTime,
		UpdatedTime:       accountProto.UpdatedTime,
		Id:                accountProto.ID,
		Type:              accountProto.Type,
		DisplayName:       accountProto.DisplayName,
		Avatar:            accountProto.Avatar,
		PermanentAvatar:   accountProto.PermanentAvatar,
		Email:             accountProto.Email,
		Phone:             accountProto.Phone,
		Location:          accountProto.Location,
		Address:           accountProto.Address,
		Affiliation:       accountProto.Affiliation,
		Title:             accountProto.Title,
		IdCardType:        accountProto.IdCardType,
		IdCard:            accountProto.IdCard,
		Homepage:          accountProto.Homepage,
		Bio:               accountProto.Bio,
		Tag:               accountProto.Tag,
		Region:            accountProto.Region,
		Language:          accountProto.Language,
		Gender:            accountProto.Gender,
		Birthday:          accountProto.Birthday,
		Education:         accountProto.Education,
		Score:             int(accountProto.Score),
		Ranking:           int(accountProto.Ranking),
		IsDefaultAvatar:   accountProto.IsDefaultAvatar,
		IsOnline:          accountProto.IsOnline,
		IsAdmin:           accountProto.IsAdmin,
		IsGlobalAdmin:     accountProto.IsGlobalAdmin,
		IsForbidden:       accountProto.IsForbidden,
		IsDeleted:         accountProto.IsDeleted,
		SignupApplication: accountProto.SignupApplication,
		Hash:              accountProto.Hash,
		PreHash:           accountProto.PreHash,
		CreatedIp:         accountProto.CreatedIp,
		LastSigninTime:    accountProto.LastSigninTime,
		LastSigninIp:      accountProto.LastSigninIp,
		Github:            accountProto.Github,
		Google:            accountProto.Google,
		QQ:                accountProto.QQ,
		WeChat:            accountProto.WeChat,
		Facebook:          accountProto.Facebook,
		DingTalk:          accountProto.DingTalk,
		Weibo:             accountProto.Weibo,
		Gitee:             accountProto.Gitee,
		LinkedIn:          accountProto.LinkedIn,
		Wecom:             accountProto.Wecom,
		Lark:              accountProto.Lark,
		Gitlab:            accountProto.Gitlab,
		Ldap:              accountProto.Ldap,
		Properties:        accountProto.Properties,
	}
	return account, nil
}

func Account2AccountProto(account *auth.User) (accountProto *Account, err error) {
	if account == nil {
		return nil, nil
	}

	accountProto = &Account{
		Owner:             account.Owner,
		Name:              account.Name,
		CreatedTime:       account.CreatedTime,
		UpdatedTime:       account.UpdatedTime,
		ID:                account.Id,
		Type:              account.Type,
		DisplayName:       account.DisplayName,
		FirstName:         "",
		LastName:          "",
		Avatar:            account.Avatar,
		PermanentAvatar:   account.PermanentAvatar,
		Email:             account.Email,
		EmailVerified:     false,
		Phone:             account.Phone,
		Location:          account.Location,
		Address:           account.Address,
		Affiliation:       account.Affiliation,
		Title:             account.Title,
		IdCardType:        account.IdCardType,
		IdCard:            account.IdCard,
		Homepage:          account.Homepage,
		Bio:               account.Bio,
		Tag:               account.Tag,
		Region:            account.Region,
		Language:          account.Language,
		Gender:            account.Gender,
		Birthday:          account.Birthday,
		Education:         account.Education,
		Score:             int64(account.Score),
		Karma:             int64(account.Karma),
		Ranking:           int64(account.Ranking),
		IsDefaultAvatar:   account.IsDefaultAvatar,
		IsOnline:          account.IsOnline,
		IsAdmin:           account.IsAdmin,
		IsGlobalAdmin:     account.IsGlobalAdmin,
		IsForbidden:       account.IsForbidden,
		IsDeleted:         account.IsDeleted,
		SignupApplication: account.SignupApplication,
		Hash:              account.Hash,
		PreHash:           account.PreHash,
		CreatedIp:         account.CreatedIp,
		LastSigninTime:    account.LastSigninTime,
		LastSigninIp:      account.LastSigninIp,
		Github:            account.Github,
		Google:            account.Google,
		QQ:                account.QQ,
		WeChat:            account.WeChat,
		WeChatUnionId:     "",
		Facebook:          account.Facebook,
		DingTalk:          account.DingTalk,
		Weibo:             account.Weibo,
		Gitee:             account.Gitee,
		LinkedIn:          account.LinkedIn,
		Wecom:             account.Wecom,
		Lark:              account.Lark,
		Gitlab:            account.Gitlab,
		Adfs:              "",
		Baidu:             "",
		Alipay:            "",
		Casdoor:           "",
		Infoflow:          "",
		Apple:             "",
		AzureAD:           "",
		Slack:             "",
		Steam:             "",
		Custom:            "",
		Ldap:              account.Ldap,
		Properties:        account.Properties,
	}
	return accountProto, nil
}

func AccountSettingsProto2AccountSettings(accountSettingsProto *AccountSettings) (accountSettings *account.AccountSettings, err error) {
	if accountSettingsProto == nil {
		return nil, nil
	}

	accountSettings = &account.AccountSettings{
		LangID: uint(accountSettingsProto.LangID),
	}
	return accountSettings, nil
}

func AccountSettings2AccountSettingsProto(accountSettings *account.AccountSettings) (accountSettingsProto *AccountSettings, err error) {
	if accountSettings == nil {
		return nil, nil
	}

	accountSettingsProto = &AccountSettings{
		LangID: uint64(accountSettings.LangID),
	}
	return accountSettingsProto, nil
}

func ClaimsProto2Claims(claimsProto *JwtClaims) (claims *auth.Claims, err error) {
	if claimsProto == nil {
		return nil, nil
	}

	account, err := AccountProto2Account(claimsProto.User)
	if err != nil {
		return nil, err
	}
	claims = &auth.Claims{
		User:        *account,
		AccessToken: claimsProto.AccessToken,
	}
	return claims, nil
}

func Claims2ClaimsProto(claims *auth.Claims) (claimsProto *JwtClaims, err error) {
	if claims == nil {
		return nil, nil
	}

	accountProto, err := Account2AccountProto(&claims.User)
	if err != nil {
		return nil, err
	}
	claimsProto = &JwtClaims{
		User:        accountProto,
		AccessToken: claims.AccessToken,
	}
	return claimsProto, nil
}
