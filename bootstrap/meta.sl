






# Holds values patterned by labels
struct = {
  state []     n -> state n
  state [m:ms] n -> struct {
    m -> n
    o -> state o
  } ms
}

struct = struct { _ -> .nil }

