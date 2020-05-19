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
  }

  print (f 10 1)

# Turns into
#_ = 
#  a = 10

#  _f = {
#    _a b -> _a + b
#  }

#  f = {
#    a ->
#      a = 15
#      _f a
#  }

#  print (f 10 1)

print "ok"