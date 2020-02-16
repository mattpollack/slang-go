package main

import "./bootstrap/std.sl"
import "./bootstrap/data.sl"

lexer =
  token_t = data.struct [
    .value,
    .src
  ]

  scan_word = {
    prefix str ->
      scan_loop = {
        ""     suffix -> data.some suffix
        [p:ps] [p:ss] -> scan_loop ps ss
                      => data.none
      }

      match scan_loop prefix str {
        [.none]         -> data.none
        [.some, suffix] -> data.some (token_t prefix suffix)
      }
  }

  # Ordered by length
  scanners = std.map scan_word [
    "{",
    "}",
    "[",
    "]",
  ]

  lexer = {
    current src -> {
      .current -> current
      .next    ->
        skip_ws = {
          [" ":str] -> skip_ws str
          str       -> str
        }

        src = skip_ws src

        match src {
          "" -> lexer (token_t "EOF" "") ""
             =>
            next_token = std.foldr {
              scan [.none] -> scan src
              _    token   -> token
            } (data.none) scanners

            match next_token {
              [.none]        -> lexer (token_t "NO_MATCH" src) src
              [.some, token] -> lexer token                    (token.src)
            }
        }
    }
  }

  lexer (token_t "SOF" "")

parser = lexer "{{[[]]}}"

parser = parser.next
_ = print (parser.current.value)

parser = parser.next
_ = print (parser.current.value)

parser = parser.next
_ = print (parser.current.value)

parser = parser.next
_ = print (parser.current.value)

parser = parser.next
_ = print (parser.current.value)

parser = parser.next
_ = print (parser.current.value)

parser = parser.next
_ = print (parser.current.value)

parser = parser.next
_ = print (parser.current.value)

parser = parser.next
_ = print (parser.current.value)

parser = parser.next
_ = print (parser.current.value)

parser = parser.next
_ = print (parser.current.value)


print "ok"