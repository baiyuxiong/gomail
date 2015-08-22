package status
import (
	"html/template"
	"time"
	"github.com/baiyuxiong/gomail/sender"
	"github.com/baiyuxiong/gomail/model"
	"encoding/json"
	"strings"
	"github.com/go-xweb/log"
)

var Substitution_vars string = `{
    "to": [
        "ok@baiyuxiong.com"
    ],
    "sub": {
        "%messages%": [
            "{{message}}"
        ]
    }
}`

func CheckInterval()  {
	timer := time.NewTicker(15 * time.Minute)
	for {
		select {
		case <-timer.C:
			data := Check()
			messagesInterface,exists := data["messages"]
			if exists{
				log.Println("CheckInterval: checking message")
				messages := messagesInterface.([]string)
				if len(messages) >0{
					message := ""
					for _,v := range messages{
						message += v+"<br/>"
					}
					m := model.EmailSendCloud{}
					m.MailType = "batch"
					m.Template_invoke_name = "server_error"
					m.Substitution_vars =strings.Replace(Substitution_vars,"{{message}}",message,-1)
					mail,_ := json.Marshal(m);
					sender.SendEmailBySendcloud(string(mail));
				}
			}else{
				log.Println("CheckInterval: messages not exists")
			}
		}
	}
}


func Check() map[string]interface{} {
	messages := make([]string,0)

	CPUStats,CPUStatsMsg := CPUStats()
	if len(CPUStatsMsg) >0{
		messages = append(messages,CPUStatsMsg)
	}

	FSInfos,FSInfosMsg := FSInfos()
	if len(FSInfosMsg) >0{
		messages = append(messages,FSInfosMsg)
	}

	memcachedStatus,MemcachedStatusMsg := MemcachedStatus()
	if len(MemcachedStatusMsg) >0{
		messages = append(messages,MemcachedStatusMsg)
	}

	redisStatus,RedisStatusMsg := RedisStatus()
	if len(RedisStatusMsg) >0{
		messages = append(messages,RedisStatusMsg)
	}

	phpStatus,PHPStatusMsg := PHPStatus()
	if len(PHPStatusMsg) >0{
		messages = append(messages,PHPStatusMsg)
	}

	mysqlStatus,MysqlStatusMsg := MysqlStatus()
	if len(MysqlStatusMsg) >0{
		messages = append(messages,MysqlStatusMsg)
	}

	httpdStatus,HttpdStatusMsg := HttpdStatus()
	if len(HttpdStatusMsg) >0{
		messages = append(messages,HttpdStatusMsg)
	}

	data := map[string]interface{}{
		"messages" : messages,
		"memcached":template.HTML(memcachedStatus),
		"redis":template.HTML(redisStatus),
		//"hostInfos":template.HTML(status.HostInfos()),
		"PHPStatus":template.HTML(phpStatus),
		"MysqlStatus":template.HTML(mysqlStatus),
		"HttpdStatus":template.HTML(httpdStatus),
		"CPUStats":template.HTML(CPUStats),
		"FSInfos":template.HTML(FSInfos),
		"MemStats":template.HTML(MemStats()),
		"NetIOStats":template.HTML(NetIOStats()),
		"ProcessStats":template.HTML(ProcessStats()),
		"PagesStats":template.HTML(PagesStats()),
	}

	return data
}