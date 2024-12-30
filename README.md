A social media site for compiler enthusiasts

### How to run it
```
$ make main
```

#### Compiler Webservice
The compiler web service will be core to the social media :-)

```python
>>> url = "http://localhost:1234/api/compiler"
>>> re = requests.post(url, json={"exp":"(let ((i 0)) (begin (while (< i 4) (set i (+ i 1))) i))"})

>>> re.json()
{'exp': '[[movq 0 i] [label loop] [movq 1 %rax] [addq %rax i] [cmpq 4 i] [jl loop] [movq i %rdi] [callq print_int]]'}
>>> 

