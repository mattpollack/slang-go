package parse

data = {
  .some -> { v -> [.some, v] }
	.none -> [.none]
}

std = {
  .map ->
    map = {
      _ []     -> []
      f [m:ms] -> [f m] ++ map f ms
    }
    map

  .foldr ->
    foldr = {
      _ z []     -> z
      f z [m:ms] -> f m (foldr f z ms)
    }
    foldr
}

token = {
  v s -> {
    .val -> v
    .src -> s
  }
}

is_alpha =
  alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
  cs = std.foldr { c list -> [c] ++ list } [] #alphabet
  _ = print cs
{
  c : (
    c == "a" ||
    c == "b" ||
    c == "c"
  )
}

test =
  print (is_alpha "c")

{}
