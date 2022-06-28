package sn_account

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/yaruz/app/internal/domain/sn_account"
)

func SNAccountProto2SNAccount(sNAccountProto *SNAccount) (sNAccount *sn_account.SNAccount, err error) {
	sNAccount = &sn_account.SNAccount{
		ID:     uint(sNAccountProto.ID),
		TypeID: uint(sNAccountProto.TypeID),
		SNID:   sNAccountProto.SNID,
	}
	if sNAccountProto.CreatedAt != nil && sNAccountProto.CreatedAt.IsValid() {
		sNAccount.CreatedAt = sNAccountProto.CreatedAt.AsTime()
	}
	return sNAccount, nil
}

func SNAccount2SNAccountProto(sNAccount *sn_account.SNAccount) (sNAccountProto *SNAccount, err error) {
	sNAccountProto = &SNAccount{
		ID:        uint64(sNAccount.ID),
		TypeID:    uint64(sNAccount.TypeID),
		SNID:      sNAccount.SNID,
		CreatedAt: timestamppb.New(sNAccount.CreatedAt),
	}
	return sNAccountProto, nil
}
