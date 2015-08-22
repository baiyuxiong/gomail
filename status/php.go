package status

import (
	"os/exec"
	"strings"
)
func PHPStatus() (string,string){
	message := ""
	status := PHPVersion()

	module := PHPModule()
	status += "\n" + module

	if !strings.Contains(module,"redis") || !strings.Contains(module,"memcached"){
		message = "请检查PHP的redis和memcached扩展"
	}

	return status,message
}

func PHPVersion() string {
	res,_ := exec.Command("php", "-v").Output()
	return string(res)
}

func PHPModule() string {
	res,_ := exec.Command("php", "-m").Output()
	return string(res)
}