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

# string helpers
strings = {
  .split -> {
    "" _ ->
      {
        .head -> ""
        .tail -> ""
      }
    src 0 ->
      {
        .head -> ""
        .tail -> src
      }
    [x:xs] i ->
      rest = strings.split xs (i - 1)
      {
        .head -> x ++ rest.head
        .tail -> rest.tail
      }
  }
}

map = {
  []     _  -> []
  [x:xs] fn -> [fn x] ++ map xs fn
}

find = {
  []     _  -> .nil
  [x:xs] fn -> match fn x {
    .true -> x
          => find xs fn
  }
}

# MAIN

parser_new =
  is_alpha = {
    x : ((x >= "a" && x <= "z") || (x >= "A" && x <= "Z"))
  }

  scanner = {
    fn [x:xs] ->
      match fn x {
        .true -> 1 + scanner fn xs
              => 0
      }
  }

  scanner_literal = {
    [l:ls] [x:xs] -> match l == x {
      .true -> 1 + scanner_literal ls xs
            => 0
    }

    => 0
  }

  scanners = [
    { .scan -> scanner { x -> is_alpha x }; .kind -> .TOKEN_KIND_IDENTIFIER  },
    { .scan -> scanner_literal "{";         .kind -> .TOKEN_KIND_BRACE_OPEN  },
  ]

  parser_new = {
    src -> {
      next -> find scanners {
        scanner ->
          _ = printf "{}\n" (scanner.scan src)
          .false
      }
    }
  }

  parser_new

parser = parser_new "{}->1"
_ = print (parser.next)

print ("\n## main.sl\n")