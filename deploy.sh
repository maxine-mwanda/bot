# 1. Stop the service on the server
echo "Stop the service on the server"
ssh root@108.61.198.87 "systemctl stop telegrambot"
# 2. Copy from local to server
echo "Compile and send to server"
go build main.go
scp main root@108.61.198.87:/root/telegram/main
rm main
# 4. Restart the service
echo "Restart service on server"
ssh root@108.61.198.87 "systemctl restart telegrambot"
ssh root@108.61.198.87 "tail -f /root/telegram/logs/*"

