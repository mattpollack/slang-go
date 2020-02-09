

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

token = struct [
  .char,
  .line,
  .val
]

next = token 1 1 "{"

_ = print (next.char)
_ = print (next.line)
_ = print (next.val)

print "done\n"

#derp = {
#  v -> {
#    v -> v + 1
#  }
#}

#print (derp 10 10)