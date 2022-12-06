package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/lugondev/mpc-tss-lib/internal/config"
	"github.com/lugondev/mpc-tss-lib/pb"
	grpcclient "github.com/lugondev/mpc-tss-lib/pkg/grpc/client"
	"github.com/lugondev/mpc-tss-lib/pkg/mpc/networking/server"
	zerolog "github.com/rs/zerolog/log"
	"github.com/thoas/go-funk"
	"net/http"
	"strconv"
)

var logger = zerolog.Logger

const MessageResponse = "MPC Server"

func main() {
	flagConfigPath := flag.String("config", "configuration.yml", "config yml path file")
	flagAccessToken := flag.String("accessToken", "", "url to which a client will connect")
	flag.Parse()

	cfg, err := config.LoadConfig(flagConfigPath)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to LoadConfig")
	}
	fmt.Println("=======================================================")
	fmt.Println("You are about to Gateway the MPC Server")
	fmt.Println("Server Port: ", cfg.Server.Port)
	for i, client := range cfg.Server.Clients {
		fmt.Println(fmt.Sprintf("Client %d: %s", i+1, client))
	}
	fmt.Println("=======================================================")

	createNetOperation(cfg.Server.Port, cfg, *flagAccessToken)
}

func createNetOperation(port int64, cfg *config.Config, accessToken string) {
	netOperation := server.NewServer(nil, accessToken, &logger)

	e := echo.New()
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
	}))

	//ctx, _ := context.WithTimeout(context.Background(), 1*time.Minute)
	ctx := context.Background()
	e.POST("/sign/:pubkey", func(c echo.Context) error {
		pubkey := c.Param("pubkey")
		var body struct{ Message string }
		if err := c.Bind(&body); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		if body.Message == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "message is empty")
		}
		msg := common.FromHex(body.Message)
		if len(msg) == 0 {
			return echo.NewHTTPError(http.StatusBadRequest, "message is not hex")
		}

		parties, err := grpcclient.CallClientGRPCs(cfg.Server.Clients, func(client pb.MpcPartyClient, i int) (*pb.GetPartyResponse, error) {
			return client.GetParty(ctx, &pb.GetPartyParams{Pubkey: pubkey})
		})
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		partyIDs := grpcclient.GetPartyIDs(parties)
		for _, partyID := range partyIDs {
			if netOperation.IsClientConnected(partyID) {
				return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("client %s is already connected", partyID))
			}
		}
		go func() {
			logger.Info().Msgf("start operation: %v", partyIDs)
			if err := netOperation.StartOperation(ctx, partyIDs); err != nil {
				logger.Error().Err(err).Msg("failed to StartOperation")
			}
		}()

		signatures, err := grpcclient.CallPartiesGRPCs(cfg.Server.Clients, partyIDs, func(client pb.MpcPartyClient, parties []string, i int) (*pb.SignResponse, error) {
			return client.Sign(ctx, &pb.SignParams{
				Id:      parties[i],
				Parties: parties,
				Message: msg,
				Pubkey:  pubkey,
			})
		})
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"message":    MessageResponse,
			"signatures": signatures,
		})
	})

	e.POST("/keygen", func(c echo.Context) error {
		partyIDs := []string{
			uuid.New().String(),
			uuid.New().String(),
		}

		go func() {
			logger.Info().Msgf("start operation: %v", partyIDs)
			if err := netOperation.StartOperation(ctx, partyIDs); err != nil {
				logger.Error().Err(err).Msg("failed to StartOperation")
			}
		}()

		publicKeys, err := grpcclient.CallPartiesGRPCs(cfg.Server.Clients, partyIDs, func(client pb.MpcPartyClient, parties []string, i int) (*pb.KeygenGeneratorResponse, error) {
			return client.KeygenGenerator(ctx, &pb.KeygenGeneratorParams{
				Id:  parties[i],
				Ids: parties,
			})
		})
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message":    MessageResponse,
			"parties":    partyIDs,
			"publicKeys": publicKeys,
		})
	})

	e.GET("/get-party/:pubkey", func(c echo.Context) error {
		pubkey := c.Param("pubkey")
		parties, err := grpcclient.CallClientGRPCs(cfg.Server.Clients, func(client pb.MpcPartyClient, i int) (*pb.GetPartyResponse, error) {
			return client.GetParty(ctx, &pb.GetPartyParams{Pubkey: pubkey})
		})
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": MessageResponse,
			"parties": funk.Filter(parties, func(p *pb.GetPartyResponse) bool {
				return p != nil
			}),
		})
	})

	e.GET("/get-parties/:clientIndex", func(c echo.Context) error {
		clientIndexParam := c.Param("clientIndex")
		clientIndex, err := strconv.Atoi(clientIndexParam)
		if clientIndex-1 > len(cfg.Server.Clients) {
			return echo.NewHTTPError(http.StatusBadRequest, "clientIndex is invalid")
		}
		parties, err := grpcclient.CallClientGRPC(cfg.Server.Clients[clientIndex], func(client pb.MpcPartyClient) (*pb.GetPartiesResponse, error) {
			return client.GetParties(ctx, &pb.GetPartiesParams{})
		})
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": MessageResponse,
			"parties": parties,
		})
	})

	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": MessageResponse,
		})
	})

	e.GET("/ws", func(c echo.Context) error {
		netOperation.ServeHTTP(c.Response(), c.Request())
		return nil
	})

	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), e); err != nil {
		logger.Fatal().Err(err).Msg("ListenAndServe")
	}
}
