package main

import "os"
import "fmt"

import "github.com/joho/godotenv"
import "github.com/labstack/gommon/log"

import _ "github.com/jinzhu/gorm/dialects/mysql"

import "github.com/sizethree/miritos.api/db"
import "github.com/sizethree/miritos.api/net"
import "github.com/sizethree/miritos.api/routes"
import "github.com/sizethree/miritos.api/activity"
import "github.com/sizethree/miritos.api/middleware"

func main() {
	err := godotenv.Load()

	if err != nil {
		fmt.Printf("bad env: %s\n", err.Error())
		return
	}

	dbconf := db.Config{
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOSTNAME"),
		os.Getenv("DB_DATABASE"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_DEBUG") == "true",
	}

	port := os.Getenv("PORT")

	if len(port) < 1 {
		port = "8080"
	}

	database, err := db.Open(dbconf)

	if err != nil {
		panic(err)
	}

	// create the logger that will be shared by the server and the activity processor
	logger := log.New("miritos")
	logger.SetLevel(0)
	logger.SetHeader("[${time_rfc3339} ${level} ${short_file}]")

	// create the channel that will be used by the server runtime and activity processor
	stream := make(chan activity.Message, 100)

	// create our multiplexer and add our routes
	mux := net.Multiplexer{}

	mux.Use(middleware.InjectClient)
	mux.Use(middleware.InjectUser)

	mux.GET("/system", routes.System)

	mux.GET("/auth/user", routes.PrintAuth, middleware.RequireUser)
	mux.GET("/auth/roles", routes.PrintUserRoles, middleware.RequireUser)
	mux.GET("/auth/tokens", routes.PrintClientTokens, middleware.RequireClient)

	mux.GET("/activity", routes.FindActivity)

	// special route - returns all active activities based on their display schedules.
	mux.GET("/activity/live", routes.FindLiveActivity, middleware.RequireClient)

	mux.GET("/oauth/google/prompt", routes.GoogleOauthRedirect)
	mux.GET("/oauth/google/auth", routes.GoogleOauthReceiveCode)

	mux.GET("/user-roles", routes.FindRoles, middleware.RequireClient)

	mux.GET("/client-admins", routes.FindClientAdmins, middleware.RequireUser)

	mux.GET("/display-schedules", routes.FindDisplaySchedules, middleware.RequireClient)
	mux.PATCH("/display-schedules/:id", routes.UpdateDisplaySchedule, middleware.RequireUser)

	mux.GET("/users", routes.FindUser, middleware.RequireClient)
	mux.POST("/users", routes.CreateUser, middleware.RequireClient)
	mux.PATCH("/users/:id", routes.UpdateUser, middleware.RequireUser)

	mux.POST("/photos", routes.CreatePhoto, middleware.RequireClient)
	mux.GET("/photos", routes.FindPhotos, middleware.RequireClient)
	mux.GET("/photos/:id/view", routes.ViewPhoto, middleware.RequireClient)

	// create the server runtime and the activity processor runtime
	runtime := net.ServerRuntime{logger, dbconf, stream, &mux}
	processor := activity.Processor{logger, database, stream}
	server := net.Server{nil, &runtime}

	// start the server & processor
	server.Logger().Debugf(fmt.Sprintf("starting"))
	go processor.Begin()
	server.Run(fmt.Sprintf(":%s", port))
}
