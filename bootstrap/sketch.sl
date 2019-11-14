

# recursion binds to expression

map = {
  state .get s -> state s
  state .set s -> {
    v -> map {
      s -> v
      k -> state k
    }
  }
}

map = map { _ -> .nil }

_ = print_ast map     

test = map

_ = print_ast (test.set.derp)

#test = test.set.derp 10
#_ = print_ast (test.get.derp)

print_ast "\ndone!\n"


