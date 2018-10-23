
----------------
-- Run server --
----------------
* Development 
... go/src/github.com/hoopra/api-base_go/go run *.go

* Production
GO_ENV=prod
... go/src/github.com/hoopra/api-base_go/go run *.go

* Overwrites for key paths
AUTH_PRIVATE_KEY_PATH={key to private key}
AUTH_PUBLIC_KEY_PATH={key to public key}

----------------
-- Race check --
----------------
... go/src/github.com/hoopra/api-base_go/go run *.go -race


----------------
-- Run tests  --
----------------
... go/src/github.com/hoopra/api-base_go/test/go run test.go
