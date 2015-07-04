package maillog

import (
	"github.com/baiyuxiong/gomail/model"
	"github.com/baiyuxiong/gomail/config"
	"github.com/garyburd/redigo/redis"
	"encoding/json"
	"log"
)

func PutLog(conn redis.Conn,l model.EmailLog) {
	bs, err := json.Marshal(l)
	if err != nil {
		log.Println("PutLog - err : ", err.Error())
		return
	}
	conn.Do("LPUSH",config.Config().LogKey,string(bs))
}