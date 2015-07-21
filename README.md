# gomail
使用go语言编写的邮件发送器，读取redis list中的数据发送邮件。提供WEB页面查看发送日志。

Doc in English click [here](https://github.com/baiyuxiong/gomail/blob/master/README_en.md)

# 安装

````
go get github.com/baiyuxiong/gomail
cd $GOPATH/src/github.com/baiyuxiong/gomail
bower install
go run gomain.go
````

# 配置
cfg.json

* httpAddress: 查看邮件日志的WEB地址
* username: 登录用户名
* password: 登录密码
* sender: 邮件发送方式，支持"smtp"或者"sendcloud"
* jobKey: redis中的读取邮件记录的LIST名称，通过"BRPOP $jobKey 0" 的方式使用
* logKey: redis中存取邮件发送日志的LIST名称，通过"LPUSH $logKey %log%"方式存入日志

# 截图
![](http://baiyuxiong.com/download/screenshot.png)
