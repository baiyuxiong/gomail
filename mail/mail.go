package mail

import (
	"github.com/garyburd/redigo/redis"
	"github.com/baiyuxiong/gomail/config"
	"log"
	"time"
	"os"
)

var MailConn redis.Conn
var LogConn redis.Conn

func GetMailLog() (reply interface{}, err error){
	log.Println("GetMailLog...")
	reply,err = LogConn.Do("LRANGE",config.Config().LogKey,0,20)
	if err != nil{
		log.Println("initRedisConn - reply ", reply, ",error : ", err.Error())
	}

	return
}

func Start() {

	initRedisConn()
	defer func(){
		MailConn.Close()
		MailConn = nil
	}()

	startSendMail()
}

func initRedisConn() {
	mc, err := redis.Dial(config.Config().Redis.Network, config.Config().Redis.Address)
	if err != nil {
		log.Println("initRedisConn mc - error : ", err.Error())
	}

	lc, err := redis.Dial(config.Config().Redis.Network, config.Config().Redis.Address)
	if err != nil {
		log.Println("initRedisConn lc - error : ", err.Error())
	}

	if len(config.Config().Redis.Password) > 0{
		if _, err := mc.Do("AUTH", config.Config().Redis.Password); err != nil {
			mc.Close()
			log.Println("initRedisConn mc - Do AUTH error : ", err.Error())
			os.Exit(0)
			return
		}
		if _, err := lc.Do("AUTH", config.Config().Redis.Password); err != nil {
			lc.Close()
			log.Println("initRedisConn lc - Do AUTH error : ", err.Error())
			os.Exit(0)
			return
		}
	}
	MailConn = mc
	LogConn = lc
	log.Println("initRedisConn - connect to redis ok")
}

func startSendMail() {
	log.Println("startSendMail....")
	for {
		if nil == MailConn {
			initRedisConn()
			time.Sleep(time.Second*2)
		}

		data,err := MailConn.Do("BRPOP",config.Config().JobKey,0)

		if err != nil {
			log.Println("startSendMail - redis BRPOP error : ", err.Error())
			time.Sleep(time.Second*10)
		}else{
			log.Println("startSendMail - get job, start send...")

			log.Printf("data is %s",data)
			mails := data.([]interface{})
			if len(mails) == 2{
				go send(string(mails[1].([]byte)))
			}else{
				log.Printf("%s",mails)
			}
		}
	}
}