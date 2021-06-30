package app

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/bhatneha/unpackker-api/api"
	// "github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

var (
	defaultconfig = `{
		"port":"8080",
		"logpath":"logdata.log",
		"ui":"false"
	}`
)

type ConfigData struct {
	AppPort   string    `json:"port"`
	UIPort    string    `json:"uiport"`
	LogPath   string    `json:"logpath"`
	UI        string    `json:"ui"`
	AppConfig string    `json:"appconfig"`
	Logger    io.Writer `json:"logger"`
	errgrp    errgroup.Group
}

func Run() {
	con, err := getConfig()
	if err != nil {
		fmt.Fprintf(os.Stdout, "%s \n", err.(error).Error())
		os.Exit(1)
	}

	if err := con.mergeConfig(); err != nil {
		fmt.Fprintf(os.Stdout, "%s \n", err.(error).Error())
		os.Exit(1)
	}

	if len(con.LogPath) != 0 {
		log, err := con.configLog()
		if err != nil {
			fmt.Fprintf(os.Stdout, "%s \n", err.(error).Error())
		}
		con.Logger = log

		con.runApp()

		// f,err := os.OpenFile(con.LogPath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
		// Formatter := new(log.TextFormatter)
		// Formatter.TimestampFormat = "02-01-2006 15:04:05"
		// Formatter.FullTimestamp = true
		// log.SetFormatter(Formatter)
		// if err != nil {
		// 	// Cannot open log file. Logging to stderr
		// 	fmt.Println(err)
		// } else {
		// 	log.SetOutput(f)
		// }
		// log.Debug("started packing!")
	}
}

func (c *ConfigData) runApp() {
	(c.errgrp).Go(func() error {
		err := (c.enableApp()).ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
		return err
	})

	// if c.UI == "true" {
	// 	(c.errgrp).Go(func() error {
	// 		err := (c.enableUI()).ListenAndServe()
	// 		if err != nil && err != http.ErrServerClosed {
	// 			log.Fatal(err)
	// 		}
	// 		return err
	// 	})
	// }
	if err := (c.errgrp).Wait(); err != nil {
		log.Fatal(err)
	}
}

func (c *ConfigData) enableApp() *http.Server {
	r := gin.New()
	if c.Logger != nil {
		gin.DefaultWriter = io.MultiWriter(c.Logger)
	}
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(CORS())

	r.MaxMultipartMemory = 30
	routes := &api.Router{Router: r}
	routes.AllRoutes()

	if len(c.AppPort) == 0 {
		c.AppPort = "80"
	}
	return &http.Server{
		Addr:    ":" + c.AppPort,
		Handler: r,
	}
}

// func (c *ConfigData) enableUI() *http.Server {
// 	ui := gin.New()
// 	if c.Logger != nil {
// 		gin.DefaultWriter = io.MultiWriter(c.Logger)
// 	}
// 	ui.Use(gin.Logger())
// 	ui.Use(gin.Recovery())

// 	routes := &api.Router{Router: ui}
// 	routes.AllRoutes()

// 	if len(c.UIPort) == 0 {
// 		c.UIPort = "8040"
// 	}
// 	return &http.Server{
// 		Addr:    ":" + c.UIPort,
// 		Handler: ui,
// 	}
// }

// func getCORS() cors.Config {
// 	corsConfig := cors.Config{
// 		AllowMethods:  []string{"POST", "GET"},
// 		AllowHeaders:  []string{"Origin", "Access-Control-Allow-Headers", "Access-Control-Allow-Origin", "Content-Length", "Content-Type"},
// 		ExposeHeaders: []string{"Content-Length"},
// 	}
// 	corsConfig.AllowAllOrigins = true
// 	// corsConfig.AllowCredentials = true
// 	// corsConfig.AddAllowHeaders("authorization")
// 	return corsConfig
// }

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers",
			"Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			// c.AbortWithStatus(204)
			c.Status(http.StatusNoContent)
			return
		}
		c.Next()
	}
}
