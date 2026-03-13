# mini-asciinema-store

A minimalistic, lightweight server designed to store and serve [asciinema](https://asciinema.org/) recordings.

---

## Features
* **Simple & Fast:** Minimal setup required.
* **Storage Agnostic:** Stores recordings locally, in-memory, or on AWS S3, powered by the [gocloud.dev blob driver](https://gocloud.dev/howto/blob/).
* **Basic Authentication:** Supports token-based upload authentication out of the box.
* **HTTPS Support:** Optional built-in TLS for secure communication.
* **Containerized:** Ready-to-use Docker environment for effortless deployment.

---

## Installation

### Pre-compiled Binaries
You can download pre-compiled binaries from the [Releases page](https://github.com/Gerifield/mini-asciinema-store/releases).

### Building from Source

Ensure you have **Go 1.25** or newer installed.

```bash
git clone https://github.com/Gerifield/mini-asciinema-store.git
cd mini-asciinema-store
go build -o mini-asciinema-store src/cmd/server/server.go
```

The compiled binary `mini-asciinema-store` will be ready to run in your directory.

### Docker

You can easily build and run the server using Docker. The default container stores uploaded files in a `/uploads` folder in the root of the container.

```bash
docker build -t mini-store .
docker run -p 8080:8080 -v $(pwd)/uploads:/uploads mini-store:latest
```

---

## Configuration and Usage

Start the server using default settings or pass command-line flags.

### Available Flags

```text
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

### Storage Configuration

To store uploads in a local directory, create the folder and set `-uploadBucket`:

```bash
mkdir uploads
./mini-asciinema-store -uploadBucket="file://$(pwd)/uploads"
```

*Note: The server leverages the [Go Cloud Development Kit (Blob)](https://gocloud.dev/howto/blob/) to support different storage providers.*

---

## Asciinema Client Configuration

To use your self-hosted server with the `asciinema` CLI tool, update your `asciinema` configuration file (usually `~/.config/asciinema/config`) to point to your instance.

```ini
[api]
; API server URL, default: https://asciinema.org
; Set this to the address of your mini-asciinema-store instance
url = http://127.0.0.1:8080
```
*(Replace `http://127.0.0.1:8080` with your actual server base URL)*

Alternatively, you can override the API URL directly using an environment variable when recording:

```bash
ASCIINEMA_API_URL=http://127.0.0.1:8080 asciinema rec
```

---

## Authentication

You can restrict uploads to specific users by utilizing a simple token-based authentication mechanism.

1. **Locate your install ID:** By default, `asciinema` generates a unique installation ID located at `~/.config/asciinema/install-id`.
2. **Create an auth file:** Add the allowed install IDs (one per line) into a text file.

**Example:**
```bash
echo "<my_install_id>" > authFile.txt
./mini-asciinema-store -uploadBucket="file://$(pwd)/uploads" -authFile authFile.txt
```
