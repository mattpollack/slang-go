package main

import "./bootstrap/std.sl"

token_t = std.struct [
  .line,
  .char,
  .value
]

scanners =
  scan_word = {
    word -> {
      [word:next] -> [word, next]
                  => ["",   next]
    }
  }

  std.map scan_word ["{", "}"]

_ = print scanners

print "ok"