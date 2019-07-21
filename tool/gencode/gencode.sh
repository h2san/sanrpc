mkdir pb
mkdir client
protoc -I ./  $1 --go_out=./pb --sanrpc-client_out=./client --sanrpc-server_out=./  --sanrpc-main_out=./

mod=$1
modname=`cut -d '.' ${mod} -f 0`
go mod init ${modname}