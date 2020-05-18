package main

import "./bootstrap/std.sl"
import "./bootstrap/data.sl"

fizzbuzz = {
  n ->
    pairs = [
      [3, "fizz"],
      [5, "buzz"]
    ]

    # Interesting parse case
    match (std.foldr {
      [m, a] : (match n % m {0}) s -> a ++ s
                               _ s -> s
    } "" pairs) {
      "" -> n
       s -> s
    }
}

print (std.unfoldr {
  30 -> data.none
  n -> data.some [fizzbuzz n, n + 1]
} 1)