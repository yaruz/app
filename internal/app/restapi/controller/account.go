// Copyright 2021 The casbin Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package controller

import (
	_ "embed"
	routing "github.com/go-ozzo/ozzo-routing/v2"
	"github.com/minipkg/log"
	"github.com/yaruz/app/internal/domain/user"

	//"github.com/casbin/casnode/object"
	//"github.com/casbin/casnode/util"
	"github.com/casdoor/casdoor-go-sdk/auth"
)

type accountController struct {
	Logger  log.ILogger
	Service user.IService
}

// RegisterHandlers sets up the routing of the HTTP handlers.
func RegisterAccountHandlers(r *routing.RouteGroup, service user.IService, logger log.ILogger, authHandler routing.Handler) {
	c := accountController{
		Logger:  logger,
		Service: service,
	}

	r.Get(`/user/<id:\d+>`, c.get)
	//r.Get("/users", c.list)

}

func init() {
	InitAuthConfig()
}

func InitAuthConfig() {
	casdoorEndpoint := "http://localhost:80"
	clientId := "clientId"
	clientSecret := "clientSecret"
	casdoorOrganization := "org"
	casdoorApplication := "app"

	auth.InitConfig(casdoorEndpoint, clientId, clientSecret, "JwtPublicKey", casdoorOrganization, casdoorApplication)
}

// @Title Signin
// @Description sign in as a member
// @Param   code     QueryString    string  true        "The code to sign in"
// @Param   state     QueryString    string  true        "The state"
// @Success 200 {object} controllers.api_controller.Response The Response object
// @router /signin [post]
// @Tag Account API
func (c *accountController) Signin(ctx *routing.Context) error {
	code := "code"
	state := "state"

	token, err := auth.GetOAuthToken(code, state)
	if err != nil {
		return err
	}

	claims, err := auth.ParseJwtToken(token.AccessToken)
	if err != nil {
		return err
	}

	//affected, err := object.UpdateMemberOnlineStatus(&claims.User, true, util.GetCurrentTime())
	//if err != nil {
	//	c.ResponseError(err.Error())
	//	return
	//}

	claims.AccessToken = token.AccessToken
	//c.SetSessionClaims(claims)

	return ctx.Write(claims)
}

// @Title Signout
// @Description sign out the current member
// @Success 200 {object} controllers.api_controller.Response The Response object
// @router /signout [post]
// @Tag Account API
func (c *accountController) Signout(ctx *routing.Context) error {
	//claims := c.GetSessionClaims()
	//if claims != nil {
	//_, err := object.UpdateMemberOnlineStatus(&claims.User, false, util.GetCurrentTime())
	//if err != nil {
	//	c.ResponseError(err.Error())
	//	return
	//}
	//}
	//
	//c.SetSessionClaims(nil)
	//
	//c.ResponseOk()
	return ctx.Write(true)
}

// @Title GetAccount
// @Description Get current account
// @Success 200 {object} controllers.api_controller.Response The Response object
// @router /get-account [get]
// @Tag Account API
//func (c *accountController) GetAccount() {
//	if c.RequireSignedIn() {
//		return
//	}
//
//	claims := c.GetSessionClaims()
//
//	c.ResponseOk(claims)
//}
//
//func (c *accountController) UpdateAccountBalance(amount int) {
//	user := c.GetSessionUser()
//	user.Score += amount
//	c.SetSessionUser(user)
//}
