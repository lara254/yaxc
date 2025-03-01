# WEB API
A an api for making social medias involving a Scheme compiler

### Scheme supported

def ::= (define (<var> <formals>) <exp>)

exp ::= <var>
     |  <int>
     |  <bool>
     |  (let ((<var> <exp>)) <exp>)
     |  (if <exp> <exp> <exp>)
     |  (lambda <formals> <exp>)
     |  <prim>
     |  (while <exp> <exp>)
     |  (set! <var> <exp>)
     |  (begin <exp>*)

*a-Normal form*

aexp ::= <var>
      |  <bool>
      |  <int>
      |  (lambda <formals> <exp>)

cexp ::= (if <aexp> <exp> <exp>)
      |  (set <var> <exp>)
      |  (<aexp> <aexp> ...)

<exp> ::= (let ((<var> <cexp>)) <exp>)
       |  <aexp>
       |  <cexp>


### How to run it
```
$ make server 
```

```python
>>> url = "http://localhost:1234/api/compiler"
>>> re = requests.post(url, json={"exp":"(let ((i 0)) (begin (while (< i 4) (set i (+ i 1))) i))"})

>>> re.json()
{'exp': '[[movq 0 i] [label loop] [movq 1 %rax] [addq %rax i] [cmpq 4 i] [jl loop] [movq i %rdi] [callq print_int]]'}
>>> 

