
# USER LIBRARY

record = {
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

map = {
  _  []     -> []
  fn [x:xs] -> [fn x] ++ map fn xs
}

# PROGRAM

token_t = record [
  .line,
  .char,
  .value
]

scanners =
  scan_word = {
    word -> {
      [word:_] -> len word
               => 0
    }
  }

  scan_identifier = {
    [c:cs] : ((c >= "a" && c <= "z") ||
              (c >= "A" && c <= "Z") ||
               c == "_") -> 1 + scan_identifier cs
                         => 0
  }

  map scan_word ["{", "}"] ++ [scan_identifier]

_ = map { n -> print (n "{") } scanners

print "done!"