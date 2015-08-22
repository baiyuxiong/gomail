package status

import (
	"github.com/baiyuxiong/gomail/utils"
	"os/exec"
	"bytes"
)

func checkProcess(processName string) (string,string) {
	c1 := exec.Command("ps", "aux")
	c2 := exec.Command("grep", processName)
	c3 := exec.Command("grep", "-v","grep")

	c2.Stdin, _ = c1.StdoutPipe()
	c3.Stdin, _ = c2.StdoutPipe()

	var b3 bytes.Buffer
	c3.Stdout = &b3

	c1.Start()
	c2.Start()
	c3.Start()
	c1.Wait()
	c2.Wait()
	c3.Wait()

	message := ""
	result := b3.String()
	if len(result) < 1{

		message = "进程"+processName+"未运行"
	}

	return utils.Nl2br(b3.String()),message
}

func MemcachedStatus() (string,string) {
	return checkProcess("memcached")
}

func RedisStatus() (string,string) {
	return checkProcess("redis")
}

func MysqlStatus() (string,string) {
	return checkProcess("mysql")
}

func HttpdStatus() (string,string) {
	return checkProcess("httpd")
}