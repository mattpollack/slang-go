package std

map = {
  _ []     -> []
  f [m:ms] -> [f m] ++ map f ms
}

filter = {
  _ []             -> []
  f [m:ms] : (f m) -> [m] ++ filter f ms
  f [_:ms]         -> filter f ms
}

foldr = {
  _ z []     -> z
  f z [m:ms] -> f m (foldr f z ms)
}

foldl = {
  _ z []     -> z
  f z [m:ms] -> foldl f (f z m) ms
}

unfoldr = {
  f z ->
    match f z {
      [.some, [v, s]] -> [v] ++ unfoldr f s
                      => []
    }
}

unfoldl = {
  f z ->
    match f z {
      [.some, [v, s]] -> unfoldl f s ++ [v]
                      => []
    }
}

# Reflection use-case
{
  .map     -> map
  .filter  -> filter
  .foldr   -> foldr
  .foldl   -> foldl
  .unfoldr -> unfoldr
  .unfoldl -> unfoldl
}