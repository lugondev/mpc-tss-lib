package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	dbGateway "github.com/lugondev/mpc-tss-lib/db/gateway"
	dbgateway "github.com/lugondev/mpc-tss-lib/db/gateway/sqlc"
	"github.com/lugondev/mpc-tss-lib/internal/config"
	"github.com/lugondev/mpc-tss-lib/pb"
	rabbitmq "github.com/lugondev/mpc-tss-lib/pkg/ampq"
	grpcclient "github.com/lugondev/mpc-tss-lib/pkg/grpc/client"
	"github.com/lugondev/mpc-tss-lib/pkg/mpc/networking/server"
	zerolog "github.com/rs/zerolog/log"
	"github.com/thoas/go-funk"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

var logger = zerolog.Logger

const MessageResponse = "MPC Server"

func main() {
	flagConfigPath := flag.String("config", "configuration.gateway.yml", "config yml path file")
	flagAccessToken := flag.String("accessToken", "", "url to which a client will connect")
	flag.Parse()

	cfg, err := config.LoadConfig(flagConfigPath)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to LoadConfig")
	}
	sqlStore, err := config.NewGatewayDB(cfg.DB.Postgresql, &logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to Postgresql")
	}
	fmt.Println("=======================================================")
	fmt.Println("You are about to Gateway the MPC Server")
	fmt.Println("Server Port: ", cfg.Server.Port)
	for i, client := range cfg.Server.Clients {
		fmt.Println(fmt.Sprintf("Client %d: %s", i+1, client))
	}
	fmt.Println("=======================================================")

	createNetOperation(cfg, sqlStore, *flagAccessToken)
}

