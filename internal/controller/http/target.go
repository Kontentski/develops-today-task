package httpcontroller

import (
	"github.com/Kontentski/develops-today-task/internal/service"
	"github.com/Kontentski/develops-today-task/pkg/errs"
	"github.com/gin-gonic/gin"
)

type targetRoutes struct {
	routerContext
}

func newTargetRoutes(options RouterOptions) {
	r := &targetRoutes{
		routerContext{
			services: options.Services,
			logger:   options.Logger.Named("targetRoutes"),
			cfg:      options.Config,
		},
	}

	// Standalone target operations
	p := options.Handler.Group("/targets")
	{
		p.PUT("/:id", errorHandler(options, r.updateTarget))
		p.DELETE("/:id", errorHandler(options, r.deleteTarget))
	}

	m := options.Handler.Group("/missions/:id")
	{
		m.POST("/targets", errorHandler(options, r.createTarget))
		m.GET("/targets", errorHandler(options, r.listTargets))
	}
}

type createTargetRequest struct {
	Name      string `json:"name" binding:"required"`
	Country   string `json:"country" binding:"required"`
	Notes     string `json:"notes"`
	Completed bool   `json:"completed"`
}

func (r *targetRoutes) createTarget(c *gin.Context) (interface{}, *httpErr) {
	missionID := c.Param("id")
	var req createTargetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, &httpErr{Type: httpErrTypeClient, Message: "invalid request body", Details: err}
	}

	opts := service.CreateTargetOptions{
		Name:      req.Name,
		Country:   req.Country,
		Notes:     req.Notes,
		Completed: req.Completed,
	}

	target, err := r.services.Target.CreateTarget(c, missionID, opts)
	if err != nil {
		if errs.IsExpected(err) {
			return nil, &httpErr{Type: httpErrTypeClient, Message: err.Error()}
		}
		return nil, &httpErr{Type: httpErrTypeServer, Message: "failed to create target", Details: err}
	}

	return target, nil
}

type updateTargetRequest struct {
	Notes     *string `json:"notes,omitempty"`
	Completed *bool   `json:"completed,omitempty"`
}

func (r *targetRoutes) updateTarget(c *gin.Context) (interface{}, *httpErr) {
	id := c.Param("id")
	var req updateTargetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, &httpErr{Type: httpErrTypeClient, Message: "invalid request body", Details: err}
	}

	opts := service.UpdateTargetOptions{
		Notes:     req.Notes,
		Completed: req.Completed,
	}

	target, err := r.services.Target.UpdateTarget(c, id, opts)
	if err != nil {
		if errs.IsExpected(err) {
			return nil, &httpErr{Type: httpErrTypeClient, Message: err.Error()}
		}
		return nil, &httpErr{Type: httpErrTypeServer, Message: "failed to update target", Details: err}
	}

	return target, nil
}

func (r *targetRoutes) deleteTarget(c *gin.Context) (interface{}, *httpErr) {
	id := c.Param("id")

	if err := r.services.Target.DeleteTarget(c, id); err != nil {
		if errs.IsExpected(err) {
			return nil, &httpErr{Type: httpErrTypeClient, Message: err.Error()}
		}
		return nil, &httpErr{Type: httpErrTypeServer, Message: "failed to delete target", Details: err}
	}

	return gin.H{"message": "target deleted successfully"}, nil
}

func (r *targetRoutes) listTargets(c *gin.Context) (interface{}, *httpErr) {
	missionID := c.Param("id")

	targets, err := r.services.Target.ListTargets(c, missionID)
	if err != nil {
		return nil, &httpErr{Type: httpErrTypeServer, Message: "failed to list targets", Details: err}
	}

	return targets, nil
} 