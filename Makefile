all:
	mkdir -p out/Release
	GOOS=linux GOARCH=amd64 go build -o out/Release/wepkg_linux_amd64 ./cmds
	GOOS=linux GOARCH=arm go build -o out/Release/wepkg_linux_armeabi ./cmds
	GOOS=linux GOARCH=arm64 go build -o out/Release/wepkg_linux_arm64 ./cmds
	GOOS=darwin GOARCH=amd64 go build -o out/Release/wepkg_darwin ./cmds
	GOOS=windows GOARCH=amd64 go build -o out/Release/wepkg.exe ./cmds
