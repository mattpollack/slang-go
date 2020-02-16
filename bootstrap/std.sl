package std

{
  .map ->
    map = {
      _ []     -> []
      f [m:ms] -> [f m] ++ map f ms
    }
    map

  .filter ->
    filter = {
      _ []             -> []
      f [m:ms] : (f m) -> [m] ++ filter f ms
      f [_:ms]         -> filter f ms
    }
    filter

  .foldr ->
    foldr = {
      _ z []     -> z
      f z [m:ms] -> f m (foldr f z ms)
    }
    foldr

  .foldl ->
    foldl = {
      _ z []     -> z
      f z [m:ms] -> foldl f (f z m) ms
    }
    foldl

  .struct ->
    struct = {
      state []     n -> state n
      state [m:ms] n -> struct {
        m -> n
        o -> state o
      } ms
    }
    struct { _ -> .no_record }
}