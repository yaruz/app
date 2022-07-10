package controller

import (
	routing "github.com/go-ozzo/ozzo-routing/v2"

	"github.com/minipkg/log"

	"github.com/yaruz/app/internal/domain/user"
	"github.com/yaruz/app/internal/pkg/auth"
	"github.com/yaruz/app/internal/pkg/socnets/tg"
)

type telegramController struct {
	RouteGroup *routing.RouteGroup
	Logger     log.ILogger
	Tg         tg.IService
	User       user.IService
	Auth       auth.Service
}

func NewTelegramController(r *routing.RouteGroup, logger log.ILogger, telegramService tg.IService, authService auth.Service, userService user.IService) *telegramController {
	return &telegramController{
		RouteGroup: r,
		Logger:     logger,
		Tg:         telegramService,
		User:       userService,
		Auth:       authService,
	}
}

// RegisterHandlers sets up the routing of the HTTP handlers.
func (c *telegramController) RegisterHandlers() {

	c.RouteGroup.Use(c.Auth.CheckAuthMiddleware)

	c.RouteGroup.Get(`/send-code`, c.sendCode)
	c.RouteGroup.Get(`/signin`, c.signin)

}

func (c *telegramController) sendCode(rctx *routing.Context) error {

	return rctx.Write(true)
}

func (c *telegramController) signin(rctx *routing.Context) error {
	return rctx.Write(true)
}
