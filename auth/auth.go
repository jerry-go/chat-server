package auth

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/context"

	apiproto "github.com/caojunxyz/mimi-api/proto"
	"github.com/caojunxyz/mimi-server/utils"
	jwt "github.com/dgrijalva/jwt-go"
)

var KEY = []byte("Fa^%#$4542+099x")

const AUTH_HEADER_FIELD = "Authorization"

const (
	REVIEW_VERSION = "1.0.3"
	LATEST_VERSION = "1.0.3"
	APP_URL        = "http://www.baidu.com/"
)

func convertToIntVersion(ver string) int {
	list := strings.Split(ver, ".")
	if len(list) != 3 {
		return 0
	}
	str := strings.Replace(ver, ".", "", -1)
	n, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}
	return n
}

func checkVersion(version string) (bool, error) {
	clientIntVer := convertToIntVersion(version)
	if clientIntVer == 0 {
		return false, fmt.Errorf("无效版本号: %s", version)
	}
	reviewIntVer := convertToIntVersion(REVIEW_VERSION)
	if reviewIntVer == 0 {
		return false, fmt.Errorf("无效版本号: %s", REVIEW_VERSION)
	}
	latestIntVer := convertToIntVersion(LATEST_VERSION)
	if reviewIntVer == 0 {
		return false, fmt.Errorf("无效版本号: %s", LATEST_VERSION)
	}

	log.Println(clientIntVer, reviewIntVer, latestIntVer)
	isForceUpgrade := (clientIntVer != reviewIntVer && clientIntVer < latestIntVer)
	return isForceUpgrade, nil
}

type Claims struct {
	AccountId int64  `json:"accountId"`
	DeviceId  string `json:"deviceId"`
	jwt.StandardClaims
}

func SetHeader(w http.ResponseWriter, accountId int64, deviceId string) {
	claims := Claims{
		accountId,
		deviceId,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24 * 15).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "cp.kxkr.com",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(KEY)
	if err != nil {
		log.Println(err)
		return
	}
	// log.Printf("token: %s\n", signedToken)
	w.Header().Add(AUTH_HEADER_FIELD, signedToken)
}

// TODO: before hook, after hook
func Validate(protected http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// log.Println("validate...")
		authHeader := r.Header.Get(AUTH_HEADER_FIELD)
		if authHeader == "" {
			log.Println("缺少token:", r.URL.Path)
			http.Error(w, "Token验证未通过", http.StatusForbidden)
			return
		}
		// log.Printf("token: %s\n", authHeader)
		claims := Claims{}
		token, err := jwt.ParseWithClaims(authHeader, &claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method %v", token.Header["alg"])
			}
			return KEY, nil
		})

		if err != nil {
			log.Println(err)
			http.Error(w, "Token验证未通过", http.StatusForbidden)
			return
		}

		appVersion := r.Header.Get("appVersion")
		log.Printf("appVersion: %s", appVersion)
		isForceUpgrade, err := checkVersion(appVersion)
		// if err != nil {
		// 	log.Println(err)
		// 	http.Error(w, err.Error(), http.StatusForbidden)
		// 	return
		// }
		w.Header().Add("reviewVersion", REVIEW_VERSION)
		w.Header().Add("latestVersion", LATEST_VERSION)
		w.Header().Add("APP_URL", APP_URL)
		if isForceUpgrade {
			log.Println("Force Upgrade:", appVersion, REVIEW_VERSION, LATEST_VERSION)
			info := &apiproto.StringValue{Value: APP_URL}
			desc := fmt.Sprintf("需要升级到%s版本才能继续使用", LATEST_VERSION)
			utils.WriteHttpResponse(w, r, apiproto.RespCode_Upgrade, desc, info)
			return
		}

		if claims, ok := token.Claims.(*Claims); ok && token.Valid {
			ctx := context.WithValue(context.Background(), "accountId", claims.AccountId)
			protected(w, r.WithContext(ctx))
			return
		}
		http.Error(w, "Token验证未通过", http.StatusForbidden)
	})
}

func WsValidate(protected http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		r.ParseForm()
		authArr := r.Form["_token"]

		if len(authArr) != 1 {
			log.Println("缺少token:", r.URL.Path)
			http.Error(w, "Token验证未通过", http.StatusForbidden)
			return
		}
		log.Println("_token:", authArr, "len token", len(authArr[0]))
		authHeader := authArr[0]

		claims := Claims{}
		token, err := jwt.ParseWithClaims(authHeader, &claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method %v", token.Header["alg"])
			}
			return KEY, nil
		})
		if err != nil {
			log.Println(err)
			http.Error(w, "Token验证未通过", http.StatusForbidden)
			return
		}

		if claims, ok := token.Claims.(*Claims); ok && token.Valid {
			ctx := context.WithValue(context.Background(), "accountId", claims.AccountId)
			protected(w, r.WithContext(ctx))
			return
		}
		http.Error(w, "Token验证未通过", http.StatusForbidden)
	})
}
