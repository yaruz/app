package controller

import (
	"github.com/yaruz/app/internal/pkg/apperror"
	"github.com/minipkg/selection_condition"

	"github.com/minipkg/log"
	ozzo_routing "github.com/minipkg/ozzo_routing"
	"github.com/minipkg/ozzo_routing/errorshandler"

	"github.com/yaruz/app/internal/domain/mark"

	routing "github.com/go-ozzo/ozzo-routing/v2"
)

type markController struct {
	Logger  log.ILogger
	Service mark.IService
}

// RegisterHandlers sets up the routing of the HTTP handlers.
//	GET /api/marks/ - список всех моделей
//	GET /api/mark/{ID} - детали марки
func RegisterMarkHandlers(r *routing.RouteGroup, service mark.IService, logger log.ILogger, authHandler routing.Handler) {
	c := markController{
		Logger:  logger,
		Service: service,
	}

	r.Get("/marks", c.list)
	r.Get(`/mark/<id>`, c.get)
}

// get method is for getting a one entity by ID
func (c markController) get(ctx *routing.Context) error {
	id, err := ozzo_routing.ParseUintParam(ctx, "id")
	if err != nil {
		errorshandler.BadRequest("ID is required to be uint")
	}

	entity, err := c.Service.Get(ctx.Request.Context(), id)
	if err != nil {
		if err == apperror.ErrNotFound {
			c.Logger.With(ctx.Request.Context()).Info(err)
			return errorshandler.NotFound("not found")
		}
		c.Logger.With(ctx.Request.Context()).Error(err)
		return errorshandler.InternalServerError("")
	}

	return ctx.Write(entity)
}

// list method is for a getting a list of all entities
func (c markController) list(ctx *routing.Context) error {
	e := c.Service.NewEntity()
	e.TypeID = 2
	cond := &selection_condition.SelectionCondition{
		Where: e,
		SortOrder: []map[string]string{{
			"name": selection_condition.SortOrderAsc,
		}},
	}
	items, err := c.Service.Query(ctx.Request.Context(), cond)
	if err != nil {
		if err == apperror.ErrNotFound {
			c.Logger.With(ctx.Request.Context()).Info(err)
			return errorshandler.NotFound("")
		}
		c.Logger.With(ctx.Request.Context()).Error(err)
		return errorshandler.InternalServerError("")
	}
	return ctx.Write(items)
}
