package token

import (
	"golang.org/x/oauth2"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TokenProto2Token(tokenProto *Token) (token *oauth2.Token, err error) {
	token = &oauth2.Token{
		AccessToken:  tokenProto.AccessToken,
		TokenType:    tokenProto.TokenType,
		RefreshToken: tokenProto.RefreshToken,
	}
	if tokenProto.Expiry != nil && tokenProto.Expiry.IsValid() {
		token.Expiry = tokenProto.Expiry.AsTime()
	}
	return token, nil
}

func Token2TokenProto(token *oauth2.Token) (tokenProto *Token, err error) {
	tokenProto = &Token{
		AccessToken:  token.AccessToken,
		TokenType:    token.TokenType,
		RefreshToken: token.RefreshToken,
		Expiry:       timestamppb.New(token.Expiry),
	}
	return tokenProto, nil
}
