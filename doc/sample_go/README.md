# How to try test-app?
* test-app is http server written by Golang.
* So this tour need Go development environment to execute `go build`

## 1. Docker pull circusd and nginx-proxy-socket
```shell
$ docker pull sioncojp/circusd:latest
$ docker pull sioncojp/nginx-proxy-socket:latest
```

## 2. Setting config.yml and External shell
```
$ vi /tmp/config.yml
# reference examples/config.yml

$ vi /tmp/hoge.sh
GOOS=linux GOARCH=amd64 go build *.go
mv test-app $1/app
```

## 3. Running App
```shell
### create local database
$ mysql -uroot
> create database galaxy
> create table commits(
id bigint(20) unsigned NOT NULL PRIMARY KEY AUTO_INCREMENT,
number char(40) UNIQUE NOT NULL,
created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
INDEX idx_number(number)
);

### glide install & run app
$ go install github.com/sioncojp/galaxy/cmd/galaxy
$ galaxy -c /tmp/config.yml
```

## 4. API Request to galaxy
```shell
$ curl http://localhost:8080/repository -X POST
$ curl http://localhost:8080/container_proxy -X POST
$ curl -F "commit_number=99c6894" http://localhost:8080/container/:commit_number -X POST
$ curl -F "commit_number=28ea6b0" http://localhost:8080/container/:commit_number -X POST
$ curl http://localhost:8080/container/list
$ curl http://localhost:8080/url/list
```

## 5. show from browser after edit /etc/hosts
```shell
$ vi /etc/hosts
127.0.0.1 99c6894-test-app.com
127.0.0.1 28ea6b0-test-app.com

### curl or browser
$ curl http://99c6894-test-app.com/
It works!

$ curl http://28ea6b0-test-app.com/
It hard works!
```

## 6. if you want to remove galaxy
```shell
$ curl -F "commit_number=99c6894" http://localhost:8080/container/:commit_number -X DELETE
$ curl -F "commit_number=28ea6b0" http://localhost:8080/container/:commit_number -X DELETE
$ curl http://localhost:8080/container_proxy -X DELETE

### delete repository directory
* rm to workdir written by config.yml

### delete database
$ mysql -uroot
> drop database galaxy
```
