clean:
	rm parser y.output compiler server
yacc:
	 goyacc -o parser.go parser.y

server:
	go build -o server server.go compiler.go parser.go
compiler:
	go build -o compiler main.go compiler.go parser.go



