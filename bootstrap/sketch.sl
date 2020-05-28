package sketch

f = {
  str ->
    g = {
      9 -> str
      n -> g (n + 1)
    }
    g
}

foldr = {
  _ z []     -> z
  f z [m:ms] -> f m (foldr f z ms)
}

foldl = {
  _ z []     -> z
  f z [m:ms] -> foldl f (f z m) ms
}

#{
#  .foldr -> foldr
#  .foldl -> foldl
#}

print ( f "done!" 0 )