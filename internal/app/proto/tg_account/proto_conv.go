package tg_account

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/yaruz/app/internal/domain/tg_account"
)

func TgAccountProto2TgAccount(tgAccountProto *TgAccount) (tgAccount *tg_account.TgAccount, err error) {
	tgAccount = &tg_account.TgAccount{
		ID:   uint(tgAccountProto.ID),
		TgID: tgAccountProto.TgID,
	}
	if tgAccountProto.CreatedAt != nil && tgAccountProto.CreatedAt.IsValid() {
		tgAccount.CreatedAt = tgAccountProto.CreatedAt.AsTime()
	}
	return tgAccount, nil
}

func TgAccount2TgAccountProto(tgAccount *tg_account.TgAccount) (tgAccountProto *TgAccount, err error) {
	tgAccountProto = &TgAccount{
		ID:        uint64(tgAccount.ID),
		TgID:      tgAccount.TgID,
		CreatedAt: timestamppb.New(tgAccount.CreatedAt),
	}
	return tgAccountProto, nil
}
