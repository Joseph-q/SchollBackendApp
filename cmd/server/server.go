package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	assistanceController "github.com/juseph-q/SchoolPr/internal/assistance/controller"
	assistanceServices "github.com/juseph-q/SchoolPr/internal/assistance/services"
	"github.com/juseph-q/SchoolPr/internal/config"
	coursesController "github.com/juseph-q/SchoolPr/internal/courses/controller"
	coursesService "github.com/juseph-q/SchoolPr/internal/courses/service"
	studentController "github.com/juseph-q/SchoolPr/internal/student/controller"
	studentService "github.com/juseph-q/SchoolPr/internal/student/services"

	"github.com/juseph-q/SchoolPr/internal/database"

	"github.com/juseph-q/SchoolPr/internal/validations"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

var serverCmd = &cobra.Command{
	Use: "server",
	Run: func(cmd *cobra.Command, args []string) {
		runApplication()
	},
}

func runApplication() {
	conf, _ := config.Load("config/config_develop.yaml")

	app := fx.New(
		fx.Supply(conf),
		fx.Provide(
			//setup database
			database.NewDataBase,
			//student
			studentService.NewStudentService,
			studentController.NewStudentController,
			//courses
			coursesService.NewCourseService,
			coursesController.NewCourseController,
			//assistances
			assistanceServices.NewAssitanceService,
			assistanceController.NewAssistanceController,
			//server
			serverOn,
		),
		fx.StopTimeout(conf.Server.GracefulTimeout+time.Second),
		fx.Invoke(
			assistanceController.HandleRoutes,
			coursesController.HandleRoutes,
			studentController.HandleRoutes,
			func(*gin.Engine) {},
		),
	)

	app.Run()
}

func serverOn(lc fx.Lifecycle, conf *config.Config) *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     conf.Cors.AllowedOrigins,
		AllowMethods:     conf.Cors.AllowedMethods,
		AllowHeaders:     conf.Cors.AllowedHeaders,
		AllowCredentials: conf.Cors.AllowCredentials,
		MaxAge:           time.Duration(conf.Cors.MaxAge),
	}))

	//Custom Validations
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("gender", validations.ValidateGender)
		v.RegisterValidation("dateformat", validations.DateValidation)

	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", conf.Server.Port),
		Handler:      r,
		ReadTimeout:  conf.Server.MaxReadTimeout,
		WriteTimeout: conf.Server.MaxWriteTimeout,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					fmt.Println("failed to close http server", "err", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {

			return srv.Shutdown(ctx)
		},
	})

	return r
}
