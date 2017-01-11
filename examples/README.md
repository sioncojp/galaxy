# How to try test-app?
* test-app is http server written by Golang.

## 1. Make a circusd container image
```shell
### run container
$ docker pull centos:6.8
$ docker run --name circusd -it centos:6.8

### operation in the container
### install python2.7
$ yum install gcc-c++ zlib-devel openssl-devel -y
$ cd /usr/local/src
$ curl -O https://www.python.org/ftp/python/2.7.13/Python-2.7.13.tgz
$ tar zxf Python-2.7.13.tgz
$ cd Python-2.7.13
$ ./configure
$ make && make altinstall
$ mv -f /usr/local/src/Python-2.7.13/python /usr/bin/python

### change python version for yum
$ vi /usr/bin/yum
#!/usr/bin/python2.6

### install pip
$ curl -kL https://bootstrap.pypa.io/get-pip.py | python

### install circus
$ pip install circus==0.13

### Setting circus
$ mkdir /var/log/circusd
$ chmod 775 /var/log/circusd
$ mkdir /etc/circusd/
$ chmod 755 /etc/circusd/

### Setting circus.ini
$ touch /etc/circusd/circus.ini && chmod 755 /etc/circusd/circus.ini

$ vi /etc/circusd/circus.ini
[circus]
[watcher:webapp]
uid = root
cmd = /tmp/test-app --fd $(circus.sockets.web)
numprocesses = 1
use_sockets = True
stop_signal = SIGINT
rlimit_nofile = 65536

stdout_stream.filename = /var/log/circusd/test-app_stdout.log
stdout_stream.class = TimedRotatingFileStream
stdout_stream.rotate_when = D
stdout_stream.rotate_interval = 1
stdout_stream.max_bytes = 1073741824
stdout_stream.backup_count = 5

stderr_stream.filename = /var/log/circusd/test-app_stderr.log
stderr_stream.class = TimedRotatingFileStream
stderr_stream.rotate_when = D
stderr_stream.rotate_interval = 1
stderr_stream.max_bytes = 1073741824
stderr_stream.backup_count = 5

[socket:web]
so_reuseport = True
umask = 000
replace = True

###
$ yum install initscripts -y
$ vi /etc/init.d/circusd
#!/bin/bash
#
# circusd This scripts turns circusd on
#
# Author: Shohei.Koyama
#
# chkconfig: - 95 04

# source function library
. /etc/rc.d/init.d/functions

CONFIG='/etc/circusd/circus.ini'
PIDFILE='/var/run/circusd.pid'
LOG_OUTPUT='/var/log/circusd/circusd.log'
RETVAL=0

start() {
    echo -n $"Starting circusd: "
    daemon /usr/local/bin/circusd --log-output $LOG_OUTPUT --pidfile $PIDFILE $CONFIG --daemon
    RETVAL=$?
    echo
    [ $RETVAL -eq 0 ] && touch /var/lock/subsys/circusd
}

stop() {
    echo -n $"Stopping circusd: "
    killproc circusd
    echo
    [ $RETVAL -eq 0 ] && rm -f /var/lock/subsys/circusd
}

restart() {
    stop
    start
}

case "$1" in
  start)
    start
    ;;
  stop)
    stop
    ;;
  restart|force-reload|reload)
    restart
    ;;
  condrestart)
    [ -f /var/lock/subsys/circusd ] && restart
    ;;
  status)
    status circusd
    RETVAL=$?
    ;;
  *)
    echo $"Usage: $0 {start|stop|status|restart|reload|force-reload|condrestart}"
    exit 1
esac

exit $RETVAL

$ chmod 755 /etc/init.d/circusd
$ exit

### create images
$ docker commit circusd circusd:1.0
$ docker rm circusd
```

## 2. Make a nginx proxy container image
```shell
### settings nginx.tmpl(in local pc)
$ mkdir /tmp/hoge
$ vi /tmp/hoge/nginx.tmpl
# reference examples/nginx.tmpl

$ docker run --name proxy -v /var/run/docker.sock:/tmp/docker.sock -v /tmp/hoge/:/tmp/ -it jwilder/nginx-proxy /bin/bash

### operation in the container
### setting ssl key
$ cd /tmp
$ openssl genrsa 2048 > server.key
$ openssl req -new -key server.key > server.csr
$ openssl x509 -days 100000 -req -signkey server.key < server.csr > server.crt
$ mv server.* /etc/nginx/conf.d/

### cp nginx.tmpl
cp /tmp/nginx.tmpl nginx.tmpl
exit

### create images
$ docker commit proxy galaxy-proxy:1.0
$ docker rm proxy
```

## 3. Setting External shell
```
$ vi /tmp/hoge.sh
GOOS=linux GOARCH=amd64 go build *.go
mv test-app $1/
```

## 4. Running App
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
$ brew install glide
$ export GO15VENDOREXPERIMENT=1
$ go get github.com/sioncojp/galaxy
$ cd $GOPATH/src/github.com/sioncojp/galaxy
$ glide install
$ go run cmd/galaxy/galaxy.go -c examples/config.yml
```

## 5. API Request to galaxy
```shell
$ curl http://localhost:8080/repository -X POST
$ curl http://localhost:8080/container_proxy -X POST
$ curl -F "commit_number=99c6894" http://localhost:8080/container/:commit_number -X POST
$ curl -F "commit_number=28ea6b0" http://localhost:8080/container/:commit_number -X POST
$ curl http://localhost:8080/container/list
```

## 6. show from browser after edit /etc/hosts
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

## 7. if you want to remove galaxy
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
