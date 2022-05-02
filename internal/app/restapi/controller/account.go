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
	"github.com/minipkg/ozzo_routing/errorshandler"
	"github.com/yaruz/app/internal/domain/account"
	"github.com/yaruz/app/internal/domain/user"
	"github.com/yaruz/app/internal/pkg/auth"
	"net/http"
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

	c.RouteGroup.Use(c.Auth.CheckAuthMiddleware)

	c.RouteGroup.Put(`/account-settings`, c.accountSettingsUpdate)

	c.RouteGroup.Get(`/account/tg-signin`, c.tgSignin)
	//r.Get(`/user/<id:\d+>`, c.get)
	//r.Get("/users", c.list)

}

// todo: settings
// todo: все настройки + настройки по умолчанию

// @Title Signin
// @Description sign in as a member
// @Param   code     QueryString    string  true        "The code to sign in"
// @Param   state     QueryString    string  true        "The state"
// @Success 200 {object} controllers.api_controller.Response The Response object
// @router /signin [post]
// @Tag Account API
func (c *accountController) signin(rctx *routing.Context) error {
	var err error
	ctx := rctx.Request.Context()

	code := rctx.Request.URL.Query().Get("code")
	state := rctx.Request.URL.Query().Get("state")

	accountSettings, err := c.Auth.RoutingGetAccountSettingsWithDefaults(rctx)
	if err != nil {
		return routing.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if ctx, err = c.Auth.SignIn(ctx, code, state, accountSettings); err != nil {
		// todo: нужный статус
		return routing.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	sess := c.Auth.GetSession(ctx)

	return rctx.Write(struct {
		Token string `json:"token"`
	}{sess.JwtClaims.AccessToken})
}

func (c *accountController) tgSignin(ctx *routing.Context) error {

	return ctx.Write(true)
}

func (c *accountController) accountSettingsUpdate(rctx *routing.Context) (err error) {
	ctx := rctx.Request.Context()

	accountSettings := account.NewSettings()
	if err := rctx.Read(accountSettings); err != nil {
		c.Logger.With(ctx).Info(err)
		return errorshandler.BadRequest(err.Error())
	}

	if err := accountSettings.Validate(); err != nil {
		return errorshandler.BadRequest(err.Error())
	}

	if ctx, err = c.Auth.AccountSettingsUpdate(ctx, accountSettings); err != nil {
		return routing.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return rctx.Write(true)
}

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
