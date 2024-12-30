clean:
	rm parser y.output
yacc:
	 goyacc -o parser.go parser.y

main:
	go build -o server  main.go compiler.go parser.go