func createNetOperation(cfg *config.Config, dbStore *dbGateway.SQLStore, accessToken string) {
	netOperation := server.NewServer(nil, accessToken, &logger)

	e := echo.New()
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
	}))

	mqSign := rabbitmq.InitMqClient(&cfg.RabbitMQ, rabbitmq.ExchangeSign, rabbitmq.TopicSign)
	go mqSign.Subscribe(false)
	go ProcessEvent(mqSign, dbStore)

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

		dataRequisition := map[string]interface{}{
			"message": body.Message,
			"pubkey":  pubkey,
		}
		dataBytes, _ := json.Marshal(dataRequisition)
		requisitionCreated, err := dbStore.CreateRequisition(ctx, dbgateway.InsertRequisitionParams{
			Data:   dataBytes,
			Pubkey: pubkey,
		}, dbGateway.RequisitionSign)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("failed to create requisition: %s", err.Error()))
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
			if err = dbStore.FailRequisition(ctx, dbgateway.FailRequisitionParams{
				Reasons:     err.Error(),
				Requisition: requisitionCreated.Requisition,
			}); err != nil {
				logger.Error().Err(err).Msg("failed to FailRequisition")
			}

			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		dataRequisition["signatures"] = signatures
		dataBytes, _ = json.Marshal(dataRequisition)
		if err = dbStore.UpdateRequisition(ctx, dbgateway.UpdateRequisitionParams{
			Status:      "success",
			Data:        dataBytes,
			Requisition: requisitionCreated.Requisition,
		}); err != nil {
			logger.Error().Err(err).Msg("failed to UpdateRequisition")
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"message":       MessageResponse,
			"pubkey":        pubkey,
			"signatures":    signatures,
			"requisitionId": requisitionCreated.Requisition,
		})
	})

	e.POST("/keygen", func(c echo.Context) error {
		partyIDs := []string{
			uuid.New().String(),
			uuid.New().String(),
		}

		dataRequisition := map[string]interface{}{
			"parties": partyIDs,
		}
		dataBytes, _ := json.Marshal(dataRequisition)
		requisitionCreated, err := dbStore.CreateRequisition(ctx, dbgateway.InsertRequisitionParams{
			Data: dataBytes,
		}, dbGateway.RequisitionKeygen)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("failed to create requisition: %s", err.Error()))
		}

		logger.Info().Msgf("requisitionCreated: %v", requisitionCreated)

		go func() {
			logger.Info().Msgf("start operation: %v", partyIDs)
			if err := netOperation.StartOperation(ctx, partyIDs); err != nil {
				logger.Error().Err(err).Msg("failed to StartOperation")
			}
		}()

		go func(r dbgateway.Requisition, data map[string]interface{}, parties []string) {
			publicKeys, err := grpcclient.CallPartiesGRPCs(cfg.Server.Clients, parties, func(client pb.MpcPartyClient, parties []string, i int) (*pb.KeygenGeneratorResponse, error) {
				return client.KeygenGenerator(ctx, &pb.KeygenGeneratorParams{
					Id:  parties[i],
					Ids: parties,
				})
			})
			if err != nil {
				logger.Error().Err(err).Msgf("failed to KeygenGenerator: parties=%v", parties)
			}
			data["publicKeys"] = publicKeys
			dataBytes, _ = json.Marshal(data)
			if err = dbStore.UpdateRequisition(ctx, dbgateway.UpdateRequisitionParams{
				Status:      "success",
				Data:        dataBytes,
				Requisition: r.Requisition,
				Pubkey: funk.Map(publicKeys, func(pk *pb.KeygenGeneratorResponse) string {
					return pk.Pubkey
				}).([]string)[0],
			}); err != nil {
				logger.Error().Err(err).Msg("failed to FailRequisition")
			}
		}(requisitionCreated, dataRequisition, partyIDs)

		return c.JSON(http.StatusOK, map[string]interface{}{
			"message":       MessageResponse,
			"parties":       partyIDs,
			"requisitionId": requisitionCreated.Requisition,
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

	e.GET("/requisitions/:status", func(c echo.Context) error {
		rStatus := c.Param("status")
		if !funk.Contains([]string{"success", "fail", "pending"}, rStatus) {
			rStatus = ""
		}
		requisitions, err := dbStore.ListRequisitions(ctx, dbgateway.ListRequisitionsParams{
			Status: rStatus,
		})
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"message":      MessageResponse,
			"requisitions": requisitions,
		})
	})

	e.GET("/requisition/:id", func(c echo.Context) error {
		rId := c.Param("id")
		requisition, err := dbStore.GetRequisition(ctx, rId)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"message":     MessageResponse,
			"requisition": requisition,
		})
	})

	e.GET("/rabbitmq", func(c echo.Context) error {
		mqSign.Publish([]byte("{\"name\":\"Lugon\"}"), strconv.Itoa(time.Now().Second()))
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": MessageResponse,
		})
	})

	e.GET("/ws", func(c echo.Context) error {
		netOperation.ServeHTTP(c.Response(), c.Request())
		return nil
	})

	if err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.Server.Port), e); err != nil {
		logger.Fatal().Err(err).Msg("ListenAndServe")
	}
}

func ProcessEvent(mq *rabbitmq.BridgeMQ, dbStore *dbGateway.SQLStore) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGTERM)

	for {
		select {
		case msg := <-mq.MsgCh:
			if len(msg.Body) > 0 {
				logger.Info().Msgf("Received a message: %s", string(msg.Body))
				//if success {
				//	msg.Ack(false)
				//} else {
				//	tryTimesDemo++
				//	logger.Debug().Msgf("tryTimesDemo %d", tryTimesDemo)
				//	msg.Nack(false, tryTimesDemo < 5)
				//}
				retryTimes, err := dbStore.GetRetryTimes(context.Background(), string(msg.Body))
				if err != nil {
					logger.Error().Err(err).Msgf("failed to GetRetryTimes: %s", string(msg.Body))

					if err := msg.Ack(false); err != nil {
						logger.Error().Err(err).Msgf("failed to Ack: %v", err)
					}
				} else {
					logger.Debug().Msgf("retryTimes %d", retryTimes)
					if err := msg.Nack(false, retryTimes < 5); err != nil {
						logger.Error().Err(err).Msgf("failed to Nack: %v", err)
					}
				}
			}
		default:

		}
	}
}
