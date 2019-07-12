# mini-asciinema-store
Small server for store asciinema recordings


# Usage

```
 -https
    	HTTPS enable
  -httpsCert string
    	HTTPS cert (default "server.crt")
  -httpsPrivateKey string
    	HTTPS private key (default "server.key")
  -listenAddr string
    	HTTP listening address (default ":8080")
  -serverBaseURL string
    	Base URL for the server (default "http://127.0.0.1:8080")
  -uploadPath string
    	Folder to store the uploaded files (default "uploads/")
```

For using the server just edit your asciinema config file (or create one): `~/.config/asciinema/config`

Then add the following config:

```
[api]

; API server URL, default: https://asciinema.org
; If you run your own instance of asciinema-server then set its address here
; It can also be overriden by setting ASCIINEMA_API_URL environment variable
url = http://127.0.0.1:8080
```

(Change the URL to your base URL.)

