package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"

	apiproto "github.com/caojunxyz/mimi-api/proto"
	pb "github.com/golang/protobuf/proto"
)

// FromRequest extracts the user IP address from req, if present.
func RequestClientIP(req *http.Request) (net.IP, error) {
	realIp := req.Header.Get("X-Real-IP")
	// ip, _, err := net.SplitHostPort(req.RemoteAddr)
	// if err != nil {
	// 	return nil, fmt.Errorf("userip: %q is not IP:port", req.RemoteAddr)
	// }
	//
	// userIP := net.ParseIP(ip)
	// if userIP == nil {
	// 	return nil, fmt.Errorf("userip: %q is not IP:port", req.RemoteAddr)
	// }
	// return userIP, nil
	return net.ParseIP(realIp), nil
}

func ParseHttpRequest(w http.ResponseWriter, r *http.Request, msg pb.Message) (accountId int64, ip net.IP, err error) {
	ip, _ = RequestClientIP(r)
	log.Println("ip:", ip)
	if msg != nil {
		defer r.Body.Close()
		var data []byte
		data, err = ioutil.ReadAll(r.Body)
		if err != nil {
			return
		}
		dataFormat := r.Header.Get("dataFormat")
		if dataFormat == "json" {
			err = json.Unmarshal(data, msg)
		} else {
			err = pb.Unmarshal(data, msg)
		}
	}
	ctx := r.Context()
	accountId, _ = ctx.Value("accountId").(int64)
	if err != nil {
		log.Println(ip, accountId, r.URL.Path, err)
		http.Error(w, "请求解析错误!", http.StatusForbidden)
	}
	return
}

func GetRequestVersion(r *http.Request) string {
	return r.Header.Get("version")
}

func GetRequestDeviceId(r *http.Request) string {
	return r.Header.Get("deviceId")
}

func WriteHttpResponse(w http.ResponseWriter, r *http.Request, code apiproto.RespCode, desc string, result pb.Message) {
	resp := &apiproto.Response{
		Code: code,
		Desc: desc,
		Api:  r.URL.Path,
	}
	if result != nil {
		dataFormat := r.Header.Get("dataFormat")
		if dataFormat == "json" {
			resp.Result, _ = json.Marshal(result)
		} else {
			resp.Result, _ = pb.Marshal(result)
		}
	}
	data, _ := pb.Marshal(resp)
	w.Write(data)
}

//-------------------------------------------------------------------------------------------------------------------------
func IsDirExists(path string) bool {
	file, err := os.Open(path)
	if err != nil {
		log.Println(path, err)
		return false
	}
	fileInfo, err := file.Stat()
	if err != nil {
		log.Println(path, err)
		return false
	}
	return fileInfo.IsDir()
}

func JoinInt32List(list []int32, sep string) string {
	ret := ""
	for i, v := range list {
		ret += fmt.Sprint(v)
		if i < len(list)-1 {
			ret += sep
		}
	}
	return ret
}

func ParseLotteryIdArg(r *http.Request) (apiproto.LotteryId, error) {
	arg := path.Base(r.URL.Path)
	n, err := strconv.Atoi(arg)
	if err != nil {
		return 0, err
	}
	return apiproto.LotteryId(n), nil
}

func TimeBeforeDays(n int) time.Time {
	now := time.Now()
	year, month, day := now.Date()
	hour, min, sec := now.Clock()
	return time.Date(year, month, day-n, hour, min, sec, 0, time.Local)
}
