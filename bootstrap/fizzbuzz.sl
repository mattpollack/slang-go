package fizzbuzz

import "./bootstrap/std.sl"
import "./bootstrap/data.sl"

fizzbuzz = {
  a ->
    affixes = [
      [3, "fizz"],
      [5, "buzz"]
    ]

    match (std.foldr {
      [n, affix] : ((a % n) == 0) str -> affix ++ str
                                _ str -> str
    } "" affixes) {
      "" -> a
       s -> s
    }
}

std.unfoldr {
  50 -> data.none
   n ->
     _ = print (fizzbuzz n)
     data.some [.nil, n + 1]
} 1