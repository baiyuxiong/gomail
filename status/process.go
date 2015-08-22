package status
import (
	"github.com/baiyuxiong/gomail/utils"
	"os/exec"
	"bytes"
)

func MemcachedStatus() (string,string) {
	c1 := exec.Command("ps", "aux")
	c2 := exec.Command("grep", "memcached")

	c2.Stdin, _ = c1.StdoutPipe()

	var b2 bytes.Buffer
	c2.Stdout = &b2

	c1.Start()
	c2.Start()
	c1.Wait()
	c2.Wait()

	return utils.Nl2br(b2.String()),""
}

func RedisStatus() (string,string) {
	c1 := exec.Command("ps", "aux")
	c2 := exec.Command("grep", "redis")

	c2.Stdin, _ = c1.StdoutPipe()

	var b2 bytes.Buffer
	c2.Stdout = &b2

	c1.Start()
	c2.Start()
	c1.Wait()
	c2.Wait()

	return utils.Nl2br(b2.String()),""
}