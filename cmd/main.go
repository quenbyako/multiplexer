package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/quenbyako/multiplexer"
)

const defaultMaxConn = 100

func main() {
	// вообще можно поудобнее конфигурацию сделать, но будем считать что это просто тестовый стенд
	port, err := strconv.Atoi(os.Getenv("MULTIPLEXER_PORT"))
	if err != nil {
		log.Fatalf("MULTIPLEXER_PORT: %v", err)
	}

	var maxConn int
	maxConnStr := os.Getenv("MULTIPLEXER_MAX_CONNECTIONS")
	if maxConnStr != "" {
		maxConn, err = strconv.Atoi(maxConnStr)
		if err != nil {
			log.Fatalf("MULTIPLEXER_MAX_CONNECTIONS: %v", err)
		}
	} else {
		maxConn = defaultMaxConn
	}

	//

	//

	operator := multiplexer.NewDefaultOperator()

	listener, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		log.Fatalf("Start listening on :%v", port)
	}
	listener = LimitListener(listener, maxConn) // golang.org/x/net/netutil.LimitListener()
	defer listener.Close()

	srv := &http.Server{
		Handler: multiplexer.SetupHandler(multiplexer.NewDefaultRouter(operator)),
	}

	println("start server")
	setupServerCanceller(srv)
	err = srv.Serve(listener)
	if err != nil && err != http.ErrServerClosed {
		log.Fatalf("HTTP server error: %v", err)
	}
}

const shutdownPeriodSeconds = 3 * time.Second

func setupServerCanceller(srv *http.Server) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs

		ctx, cancel := context.WithTimeout(context.Background(), shutdownPeriodSeconds)
		defer cancel()

		err := srv.Shutdown(ctx)
		if err != nil {
			log.Fatalf("Shutdowning server: %v", err)
		}
	}()

}
