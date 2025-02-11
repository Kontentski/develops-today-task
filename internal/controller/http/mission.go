package httpcontroller

import (
	"github.com/Kontentski/develops-today-task/internal/service"
	"github.com/Kontentski/develops-today-task/pkg/errs"
	"github.com/gin-gonic/gin"
)

type missionRoutes struct {
	routerContext
}

func newMissionRoutes(options RouterOptions) {
	r := &missionRoutes{
		routerContext{
			services: options.Services,
			logger:   options.Logger.Named("missionRoutes"),
			cfg:      options.Config,
		},
	}

	p := options.Handler.Group("/missions")
	{
		p.POST("/", errorHandler(options, r.createMission))
		p.DELETE("/:id", errorHandler(options, r.deleteMission))
		p.GET("/", errorHandler(options, r.listMissions))
		p.GET("/:id", errorHandler(options, r.getMission))
		p.PUT("/:id", errorHandler(options, r.updateMission))
		p.POST("/:id/assign", errorHandler(options, r.assignSpyCat))
	}
}

type createMissionRequest struct {
	Completed bool                     `json:"completed"`
	Targets   []createMissionTargetReq `json:"targets" binding:"required,min=1,max=3"`
}

type createMissionTargetReq struct {
	Name      string `json:"name" binding:"required"`
	Country   string `json:"country" binding:"required"`
	Notes     string `json:"notes"`
	Completed bool   `json:"completed"`
}

func (r *missionRoutes) createMission(c *gin.Context) (interface{}, *httpErr) {
	var req createMissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, &httpErr{Type: httpErrTypeClient, Message: "invalid request body", Details: err}
	}

	opts := service.CreateMissionOptions{
		Completed: req.Completed,
		Targets:   make([]service.CreateTargetOptions, len(req.Targets)),
	}

	for i, t := range req.Targets {
		opts.Targets[i] = service.CreateTargetOptions{
			Name:      t.Name,
			Country:   t.Country,
			Notes:     t.Notes,
			Completed: t.Completed,
		}
	}

	mission, err := r.services.Mission.CreateMission(c, opts)
	if err != nil {
		if errs.IsExpected(err) {
			return nil, &httpErr{Type: httpErrTypeClient, Message: err.Error()}
		}
		return nil, &httpErr{Type: httpErrTypeServer, Message: "failed to create mission", Details: err}
	}

	return mission, nil
}

func (r *missionRoutes) deleteMission(c *gin.Context) (interface{}, *httpErr) {
	id := c.Param("id")

	if err := r.services.Mission.DeleteMission(c, id); err != nil {
		if errs.IsExpected(err) {
			return nil, &httpErr{Type: httpErrTypeClient, Message: err.Error()}
		}
		return nil, &httpErr{Type: httpErrTypeServer, Message: "failed to delete mission", Details: err}
	}

	return gin.H{"message": "mission deleted successfully"}, nil
}

type updateMissionRequest struct {
	Completed bool `json:"completed" binding:"required"`
}

func (r *missionRoutes) updateMission(c *gin.Context) (interface{}, *httpErr) {
	id := c.Param("id")
	var req updateMissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, &httpErr{Type: httpErrTypeClient, Message: "invalid request body", Details: err}
	}

	mission, err := r.services.Mission.UpdateMission(c, id, service.UpdateMissionOptions{
		Completed: req.Completed,
	})
	if err != nil {
		if errs.IsExpected(err) {
			return nil, &httpErr{Type: httpErrTypeClient, Message: err.Error()}
		}
		return nil, &httpErr{Type: httpErrTypeServer, Message: "failed to update mission", Details: err}
	}

	return mission, nil
}

func (r *missionRoutes) getMission(c *gin.Context) (interface{}, *httpErr) {
	id := c.Param("id")

	mission, err := r.services.Mission.GetMission(c, id)
	if err != nil {
		if errs.IsExpected(err) {
			return nil, &httpErr{Type: httpErrTypeClient, Message: err.Error()}
		}
		return nil, &httpErr{Type: httpErrTypeServer, Message: "failed to get mission", Details: err}
	}

	return mission, nil
}

func (r *missionRoutes) listMissions(c *gin.Context) (interface{}, *httpErr) {
	missions, err := r.services.Mission.ListMissions(c)
	if err != nil {
		return nil, &httpErr{Type: httpErrTypeServer, Message: "failed to list missions", Details: err}
	}

	return missions, nil
}

type assignSpyCatRequest struct {
	SpyCatID string `json:"spyCatID" binding:"required"`
}

func (r *missionRoutes) assignSpyCat(c *gin.Context) (interface{}, *httpErr) {
	id := c.Param("id")
	var req assignSpyCatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, &httpErr{Type: httpErrTypeClient, Message: "invalid request body", Details: err}
	}

	if err := r.services.Mission.AssignSpyCat(c, id, req.SpyCatID); err != nil {
		if errs.IsExpected(err) {
			return nil, &httpErr{Type: httpErrTypeClient, Message: err.Error()}
		}
		return nil, &httpErr{Type: httpErrTypeServer, Message: "failed to assign spy cat", Details: err}
	}

	return gin.H{"message": "spy cat assigned successfully"}, nil
}
