package httpcontroller

import (
	"github.com/Kontentski/develops-today-task/internal/service"
	"github.com/Kontentski/develops-today-task/pkg/errs"
	"github.com/gin-gonic/gin"
)

type spyCatRoutes struct {
	routerContext
}

func newSpyCatRoutes(options RouterOptions) {
	r := &spyCatRoutes{
		routerContext{
			services: options.Services,
			logger:   options.Logger.Named("spyCatRoutes"),
			cfg:      options.Config,
		},
	}

	p := options.Handler.Group("/spycats")
	{
		p.POST("/", errorHandler(options, r.createSpyCat))
		p.DELETE("/:id", errorHandler(options, r.deleteSpyCat))
		p.GET("/", errorHandler(options, r.listSpyCats))
		p.GET("/:id", errorHandler(options, r.getSpyCat))
		p.PUT("/:id/salary", errorHandler(options, r.updateSpyCatSalary))
	}
}

type createSpyCatRequest struct {
	Name              string  `json:"name" binding:"required"`
	YearsOfExperience int     `json:"yearsOfExperience" binding:"required,gt=0"`
	Breed             string  `json:"breed" binding:"required"`
	Salary            float64 `json:"salary" binding:"required,gt=0"`
}

func (r *spyCatRoutes) createSpyCat(c *gin.Context) (interface{}, *httpErr) {
	var req createSpyCatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, &httpErr{Type: httpErrTypeClient, Message: "invalid request body", Details: err}
	}

	opts := service.CreateSpyCatOptions{
		Name:              req.Name,
		YearsOfExperience: req.YearsOfExperience,
		Breed:             req.Breed,
		Salary:            req.Salary,
	}

	cat, err := r.services.SpyCat.CreateSpyCat(c, opts)
	if err != nil {
		if errs.IsExpected(err) {
			return nil, &httpErr{Type: httpErrTypeClient, Message: err.Error()}
		}
		return nil, &httpErr{Type: httpErrTypeServer, Message: "failed to create spy cat", Details: err}
	}

	return cat, nil
}

func (r *spyCatRoutes) deleteSpyCat(c *gin.Context) (interface{}, *httpErr) {
	id := c.Param("id")

	if err := r.services.SpyCat.DeleteSpyCat(c, id); err != nil {
		if errs.IsExpected(err) {
			return nil, &httpErr{Type: httpErrTypeClient, Message: err.Error()}
		}
		return nil, &httpErr{Type: httpErrTypeServer, Message: "failed to delete spy cat", Details: err}
	}

	return gin.H{"message": "spy cat deleted successfully"}, nil
}

type updateSpyCatSalaryRequest struct {
	Salary float64 `json:"salary" binding:"required,gt=0"`
}

func (r *spyCatRoutes) updateSpyCatSalary(c *gin.Context) (interface{}, *httpErr) {
	id := c.Param("id")
	var req updateSpyCatSalaryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, &httpErr{Type: httpErrTypeClient, Message: "invalid request body", Details: err}
	}

	cat, err := r.services.SpyCat.UpdateSpyCatSalary(c, id, req.Salary)
	if err != nil {
		return nil, &httpErr{Type: httpErrTypeServer, Message: "failed to update salary", Details: err}
	}

	return cat, nil
}

func (r *spyCatRoutes) getSpyCat(c *gin.Context) (interface{}, *httpErr) {
	id := c.Param("id")

	cat, err := r.services.SpyCat.GetSpyCat(c, id)
	if err != nil {
		if errs.IsExpected(err) {
			return nil, &httpErr{Type: httpErrTypeClient, Message: err.Error()}
		}
		return nil, &httpErr{Type: httpErrTypeServer, Message: "failed to get spy cat", Details: err}
	}

	return cat, nil
}

func (r *spyCatRoutes) listSpyCats(c *gin.Context) (interface{}, *httpErr) {
	cats, err := r.services.SpyCat.ListSpyCats(c)
	if err != nil {
		return nil, &httpErr{Type: httpErrTypeServer, Message: "failed to list spy cats", Details: err}
	}

	return cats, nil
}
