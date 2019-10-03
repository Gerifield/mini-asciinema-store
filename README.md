# mini-asciinema-store
Small server for store asciinema recordings

## Building

Install Go on your system and clone the repository:

```
git clone https://github.com/Gerifield/mini-asciinema-store.git
cd mini-asciinema-store
go build -o mini-asciinema-store src/cmd/server/server.go
```

The you have your binary there. Just run it! :)

## Usage

```
  -authFile string
    	Authenticated tokens in a file (one token per line)
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
  -uploadBucket string
    	Folder or bucket URL to store the uploaded files (supports: file, mem, s3) (default "file:///uploads/")
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


If you'd like to use a folder to store the uploads, just create one then set the uploadBucket at the start for example:

`mini-asciinema-store -uploadBucket="file://$(pwd)/uploads"`

The server uses the https://gocloud.dev/howto/blob/ driver for the storage.

## Authentication

You could use a simple text file with one token per line to allow only certain users to upload.
Create a new file and put the client's install ID in that.

You could find the install ID by default here: `~/.config/asciinema/install-id`

Example:

```
echo "<my install id>" > authFile.txt
go run src/cmd/server/server.go -uploadBucket="file://$(pwd)/uploads" -authFile authFile.txt
```

This will set the upload target to the `uploads` folder in this directory and set the `authFile.txt`.