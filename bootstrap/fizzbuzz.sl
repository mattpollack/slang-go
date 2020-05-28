package fizzbuzz

#import "./bootstrap/std.sl"
#import "./bootstrap/data.sl"

none = [.none]
some = { x -> [.some, x] }

foldr = {
  _ z []     -> z
  f z [m:ms] -> f m (foldr f z ms)
}

unfoldr = {
  f z ->
    match f z {
      [.some, [v, s]] -> [v] ++ unfoldr f s
                      => []
    }
}

fizzbuzz = {
  a ->
    affixes = [
      [3, "fizz"],
      [5, "buzz"]
    ]

    match (foldr {
      [n, affix] : ((a % n) == 0) str -> affix ++ str
                                _ str -> str
    } "" affixes) {
      "" -> a
       s -> s
    }
}

unfoldr {
  100 -> none
   n ->
     _ = print (fizzbuzz n)
     some [.nil, n + 1]
} 1