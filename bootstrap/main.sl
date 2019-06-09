
# STANDARD LIBRARY

let do = {
  s : (state.end) _  -> s
  s               fn -> (do (fn s) fn)
}

let range = {
  min max : (max > min) v : (&& (>= v min) (< v max))
}

#let has_prefix = {
#  str prefix :
#  (&& (>= (str.len) (prefix.len))
#      (== (str[:(prefix.len)]) prefix))
#}

# COMPILER

# parse for 

let printf = {
  #"" -> ""
  #str : (&& (>= (str.len) 2) (has_prefix str "%")) ->
  #  (print "TODO: thing")
  #[x:xs] -> 
}

let derp = {
  [3:xs] -> (print "derp\n")
  [x:xs] -> xs
}

let herp = {
  ["te":xs] -> xs
  [x:"erp"] -> x
}

let _ = (print_ast (herp "te"))
let _ = (print_ast (herp "derp"))
_