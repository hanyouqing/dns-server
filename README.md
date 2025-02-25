# dns-server
Learn DNS by creating a DNS Server.

## Quick Start

```shell
sudo apt update
sudo apt install golang-go

alias  togo="cd ${HOME}/go"
export GOPATH="${HOME}/go"
export GOBIN="${HOME}/go/bin"
export PATH=$PATH:"${HOME}/go/bin"

cd $GOPATH
# git@github.com:hanyouqing/dns-server.git
git clone https://github.com/hanyouqing/dns-server.git
cd dns-server
go mod init dns-server
```
