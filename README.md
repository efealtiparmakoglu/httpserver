# httpserver

**HTTP/1.1 server on raw TCP sockets — no net/http, zero dependencies beyond Go stdlib.**

[![License: MIT](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

Uses `net.Listen` directly, parses HTTP request lines and headers manually,
routes to handlers, and writes responses with proper Content-Length.

## Usage

```sh
httpserver -p 8080
curl localhost:8080/
curl localhost:8080/health
curl localhost:8080/hello/world
```

## License

MIT — see [LICENSE](LICENSE).
