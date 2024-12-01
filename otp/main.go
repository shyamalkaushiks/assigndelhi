package main

import (
	"fmt"
	"os"

	"otppro/config"
	"otppro/logger"
	"otppro/model"
	"otppro/services"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/rs/zerolog"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func main() {
	var err error

	// log initiation
	log := zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, NoColor: false})

	// =========================== Setup logger =========================== //
	err = logger.SetupLogger(logger.Log)
	if err != nil {
		log.Error().Err(err).Msg("error while setup logger")
		return
	}
	// ==================================================================== //

	// =========================== Setup config =========================== //
	err = config.LoadConfig()
	if err != nil {
		logger.Log.Error().Err(err).Msg("error while setup config")
		return
	}
	// ==================================================================== //

	logger.Log.Info().Msg("Enter for main()...")

	//=========================== connection to db ===========================//
	// check database parameter validation
	if config.Config.DB_USER == "" || config.Config.DB_PASSWORD == "" || config.Config.DB_SERVER == "" || config.Config.DB_PORT == "" || config.Config.DB_DATABASE == "" {
		logger.Log.Error().Err(err).Msg("database env parameters are not found")
		return
	}

	var dsn string
	switch config.Config.DATABASE {
	case "postgres":
		logger.Log.Info().Msg("database connection : postgres")
		dsn = "host=" + config.Config.DB_SERVER + " port=" + config.Config.DB_PORT + " user=" + config.Config.DB_USER + " dbname=" + config.Config.DB_DATABASE + " password=" + config.Config.DB_PASSWORD + " sslmode=disable"
	case "mysql":
		logger.Log.Info().Msg("database connection : mysql")
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True", config.Config.DB_USER, config.Config.DB_PASSWORD, config.Config.DB_SERVER, config.Config.DB_PORT, config.Config.DB_DATABASE)
	default:
		logger.Log.Error().Err(err).Msg("Invalid database selection")
		log.Error().Err(err).Msg("Invalid database selection")
		return
	}

	log.Info().Msg(dsn)
	db, err := gorm.Open(config.Config.DATABASE, dsn)
	if err != nil {
		fmt.Println("Error in db :", err.Error())
		logger.Log.Error().Err(err).Msg("Error in database connection")
		return
	}
	defer db.Close()
	db.LogMode(true)
	db.SingularTable(true)
	db.DB().SetMaxIdleConns(10)

	logger.Log.Info().Msg("database connected successfully...")
	model.DBConn = db
	//=======================================================================//

	// create router in Gin
	router := gin.Default()
	// attaches a middleware to the router.
	router.Use(CORSMiddleware())

	// Setup Middleware for Database and Log
	router.Use(func(c *gin.Context) {
		c.Set("DB", db)
	})

	// Boostrap
	routes := &services.HandlerService{}
	routes.Bootstrap(router)

	port := config.Config.SERVICE_PORT
	if port == "" {
		logger.Log.Error().Err(err).Msg("service port not found")
		return
	}

	log.Info().Msg("Starting server on :" + port)
	logger.Log.Info().Msg("Starting server on :" + port)
	router.Run(":" + port)
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3006")
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Access-Control-Allow-Origin, Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Token")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
		} else {
			c.Next()
		}
	}
}
