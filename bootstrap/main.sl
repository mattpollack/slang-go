
# USER LIBRARY

struct = {
  members ->
    next = {
      []     state v -> state v
      [m:ms] state v -> next ms {
        m -> v
        n -> state n
      }
    }
    next members { _ -> .doesnt_exist }
}

map = {
  _  []     -> []
  fn [x:xs] -> [fn x] ++ map fn xs
}

# PROGRAM

print "done!"