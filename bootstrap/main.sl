package main

import "./bootstrap/std.sl"
import "./bootstrap/data.sl"

lexer =
  lexer_t = data.struct [
    .prev_token,
    .src
  ]

  scan_word = {
    prefix str ->
      scan_loop = {
        ""     suffix -> data.some suffix
        [p:ps] [p:ss] -> scan_loop ps ss
      }

      match scan_loop prefix str {
        [.none]         -> data.none
        [.some, suffix] -> data.some (lexer_t prefix suffix)
      }
  }

  _ = match scan_word "hello" "hello, world!" {
    [.none]        -> print ":("
    [.some, lexer] ->
      _ = print (lexer.prev_token)
      print (lexer.src)
      
  }

  std.map scan_word ["{", "}", "[", "]"]

_ = print lexer     

print "ok"