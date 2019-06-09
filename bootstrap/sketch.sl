let test = {
  [(+ 0 1):[]] -> (print "ok cool\n")
  [x:xs] -> (print "ok dope\n")
}

let _ = (test [1 2])
let _ = (test [1])
(print "ok")