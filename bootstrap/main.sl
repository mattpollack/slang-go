# STANDARD LIBRARY

false = .false
true = .true

printf = {
  format -> match format {
    ["{}":xs] -> {
      x ->
        _ = print x
        _ = printf xs
        format
    }

    [x:xs] ->
      _ = print x
      _ = printf xs
      format

    => format
  }
}

for = {
  s : (s.end) fn -> s
  s           fn -> for (fn s) fn
}

map = {
  []     _  -> []
  [x:xs] fn -> [fn x] ++ map xs fn
}

# MAIN

parser =
  token_new = {
    kind value -> { .kind -> kind; .value -> value }
  }

  tokens = [
    token_new .TOKEN_KIND_BRACE_OPEN  "{",
    token_new .TOKEN_KIND_BRACE_CLOSE "}"
  ]

  _ = map tokens {
    x -> printf "{} \n" (x.value)
  }

  {
    .parse -> {
      src ->  src
    }
  }

parser.parse "{}->1"