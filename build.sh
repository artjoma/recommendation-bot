go build -ldflags "-s -w" -o build/recommendation-bot
#env GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -o build/recommendation-bot.exe
#env GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w" -o build/recommendation-bot-darwin

# use UPX compressor
if [ ! -z $1 ]; then
   upx build/recommendation-bot
fi