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

# Reflection use-case
{
  .map    -> map
  .filter -> filter
  .foldr  -> foldr
  .foldl  -> foldl
}