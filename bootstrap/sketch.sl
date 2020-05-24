package sketch


# Exploring defunctionalization
_ = 
  a = 10

  f = {
    a ->
      a = 15
      {
        b -> a + b
      }
    => 1
  }

  print (f 10 1)

print "ok"