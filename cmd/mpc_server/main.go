package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/ambrosus/ambrosus-bridge/relay/cmd"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/config"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"net/http"
	"time"

	"github.com/ambrosus/ambrosus-bridge/relay/pkg/mpc/networking/server"
	zerolog "github.com/rs/zerolog/log"
)

var logger = zerolog.Logger

func main() {
	flagConfigPath := flag.String("config", "configuration.yml", "config yml path file")
	flagAccessToken := flag.String("accessToken", "", "url to which a client will connect")
	flag.Parse()

	fmt.Println("=======================================================")
	fmt.Println("You are about to Gateway the MPC Server")
	fmt.Println("=======================================================")

	cfg, _, err := config.LoadConfig(flagConfigPath)
	if err != nil {
		log.Fatal(err)
	}

	netOperation := createCenterNetworking(fmt.Sprintf(":%d", cfg.Server.WsPort), *flagAccessToken)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	e := echo.New()
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
	}))

	e.POST("/sign", func(c echo.Context) error {
		go func() {
			_, err := netOperation.Sign(ctx, []string{"B", "C"}, []byte(""))
			if err != nil {
				fmt.Println("error: ", err)
			}
		}()

		return c.JSON(http.StatusOK, map[string]string{
			"message": "Hello, World!",
		})
	})

	e.GET("/keygen", func(c echo.Context) error {
		go func() {
			err := netOperation.Keygen(ctx, []string{"B", "C"})
			if err != nil {
				fmt.Println("error: ", err)
			}
		}()

		return c.JSON(http.StatusOK, map[string]string{
			"message": "Hello, World!",
		})
	})

	if err := e.Start(fmt.Sprintf(":%d", cfg.Server.Port)); err != nil {
		log.Fatal(err)
	}
}

func createCenterNetworking(hostUrl string, accessToken string) cmd.NetworkingOperations {
	server_ := server.NewServer(nil, accessToken, &logger)

	go func() {
		err := http.ListenAndServe(hostUrl, server_)
		if err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
	}()
	return server_
}
