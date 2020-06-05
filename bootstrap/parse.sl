package parse

import "bootstrap/data.sl"
import "bootstrap/std.sl"

is_char_class = {
  chars ->
    chars_list = std.foldr { c list -> [c] ++ list } [] chars
    {
      i : (std.foldr { c bool -> bool || c == i } .false chars_list)
    }
}

is_alpha = is_char_class "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
is_num = is_char_class "0123456789"
is_whitespace = is_char_class " \t\r\n"

token_t = {
  v s t -> {
    .val  -> v
    .src  -> s
    .type -> t
  } 
}

scan_class_t = {
  type fn in ->
    scan = {
      [c:cs] : (fn c) ->
        match scan cs {
          [str, src] -> [c ++ str, src]

        }
      cs -> ["",  cs]
    }
    match scan in {
      ["",  _]   -> data.none
      [str, src] -> data.some (token_t str src type)
    }
}

scan_word_t = {
  type word in ->
    scan = {
      ""     src    -> data.some src
      [w:ws] [w:is] -> scan ws is
                    => data.none
    }

    match scan word in {
      [.some, src] -> data.some (token_t word src type)
                   => data.none
    }
}

tokenizer_t =
  scanners = [
    scan_class_t .token_identifier   is_alpha,
    scan_class_t .token_num          is_num,
    scan_class_t .token_whitespace   is_whitespace,
    scan_word_t  .token_arrow        "->",

    # removing for now
    # scan_word_t  .token_newline      "\n",

    scan_word_t  .token_equals       "=",
    scan_word_t  .token_plus         "+",
    scan_word_t  .token_minus        "-",
    scan_word_t  .token_brace_open   "{",
    scan_word_t  .token_brace_close  "}",
    scan_word_t  .token_paren_open   "(",
    scan_word_t  .token_paren_close  ")"
  ]
  state = {
    token -> {
      .curr -> token
      .next : (token.src == "") ->
        state (token_t "" "" .end_of_file)
      .next ->
        next = std.find {
          scanner -> scanner (token.src)
        } scanners

        match next {
          [.some, next_token] -> state next_token
                              => state (token_t "" (token.src) .error_parsing)
        }
    }
  }
  {
    src -> state (token_t "" src .start_of_file)
  }

# Sample source code
source_tokenizer = tokenizer_t "abc  def  ghi

fib = {
  0 -> 1
  1 -> 1
  n -> fib (n - 1) + fib (n - 2)
}

print (fib 10)"

# scans a token type while skipping whitespace
scan_token = {
  token_type construct_fn in_tokenizer ->
    scan = {
      tokenizer ->
        tokenizer = tokenizer.next

        match tokenizer.curr {
          token : (token.type == token_type) ->
            data.some [
              tokenizer,
              construct_fn token
            ]
          => data.none
        }
    }

    match in_tokenizer.next.curr {
      token : (token.type == .token_whitespace) -> scan (in_tokenizer.next)
                                                => scan in_tokenizer
    }
}

scan_identifier = scan_token .token_identifier {
  token -> {
    .type  -> .ast_identifier
    .value -> token.val
  }
}

# scans the first of many scanners
scan_meta_or = {
  scanners tokenizer -> std.find { scanner -> scanner tokenizer } scanners
}

# scans all of many scanners 
scan_meta_and = {
  scanners tokenizer ->
    match
      (std.do {
        collection [next_tokenizer, token] -> [collection ++ [token], [next_tokenizer]]
      } scanners [] [tokenizer])
    {
      [.none]                                 -> data.none
      [.some, [collection, [next_tokenizer]]] -> data.some [collection, next_tokenizer]
    }
}

# scans many of one scanner
scan_meta_many = {
  scanner tokenizer ->
    _ = print "TODO scan_meta_many"
    data.none
}

scan_expression = {
  tokenizer ->
    std.find {
      scanner -> scanner tokenizer
    } [
      scan_identifier
    ]
}

_ = match scan_meta_and [
  scan_expression,
  scan_expression,
  scan_expression
] source_tokenizer {
  [.none]                               -> print ":("
  [.some, [[id1, id2, id3], tokenizer]] ->
    _ = print "-------------"
    _ = print (id1.value)
    _ = print (id2.value)
    _ = print (id3.value)
    _ = print "-------------"
    .nil
}

print "ok"






# grammar
# expression  = identifier
#             | number
#             | application
#             | let
# application = '(' expression+ ')'
#             | expression+ newline
# let         = identifier '=' expression newline expression

