buildWindows:
	go build -ldflags "-s -w" .

buildMac:
	set GOOS=darwin && set GOARCH=amd64 && go build -ldflags "-s -w" .

run:
	make buildWindows
	client-proxy-clientside -user_id=test
