# Slang

Slang is a dynamically typed functional programming language for expressing safe abstractions through pattern matching and preconditions. Slang is short for (s)imple (lang)uage. Typically referring to the use of informal speech and jargon, here ‘slang’ refers to the dynamic typing aspect of the language. 

```

# Non exhaustive patterns are ok
let fib = {
  0 -> 1
  1 -> 1
  n -> fib (n - 1) + fib (n - 2)
}

# Patterns without bodies return .true and .false
let range = {
  min max : (max > min) val : (val >= min && val < max)
}

# Partial application of patterns
let range_fn = range 80 90
let result = range_fn (fib 10)

# Labels and higher order patterns allow for an object-like pattern
let new_rectangle = {
  x y w h -> {
    .x -> x
    .y -> y
    .width -> w
    .height -> h
  }
}

# Is square tests that the passed pattern responds to labels .w and .h
let is_square = {
  rec : (rec.w == rec.h)
}



```

