module server

go 1.24.0

toolchain go1.24.7

replace server => ../server

require (
	github.com/go-chi/chi v1.5.5
	github.com/go-chi/chi/v5 v5.2.3
	github.com/lib/pq v1.10.9
	github.com/redis/go-redis/v9 v9.13.0
	golang.org/x/crypto v0.42.0
)

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
)
