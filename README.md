build and run
````shell
CGO_ENABLED=0 go build -o test/tc.exe && cd test && ./tc.exe init -n ch-future-world-devl -t gitops -s go-backend && cd ..
````
run with go
````shell
LOG_LEVEL=debug go run main.go
````