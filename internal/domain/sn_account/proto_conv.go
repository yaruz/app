package sn_account

import (
	"github.com/yaruz/app/internal/pkg/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func SNAccountProto2SNAccount(sNAccountProto *proto.SNAccount) (sNAccount *SNAccount, err error) {
	sNAccount = &SNAccount{
		ID:     uint(sNAccountProto.ID),
		TypeID: uint(sNAccountProto.TypeID),
		SNID:   sNAccountProto.SNID,
	}
	if sNAccountProto.CreatedAt != nil && sNAccountProto.CreatedAt.IsValid() {
		sNAccount.CreatedAt = sNAccountProto.CreatedAt.AsTime()
	}
	return sNAccount, nil
}

func SNAccount2SNAccountProto(sNAccount *SNAccount) (sNAccountProto *proto.SNAccount, err error) {
	sNAccountProto = &proto.SNAccount{
		ID:        uint64(sNAccount.ID),
		TypeID:    uint64(sNAccount.TypeID),
		SNID:      sNAccount.SNID,
		CreatedAt: timestamppb.New(sNAccount.CreatedAt),
	}
	return sNAccountProto, nil
}
