module inu1255/sshdog

go 1.13

require (
	github.com/GeertJohan/go.rice v1.0.0
	github.com/google/shlex v0.0.0-20191202100458-e7afc7fbc510
	github.com/matir/sshdog v0.0.0-20200109212941-94a466579cda
	golang.org/x/crypto v0.0.0-20200709230013-948cd5f35899
	golang.org/x/sys v0.0.0-20190412213103-97732733099d
)

replace github.com/matir/sshdog v0.0.0-20200109212941-94a466579cda => ./
