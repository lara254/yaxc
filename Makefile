clean:
	rm parser y.output compiler server
yacc:
	 goyacc -o parser.go parser.y

main:
	go build -o server  main.go compiler.go parser.go
compiler:
	go build -o compiler compiler.go parser.go



