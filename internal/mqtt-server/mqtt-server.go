package mqtt_server

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/hooks/auth"
	"github.com/mochi-mqtt/server/v2/listeners"
)

func New() {
	tcpAddr := "192.168.0.192:1883"
	//wsAddr := flag.String("ws", ":1882", "network address for Websocket listener")
	infoAddr := "192.168.0.192:8085"
	flag.Parse()

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		done <- true
	}()

	server := mqtt.New(nil)
	_ = server.AddHook(new(auth.AllowHook), nil)

	tcp := listeners.NewTCP(listeners.Config{
		ID:      "t1",
		Address: tcpAddr,
	})
	err := server.AddListener(tcp)
	if err != nil {
		log.Fatal(err)
	}

	//ws := listeners.NewWebsocket(listeners.Config{
	//	ID:      "ws1",
	//	Address: *wsAddr,
	//})
	//err = server.AddListener(ws)
	//if err != nil {
	//	log.Fatal(err)
	//}

	stats := listeners.NewHTTPStats(
		listeners.Config{
			ID:      "info",
			Address: infoAddr,
		},
		server.Info,
	)
	err = server.AddListener(stats)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		err := server.Serve()
		if err != nil {
			log.Fatal(err)
		}
	}()

	<-done
	server.Log.Warn("caught signal, stopping...")
	_ = server.Close()
	server.Log.Info("mochi mqtt shutdown complete")
}