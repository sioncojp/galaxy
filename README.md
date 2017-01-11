# Galaxy

* Galaxy provides you with target repository's http server as URL: `http://{{commit-number}}-{{url}}` by running container.

## Usage
```shell
### create database writing in config.yml
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
$ glide install
$ go build cmd/galaxy/*.go
$ ./galaxy --config=example/galaxy.yml
```

#### examples
* show `examples/README.md`

## config.yml
* You can execute external /bin/bash script after the repository will have changed a commit number.

## API
### create target repository
* curl http://localhost:8080/repository -X POST

### show running container list
* curl http://localhost:8080/container/list

### create/delete proxy server
$ curl http://localhost:8080/container_proxy -X POST
$ curl http://localhost:8080/container_proxy -X DELETE

### create/delete container
$ curl -F "commit_number=99c6894" http://localhost:8080/container/:commit_number -X POST
$ curl -F "commit_number=99c6894" http://localhost:8080/container/:commit_number -X DELETE
```

## Licence
MIT
