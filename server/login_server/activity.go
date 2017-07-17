package main

import (
	"net/http"
	"log"

	"github.com/caojunxyz/mimi-server/apps/dbagent/proto"
	"golang.org/x/net/context"
	"github.com/caojunxyz/mimi-server/utils"
	"github.com/caojunxyz/mimi-api/proto"
	"strconv"
)

//用户登录
func (srv *LoginServer) HandleLogin(w http.ResponseWriter, r *http.Request)  {
	log.Println("活动详情api.........")
	activityAccount := &dbproto.ActivityAccount{}
	activityAccount.AccountId = 1
	activityAccount.ActivityId = 1
	activityDetail ,err := srv.dbClient.ActivityDetail(context.Background(),activityAccount)
	if err != nil {
		log.Printf("%+v\n",err)
		utils.WriteHttpResponse(w, r, apiproto.RespCode_Fail, "获取活动详情失败", nil)
		return
	}
	log.Println(activityDetail)
	//utils.WriteHttpResponse(w, r, apiproto.RespCode_Success, "", activityList)
}
