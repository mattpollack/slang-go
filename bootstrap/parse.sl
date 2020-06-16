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

alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
numbers = "0123456789"

is_num = is_char_class numbers 
is_alpha = is_char_class alphabet
is_identifier = is_char_class (alphabet ++ "_")
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
    scan_class_t .token_identifier   is_identifier,
    scan_class_t .token_number       is_num,
    scan_class_t .token_whitespace   is_whitespace,
    scan_word_t  .token_arrow        "->",

    # removing for now
    # scan_word_t  .token_newline      "\n",

    scan_word_t  .token_equals        "=",
    scan_word_t  .token_plus          "+",
    scan_word_t  .token_minus         "-",
    scan_word_t  .token_multiply      "*",
    scan_word_t  .token_divide        "/",
    scan_word_t  .token_brace_open    "{",
    scan_word_t  .token_brace_close   "}",
    scan_word_t  .token_paren_open    "(",
    scan_word_t  .token_paren_close   ")",
    scan_word_t  .token_bracket_open  "[",
    scan_word_t  .token_bracket_close "]"
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

# scans a token type while skipping whitespace
scan_token = {
  token_type make_fn in_tokenizer ->
    loop = {
      tokenizer ->
        tokenizer = tokenizer.next

        match tokenizer.curr {
          token : (token.type == token_type) -> data.some [tokenizer, make_fn token]
                                             => data.none
        }
    }

    match in_tokenizer.next.curr {
      token : (token.type == .token_whitespace) -> loop (in_tokenizer.next)
                                                => loop in_tokenizer
    }
}

# scans the first of many scanners
scan_or = {
  scanners make_fn tokenizer ->
    match (std.find { scanner -> scanner tokenizer } scanners) {
      [.none]                               -> data.none
      [.some, [next_tokenizer, collection]] -> data.some [next_tokenizer, make_fn collection]
    }
}

# scans all of many scanners 
scan_and = {
  scanners make_fn tokenizer ->
    match
      (std.do {
        collection [next_tokenizer, ast] -> [[next_tokenizer], collection ++ [ast]]
      } scanners [] [tokenizer])
    {
      [.none]                                 -> data.none
      [.some, [[next_tokenizer], collection]] -> data.some [next_tokenizer, make_fn collection]
    }
}

# scans many of one scanner
scan_many = {
  scanner make_fn in_tokenizer ->
    loop = {
      tokenizer -> match scanner tokenizer {
        [.none]                        -> [tokenizer, []]
        [.some, [next_tokenizer, ast]] -> match loop next_tokenizer {
          [final_tokenizer, collection] -> [final_tokenizer, [ast] ++ collection]
        }
      }
    }

    match loop in_tokenizer {
      [_, []]                     -> data.none
      [out_tokenizer, collection] -> data.some [out_tokenizer, make_fn collection]
    }
}

scan = module tokenizer {
    identifier =
      scan_token .token_identifier {
        token -> {
          .type  -> .identifier
          .value -> token.val
        }
      }

    number =
      scan_token .token_number {
        token ->
          match data.atoi (token.val) {
            [.none]    -> data.none
            [.some, v] -> {
              .type  -> .number
              .value -> v
            }
          }
      }

    let = 
      scan_and [
        scan.identifier,
        scan_token .token_equals { id -> id },
        scan.expression,
        scan.expression
      ] {
        [id, _, value, body] -> {
          .type  -> .let
          .id    -> id
          .value -> value
          .body  -> body
        }
      }

    application =
      scan_and [
        scan_token .token_paren_open { id -> id },
        scan_many (scan.expression) { id -> id },
        scan_token .token_paren_close { id -> id }
      ] {
        [_, body, _] -> {
          .type -> .application
          .body -> body
        }
      }

    pattern =
      scan_and [
        scan_token .token_brace_open { id -> id },
        scan_many (
          scan_and [
            scan_many (scan.pmatch) { id -> id },
            scan_token .token_arrow { id -> id },
            scan.expression
          ] {
            [_matches, _, body] -> {
              .type    -> .match
              .matches -> _matches
              .body    -> body
            }
          }
        ) { id -> id },
        scan_token .token_brace_close { id -> id }
      ] {
        [_, _matchGroups, _] -> {
          .type        -> .pattern
          .matchGroups -> _matchGroups
        }
      }

    list =
      scan_and [
        scan_token .token_bracket_open { id -> id },
        scan_token .token_bracket_close { id -> id }
      ] {
        _ ->
          _ = print "TODO: parse list"
          data.none
      }

    expression =
      scan_or [
        scan.pattern,
        scan.application,
        scan.list,
        scan.let,
        scan.identifier,
        scan.number
      ] { id -> id }

    # NOTE: unfortunate naming overload
    pmatch =
      scan_or [
        scan.identifier,
        scan.number,
        scan.list
      ] { id -> id }
}

_ = print scan

# sample source code
source_tokenizer = tokenizer_t (
  std.foldr {
    str line
    -> str ++ "\n" ++ line
  } "" [
    "fib = {                                           ",
    "  0 -> 1                                          ",
    "  1 -> 1                                          ",
    "  n -> (plus (fib (minus n 1)) (fib (minus n 2))) ",
    "}                                                 ",
    "_ = [1, (a b c), 2]                               ",
    "print (fib 5)                                     "
  ]
)

_ = match scan.expression source_tokenizer {
  [.none]                   -> print ":("
  [.some, [tokenizer, ast]] ->
    _ = print "-------------"
    _ = print (ast.body)
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

