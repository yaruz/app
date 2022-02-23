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
	routing "github.com/go-ozzo/ozzo-routing/v2"
	"github.com/minipkg/log"
	"github.com/yaruz/app/internal/domain/user"
	"github.com/yaruz/app/internal/pkg/auth"
)

type accountController struct {
	RouteGroup *routing.RouteGroup
	Logger     log.ILogger
	User       user.IService
	Auth       auth.Service
}

func NewAccountController(r *routing.RouteGroup, userService user.IService, authService auth.Service, logger log.ILogger, authHandler routing.Handler) *accountController {
	return &accountController{
		RouteGroup: r,
		Logger:     logger,
		User:       userService,
		Auth:       authService,
	}
}

// RegisterHandlers sets up the routing of the HTTP handlers.
func (c *accountController) RegisterHandlers() {

	c.RouteGroup.Get(`/account/signin`, c.signin)
	c.RouteGroup.Get(`/account/fb-signin`, c.fbSignin)
	//r.Get(`/user/<id:\d+>`, c.get)
	//r.Get("/users", c.list)

}

// @Title Signin
// @Description sign in as a member
// @Param   code     QueryString    string  true        "The code to sign in"
// @Param   state     QueryString    string  true        "The state"
// @Success 200 {object} controllers.api_controller.Response The Response object
// @router /signin [post]
// @Tag Account API
func (c *accountController) signin(rctx *routing.Context) error {
	code := rctx.Request.URL.Query().Get("code")
	state := rctx.Request.URL.Query().Get("state")

	_, err := c.Auth.SignIn(rctx.Request.Context(), code, state, langId)

	return rctx.Write(err)
}

func (c *accountController) fbSignin(ctx *routing.Context) error {
	code := ctx.Request.URL.Query().Get("code")
	token := ctx.Request.URL.Query().Get("token")

	if code != "" {

	}

	return ctx.Write(true)
}

// todo: POST AccountSettings

// @Title Signout
// @Description sign out the current member
// @Success 200 {object} controllers.api_controller.Response The Response object
// @router /signout [post]
// @Tag Account API
func (c *accountController) signout(ctx *routing.Context) error {
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
//func (c *accountController) getAccount() {
//	if c.RequireSignedIn() {
//		return
//	}
//
//	claims := c.GetSessionClaims()
//
//	c.ResponseOk(claims)
//}
//
//func (c *accountController) updateAccountBalance(amount int) {
//	user := c.GetSessionUser()
//	user.Score += amount
//	c.SetSessionUser(user)
//}
