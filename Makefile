clean:
	rm parser y.output compiler server
yacc:
	 goyacc -o src/compiler/parser.go src/compiler/parser.y

server:
	go build -o server src/server.go

compiler:
	go build -o compiler src/main.go
test:
	cd tests && go test 


