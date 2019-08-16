


# Why doesn't this work? Cannot find identifier test


test = {
    x -> {
      .push -> { y -> test (x ++ [y]) }
      .vals -> x
    }
  } []

test = test.push 5
_ = print (test.vals)
_ = print "\n"

print "# sketch\n"