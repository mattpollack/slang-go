# Non exhaustive patterns are ok
fib = {
  0 -> 1
  1 -> 1
  n -> fib (n - 1) + fib (n - 2)
}

# Patterns without bodies return .true and .false
range = {
  min max : (max > min) val : (val >= min && val < max)
}

# Partial application of patterns
range_fn = range 80 90
result = range_fn (fib 10)

# Labels and higher order patterns allow for an object-like pattern
new_rectangle = {
  x y w h -> {
    .x -> x
    .y -> y
    .width -> w
    .height -> h
  }
}

# Is square tests that the passed pattern responds to labels .w and .h
is_square = {
  rec : (rec.w == rec.h)
}

print "ok!"