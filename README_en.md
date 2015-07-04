# gomail
mail sender using redis as data storage written by go(golang.org)

# Install

````
go get githuc.com/baiyuxiong/gomail
cd $GOPATH/src/githuc.com/baiyuxiong/gomail
bower install
go run gomain.go
````

# Config
cfg.json
* httpAddress: Http address for web
* sender: Mail send method, can be "smtp" or "sendcloud"
* jobKey: List key in redis, used as "BRPOP $jobKey 0" to  get mail task from redis
* logKey: List key in redis, used as "LPUSH $logKey %log%" to  put mail log in redis, will be show in web page
