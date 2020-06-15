package std

module {
  map = {
    _ []     -> []
    f [m:ms] -> [f m] ++ map f ms
  }

  filter = {
    _ []             -> []
    f [m:ms] : (f m) -> [m] ++ filter f ms
    f [_:ms]         -> filter f ms
  }

  find = {
    _ []     -> [.none]
    f [m:ms] ->
      match f m {
        [.none] -> find f ms
        some    -> some
      }
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

  apply = {
    o []     -> o
    f [m:ms] -> apply (f m) ms
  }

  # maybe broader interface?
  do = {
    next ->
      loop = {
        []       collection args -> [.some, [args, collection]]
        [fn:fns] collection args ->
          match apply fn args {
            [.none]      -> [.none]
            [.some, out] -> match next collection out {
              [next_args, next_collection] -> loop fns next_collection next_args
            }
          }
      }
      loop
  }  
}