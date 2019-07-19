mkdir pb
mkdir client
protoc -I ./  $1 --go_out=./pb --sanrpc-client_out=./client --sanrpc-server_out=./  --sanrpc-main_out=./
