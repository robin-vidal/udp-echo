# udp echo

> [!NOTE]
> Learning project for my distributed systems roadmap.
> Exploring UDP reliability, timeouts and fault detection by hand.
> See my [kvgo](https://github.com/robin-vidal/kvgo) key-value store for the applied work.

UDP echo server and client in Go. Implements manual retry logic and packet loss simulation to feel what TCP gives for free.

```bash
# Server
go run ./cmd/server

# Client
go run ./cmd/client
```

## Flags

### Client
| Flag | Default | Description |
|---|---|---|
| `--addr` | `localhost` | Server address |
| `--port` | `8080` | Server port |
| `--msg` | `Hello World!` | Message to send |
| `--timeout` | `2000` | Read/write timeout (ms) |
| `--retries` | `3` | Retries on timeout |
| `--bufferSize` | `1024` | Buffer size (bytes) |

### Server
| Flag | Default | Description |
|---|---|---|
| `--port` | `8080` | Port to listen on |
| `--bufferSize` | `1024` | Buffer size (bytes) |
