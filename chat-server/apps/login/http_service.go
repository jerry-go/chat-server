package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
	"sync"
	"github.com/caojunxyz/mimi-server/apps/dbagent/proto"
	"google.golang.org/grpc"
)

type LoginServer struct {
	sync.RWMutex
	dbClient dbproto.DbActivityAgentClient
}

func NewServer() *LoginServer {
	srv := &LoginServer{}
	return srv
}



//初始化服务器及配置路由
func (srv *LoginServer) ServeHTTP() {
	defer func() {
		if err := recover(); err != nil {
			log.Println("recover from panic:", err)
		}
	}()

	log.Println("ServeHTTP on port ", *httpPort)
	mux := http.NewServeMux()
	mux.HandleFunc("/activity/list", srv.HandleActivityList)
	mux.HandleFunc("/activity/detail", srv.HandleActivityDetail)


	httpServer := http.Server{
		Addr:         fmt.Sprintf(":%d", *httpPort),
		Handler:      mux,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	if err := httpServer.ListenAndServe(); err != nil {
		panic(err)
	}
}
