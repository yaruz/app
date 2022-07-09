package controller

import (
	routing "github.com/go-ozzo/ozzo-routing/v2"

	"github.com/minipkg/log"

	"github.com/yaruz/app/internal/domain/user"
	"github.com/yaruz/app/internal/pkg/auth"
	"github.com/yaruz/app/internal/pkg/socnets/tgservice"
)

type telegramController struct {
	RouteGroup *routing.RouteGroup
	Logger     log.ILogger
	Tg         tgservice.IService
	User       user.IService
	Auth       auth.Service
}

func NewTelegramController(r *routing.RouteGroup, logger log.ILogger, tg tgservice.IService, authService auth.Service, userService user.IService) *telegramController {
	return &telegramController{
		RouteGroup: r,
		Logger:     logger,
		Tg:         tg,
		User:       userService,
		Auth:       authService,
	}
}

// RegisterHandlers sets up the routing of the HTTP handlers.
func (c *telegramController) RegisterHandlers() {

	c.RouteGroup.Use(c.Auth.CheckAuthMiddleware)

	c.RouteGroup.Get(`/account`, c.account)
	c.RouteGroup.Get(`/signin`, c.signin)

}

func (c *telegramController) account(rctx *routing.Context) error {
	return rctx.Write(true)
}

func (c *telegramController) signin(rctx *routing.Context) error {
	return rctx.Write(true)
}
