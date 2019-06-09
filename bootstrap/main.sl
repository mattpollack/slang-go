
# STANDARD LIBRARY

let for = {
  s : (s.end) _  -> s
  s           fn -> (for (fn s) fn)
}

let responds_to = {
  val label : (match (val label) { .false -> .false => .true })
}

let printf = {
  ["%d":str] -> {
    # TODO: check it's a number
    v -> let _ = (print v) (printf str)
  }
  
  ["%s":str] -> {
    # TODO: check it's a string
    v -> let _ = (print v) (printf str)
  }

  ["%l":str] -> {
    # TODO: check it's a label
    v -> let _ = (print v) (printf str)
  }

  ["%L":str] -> {
    # TODO: check it's a list
    v ->
      let _ = (print "[")
      let _ = (for
        {
          .end  -> .false
          .list -> v
        }
        {
          s -> match (s.list) {
            [x:[]] ->
              let _ = (print x)
              { .end -> .true }

            [x:xs] ->
              let _ = (print x)
              let _ = (print " ")
              {
                .end  -> .false
                .list -> xs
              }
          }
        })
      let _ = (print "]")
      (printf str)
  }
  
  [x:str] ->
    let _ = (print x)
    (printf str)

  # Seems like this can be avoided
  => ""
}

# Cheap way to force a panic
let panic = {
  str ->
    let _ = (print str)
    ({ 0 -> 1} 1)
}

# COMPILER

let parser =
  let token_new = {
    kind val -> {
      .kind  -> kind
      .val   -> val
      .print -> (printf "(%d: '%s')\n" kind val)
    }
  }
    
  let tokens = [
    (token_new .TOKEN_KIND_PAREN_OPEN    "(")
    (token_new .TOKEN_KIND_PAREN_CLOSE   ")")
    (token_new .TOKEN_KIND_BRACKET_OPEN  "{")
    (token_new .TOKEN_KIND_BRACKET_CLOSE "}")
    (token_new .TOKEN_KIND_ARROW         "->")
    (token_new .TOKEN_KIND_PLUS          "+")
    (token_new .TOKEN_KIND_PLUS          "-")
  ]

  # Attempt to parse the next token
  let next_reserved = {
    src -> (for
      {
        .end    -> .false
        .tokens -> tokens
      }
      {
        s -> match (s.tokens) {
          []     -> { .end -> .true }
          [t:ts] -> match src {
            [(t.val):ss] ->
              { .end   -> .true
                .match -> t
                .src   -> ss
              }
            =>
              { .end    -> .false
                .tokens -> ts
              }
          }
        }
      })  
  }
  
  {
    src ->
      (print ((next_reserved src).match))

    # Always end the buffer
    => (token_new .TOKEN_KIND_BUFFER_END)
  }

let src =
"{
  0 -> 1
  1 -> 1
  n -> fib (n - 1) + fib (n - 2)
}"

(printf "%s\n" (parser src))