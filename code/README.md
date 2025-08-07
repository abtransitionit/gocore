# syntax
```shell
cd folder
# ... = recursively
go vet ./...
```


# build
```shell
cd folder
# ... = recursively
go build ./...
```


# to know
- no `main()` no `CLI`. code is build into `go env GOCACHE`