
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

test = map
test = test.set.derp 10

# _ = print_ast (test.get.derp)
_ = test.set.herp

print "\ndone!\n"


