package main

import "./bootstrap/std.sl"

nums = {
  list 0 -> list
  list i -> [i] ++ nums list (i - 1)
}

_ = print (std.filter { n : (n > 15) } (nums [] 30))

print "ok"