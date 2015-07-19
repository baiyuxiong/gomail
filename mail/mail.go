package mail

import (
	"github.com/garyburd/redigo/redis"
	"github.com/baiyuxiong/gomail/config"
	"log"
	"math"
	"time"
	"encoding/json"
	"github.com/baiyuxiong/gomail/model"
)

var mailClient *MailClient

type MailClient struct {
	MailConn redis.Conn
	LogConn  redis.Conn
}

func (m *MailClient)newRedisConn() (redis.Conn, error) {
	c, err := redis.Dial(config.Config().Redis.Network, config.Config().Redis.Address)
	if err != nil {
		log.Println("initRedisConn mc - error : ", err.Error())
		return nil, err
	}

	if len(config.Config().Redis.Password) > 0 {
		if _, err := c.Do("AUTH", config.Config().Redis.Password); err != nil {
			c.Close()
			log.Println("initRedisConn mc - Do AUTH error : ", err.Error())
			return nil, err
		}
	}
	return c, nil
}

func (m *MailClient)insureMailConn() {
	if m.MailConn != nil {
		return
	}

	var retry int = 1

	for {
		conn, err := m.newRedisConn()
		if err == nil {
			m.MailConn = conn
			log.Println("initRedisConn - MailConn connect to redis ok")
			return
		}

		if retry > 8 {
			retry = 1
		}
		time.Sleep(time.Duration(math.Pow(2.0, float64(retry))) * time.Second)
		retry++
	}
}

func (m *MailClient)insureLogConn() {
	if m.LogConn != nil {
		return
	}

	var retry int = 1

	for {
		conn, err := m.newRedisConn()
		if err == nil {
			m.LogConn = conn
			log.Println("initRedisConn - MailConn connect to redis ok")
			return
		}

		if retry > 8 {
			retry = 1
		}
		time.Sleep(time.Duration(math.Pow(2.0, float64(retry))) * time.Second)
		retry++
	}
}

func (m *MailClient) closeMailConn() {
	m.MailConn.Close()
	m.MailConn = nil
}

func (m *MailClient) closeLogConn() {
	m.LogConn.Close()
	m.LogConn = nil
}

func (m *MailClient)StartSendMail() {
	log.Println("startSendMail....")
	for {
		log.Println("conn to redis... ")

		m.insureMailConn()
		m.insureLogConn()

		log.Println("startSendMail - BRPOP... ")
		data, err := m.MailConn.Do("BRPOP", config.Config().JobKey, 0)

		if err != nil {
			log.Println("startSendMail - redis BRPOP error : ", err.Error())
			time.Sleep(time.Second*10)
		}else {
			log.Println("startSendMail - get job, start send...")

			log.Printf("data is %s", data)
			mails := data.([]interface{})
			if len(mails) == 2 {
				go m.Send(string(mails[1].([]byte)))
			}else {
				log.Printf("%s", mails)
			}
		}
	}
}

func (m *MailClient)Send(data string) {
	for _,s := range senders {
		if s.Name == config.Config().Sender{
			l := s.Run(data)
			m.PutLog(l)
			break;
		}
	}
}

func (m *MailClient)PutLog(data model.EmailLog) {
	bs, err := json.Marshal(data)
	if err != nil {
		log.Println("PutLog - err : ", err.Error())
		return
	}
	m.LogConn.Do("LPUSH",config.Config().LogKey,string(bs))
}

func (m *MailClient)GetMailLog(start int,stop int) (reply interface{}, err error) {
	log.Println("GetMailLog...")

	m.insureLogConn()

	reply, err = m.LogConn.Do("LRANGE", config.Config().LogKey, start, stop)
	if err != nil {
		log.Println("GetMailLog - reply ", reply, ",error : ", err.Error())
	}

	return
}

func GetMailLog(start int,stop int) (reply interface{}, err error) {
	return mailClient.GetMailLog(start,stop)
}


func Start() {
	mailClient = &MailClient{}
	mailClient.StartSendMail()
}
