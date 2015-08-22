package status

import (
	"github.com/akhenakh/statgo"
	"github.com/baiyuxiong/gomail/utils"
)

var s *statgo.Stat

func init()  {
	s = statgo.NewStat()
}

func HostInfos() (string) {
	//hi := s.HostInfos()
	return utils.Nl2br("")
}

func CPUStats() (string,string) {
	message := ""
	v := s.CPUStats()
	if v.Idle < 10{
		message = "CPU压力较大"
	}
	return utils.Nl2br(v.String()),message
}

func FSInfos() (string,string) {
	vs := s.FSInfos()

	message := ""
	status := ""
	for _,v := range vs{
		status += "\n"+v.String()
		if v.Size>0 && (v.Free == 0 || ( v.Free >0 && v.Size/v.Free > 9)){
			message = "硬盘空间紧张，请检查"
		}
	}
	return utils.Nl2br(status),message
}

func MemStats() string {
	v := s.MemStats()
	return utils.Nl2br(v.String())
}

func NetIOStats() string {
	vs := s.NetIOStats()

	status := ""
	for _,v := range vs{
		status += "\n"+v.String()
	}
	return utils.Nl2br(status)
}

func ProcessStats() string {
	v := s.ProcessStats()
	return utils.Nl2br(v.String())
}

func PagesStats() string {
	v := s.PageStats()
	return utils.Nl2br(v.String())
}