# web-radio-player
Go binary with an embedded web server that can play my favorite internet radios with `mpv`

# Running the player on a Raspberry Pi 3
## Compile

```shell
GOARCH=arm GOARM=7 GOOS=linux go build
```

## Create a systemd service
```
sudo cp web-radio.service /etc/systemd/system/
sudo systemctl enable web-radio.service
sudo systemctl start web-radio.service
```
