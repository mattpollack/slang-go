
# STANDARD LIBRARY

let for = {
  s : (s.end) _  -> s
  s           fn -> (for (fn s) fn)
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
          # Print the last element without a space then end
          s : (== (s.list.len) 1) ->
            let _ = (print (s.list.head))
            { .end -> .true }

          # Print each element and continue
          s : (>  (s.list.len) 0) ->
            let _ = (print (s.list.head))
            let _ = (print " ")
            {
              .end  -> .false
              .list -> (s.list.tail)
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
      .print -> (printf "{(%d): '%s'}" kind val)
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

  let _ = (printf "%L\n" tokens)

  {
    ""  -> (token_new .TOKEN_KIND_BUFFER_END)
    src -> (panic "TODO PARSER")
  }

let src = "{
  0 -> 1
  1 -> 1
  n -> fib (n - 1) + fib (n - 2)
}"

(printf "%s\n" (parser src))