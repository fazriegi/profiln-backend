package routes

import (
	"database/sql"
	"profiln-be/delivery/http"
	"profiln-be/delivery/http/middleware"
	"profiln-be/package/data"
	repository "profiln-be/package/data/repository"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func NewDataRoute(app *gin.RouterGroup, db *sql.DB, log *logrus.Logger) {
	repository := repository.NewDataRepository(db)
	usecase := data.NewDataUsecase(repository, log)
	controller := http.NewDataController(usecase)

	app.Use(middleware.Authentication())
	app.GET("/schools", controller.GetSchools)
	app.GET("/companies", controller.GetCompanies)
	app.GET("/issuing-organizations", controller.GetIssuingOrganizations)
	app.GET("/skills", controller.GetSkills)
	app.GET("/job-positions", controller.GetJobPositions)
}
