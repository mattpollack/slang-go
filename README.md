# Slang

Slang is a dynamically typed functional programming language for expressing safe abstractions through pattern matching and preconditions. Slang is short for (s)imple (lang)uage. Typically referring to the use of informal speech and jargon, here ‘slang’ refers to the dynamic typing aspect of the language. 

```
let range = {
  min (max : max > min) (v : v >= min && v < max) 
  -> .true
  => .false
}

let test = range 5 10
let val = test 6
let val = test 11 # compile time error

println "No errors!"
```

