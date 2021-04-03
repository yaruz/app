package auth

import (
	"net/http"

	routing "github.com/go-ozzo/ozzo-routing/v2"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"

	"github.com/minipkg/log"
	"github.com/minipkg/ozzo_routing/errorshandler"
)

type identity struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (i identity) Validate() error {
	return validation.ValidateStruct(&i,
		validation.Field(&i.Username, validation.Required, validation.Length(2, 100), is.Alphanumeric),
		validation.Field(&i.Password, validation.Required, validation.Length(4, 100)),
	)
}

// RegisterHandlers registers handlers for different HTTP requests.
//	POST /api/register - регистрация
//	POST /api/login - логин
func RegisterHandlers(rg *routing.RouteGroup, service Service, logger log.ILogger) {
	rg.Post("/login", login(service, logger))
	rg.Post("/register", register(service, logger))
}

func register(service Service, logger log.ILogger) routing.Handler {
	return func(c *routing.Context) error {
		var req identity

		if err := c.Read(&req); err != nil {
			logger.With(c.Request.Context()).Errorf("invalid request: %v", err)
			return errorshandler.BadRequest("")
		}

		if err := req.Validate(); err != nil {
			return err
		}
		ctx := c.Request.Context()
		token, err := service.Register(ctx, req.Username, req.Password)
		if err != nil {
			if er, ok := err.(errorshandler.Response); ok {
				logger.Errorf("Error while registering user. Status: %v; err: %q; details: %v", er.StatusCode(), er.Message, er.Details)
				return er
			}
			return err
		}
		return c.WriteWithStatus(struct {
			Token string `json:"token"`
		}{token}, http.StatusCreated)
	}
}

// login returns a handler that handles user login request.
func login(service Service, logger log.ILogger) routing.Handler {
	return func(c *routing.Context) error {
		var req identity

		if err := c.Read(&req); err != nil {
			logger.With(c.Request.Context()).Errorf("invalid request: %v", err)
			return errorshandler.BadRequest("")
		}

		if err := req.Validate(); err != nil {
			return err
		}

		token, err := service.Login(c.Request.Context(), req.Username, req.Password)
		if err != nil {
			return err
		}
		return c.Write(struct {
			Token string `json:"token"`
		}{token})
	}
}
