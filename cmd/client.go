package cmd

import (
	"fmt"
	"net/rpc"
	"strconv"

	"github.com/pawalt/kvstore/pkg/kv"
	"github.com/pawalt/kvstore/pkg/server"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(clientCmd)
	clientCmd.AddCommand(getCommand)
	clientCmd.AddCommand(putCommand)
}

var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "use client",
}

var getCommand = &cobra.Command{
	Use:   "get",
	Short: "get value",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			log.Fatal("expected an argument but did not get one")
		}

		client, err := rpc.DialHTTP(server.PROTOCOL, ":"+strconv.Itoa(server.PORT))
		if err != nil {
			log.Fatal(err)
		}

		path, err := kv.ParsePath(args[0])
		if err != nil {
			log.Fatal(err)
		}

		getReq := server.GetRequest{
			Path: path,
		}
		getResp := server.GetResponse{}
		err = client.Call("KVServer.GetRPC", &getReq, &getResp)
		if err != nil {
			log.Fatal(err)
		}

		if getResp.Value == nil {
			fmt.Println("no data at key")
		} else {
			fmt.Println("got data:")
			fmt.Println(string(getResp.Value))
		}
	},
}

var putCommand = &cobra.Command{
	Use:   "put",
	Short: "put value",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			log.Fatal("expected 2 arguments but were not provided")
		}

		client, err := rpc.DialHTTP(server.PROTOCOL, ":"+strconv.Itoa(server.PORT))
		if err != nil {
			log.Fatal(err)
		}

		path, err := kv.ParsePath(args[0])
		if err != nil {
			log.Fatal(err)
		}

		putReq := server.PutRequest{
			Path:  path,
			Value: []byte(args[1]),
		}
		putResp := server.PutResponse{}
		err = client.Call("KVServer.PutRPC", &putReq, &putResp)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("success")
	},
}
