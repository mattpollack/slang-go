package sketch

f = {
  str ->
    g = {
      9 -> str
      n -> g (n + 1)
    }
    g
}

print ( f "done!" 0 )