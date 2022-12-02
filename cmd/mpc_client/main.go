package main

import (
	"flag"
	"fmt"
	"github.com/lugondev/mpc-tss-lib/internal/config"
	grpc_client "github.com/lugondev/mpc-tss-lib/pkg/grpc/client"
	"log"
)

func main() {
	flagConfigPath := flag.String("config", "configuration.yml", "config yml path file")
	flagUrl := flag.String("url", "", "url to which a client will connect")
	flagAccessToken := flag.String("accessToken", "", "url to which a client will connect")
	flag.Parse()

	fmt.Println("=======================================================")
	fmt.Println("You are about to Client the MPC Server with access token: ", *flagAccessToken)
	fmt.Println("Gateway url: ", *flagUrl)
	fmt.Println("=======================================================")

	cfg, err := config.LoadConfig(flagConfigPath)
	if err != nil {
		log.Fatal(err)
	}
	sqlStore, err := config.NewDB(cfg.DB.Postgresql)
	if err != nil {
		log.Fatal(err)
	}

	grpc_client.StartGrpcClient(cfg.Grpc.Port, *flagUrl, sqlStore)
}

//
//func reshare(isServer bool, hostUrl, serverURL, accessToken string, id string, meInNewCommittee bool, partyIDsOld, partyIDsNew []string, thresholdOld, thresholdNew int, shareDir string) {
//	sharePath := getSharePath(shareDir, id)
//
//	fmt.Println("=======================================================")
//	fmt.Println("You are about to reshare the MPC share")
//	fmt.Println("Old IDS: ", partyIDsOld, "; threshold: ", thresholdOld)
//	fmt.Println("New IDS: ", partyIDsNew, "; threshold: ", thresholdNew)
//	fmt.Println("Your ID: ", id, "; share path: ", sharePath)
//	fmt.Println("Is your in new committee: ", meInNewCommittee)
//	fmt.Println("Is this server: ", isServer, "; server URL: ", serverURL)
//	fmt.Println("=======================================================")
//
//	mpcc := tss_wrap.NewMpc(id, thresholdOld, &logger)
//	if !meInNewCommittee {
//		readShare(mpcc, sharePath)
//	}
//
//	netOperation := createNetworking(isServer, hostUrl, serverURL, accessToken, mpcc)
//	ctx, _ := context.WithTimeout(context.Background(), 10*time.Minute)
//	err := netOperation.Reshare(ctx, partyIDsOld, partyIDsNew, thresholdNew)
//	if err != nil {
//		log.Fatal(err)
//	}
//	//saveShare(mpcc, sharePath)
//}

//func reshareBothCommittee(isServer bool, serverHost, url, accessToken, meID, meIDNew string, partyIDs, partyIDsNew []string, threshold, thresholdNew int, shareDir string) {
//	var wg sync.WaitGroup
//	if meIDNew != "" { // we are in new committee
//		wg.Add(1)
//		go func() {
//			reshare(isServer, serverHost, url, accessToken,
//				meIDNew, true, partyIDs, partyIDsNew,
//				threshold, thresholdNew, shareDir)
//			wg.Done()
//		}()
//
//		// if new committee just launched as server, old committee must be run as client, coz can't run server twice
//		time.Sleep(2 * time.Second)
//		isServer = false
//	}
//	if meID != "" { // we are in old committee
//		reshare(isServer, serverHost, url, accessToken,
//			meID, false, partyIDs, partyIDsNew,
//			threshold, thresholdNew, shareDir)
//	}
//	wg.Wait()
//}
