# How to try test-app-rails?
* test-app-rails is http server written by rails.

## 1. Docker pull circusd and nginx-proxy-socket
```shell
$ docker pull ruby:2.3
$ docker pull sioncojp/nginx-proxy-socket:latest
```

## 2. Setting config.yml and External shell
```shell
$ vi /tmp/config.yml
workdir: "/tmp/galaxy"
server:
  port: 8080
database:
  driver: "mysql"
  host: "localhost"
  port: 3306
  dbname: "galaxy"
  user: "root"
  password: ""
github:
  repository: "git@github.com:sioncojp/test-app-rails.git"
  name: "test-app-rails"
url: "test-app.com"
script: "/tmp/hoge.sh"
docker:
  image: "ruby"
  tag: "2.3"
  proxyimage: "sioncojp/nginx-proxy-socket"
  proxytag: "latest"
  exec: "cd /tmp/test-app-rails && bundle install --path .bundle && bundle exec unicorn_rails -c config/unicorn.rb -D"

### "/tmp/galaxy/bundle" has run "bundle install" on ruby:2.3-Container and copy In advance
### after $ curl http://localhost:8080/repository -X POST.
$ vi /tmp/hoge.sh
cp -R /tmp/galaxy/test-app-rails $1/
cp -R /tmp/galaxy/bundle $1/test-app-rails/.bundle
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
$ curl -F "commit_number=f546501" http://localhost:8080/container/:commit_number -X POST
$ curl -F "commit_number=327da9b" http://localhost:8080/container/:commit_number -X POST
$ curl http://localhost:8080/container/list
$ curl http://localhost:8080/url/list
```

## 5. show from browser after edit /etc/hosts
```shell
$ vi /etc/hosts
127.0.0.1 f546501-test-app.com
127.0.0.1 327da9b-test-app.com

### curl or browser
$ curl http://f546501-test-app.com/hello/index
It Works!!!

$ curl http://327da9b-test-app.com/hello/index
It Hard Works!!!
```

## 6. if you want to remove galaxy
```shell
$ curl -F "commit_number=f546501" http://localhost:8080/container/:commit_number -X DELETE
$ curl -F "commit_number=327da9b" http://localhost:8080/container/:commit_number -X DELETE
$ curl http://localhost:8080/container_proxy -X DELETE

### delete repository directory
* rm to workdir written by config.yml

### delete database
$ mysql -uroot
> drop database galaxy
```
