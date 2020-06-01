package parse

# NOTE: Because of the import bug
data = {
  .some -> { v -> [.some, v] }
  .none -> [.none]
}

# NOTE: Because of the import bug
std = {
  .map ->
    map = {
      _ []     -> []
      f [m:ms] -> [f m] ++ map f ms
    }
    map

  .foldr ->
    foldr = {
      _ z []     -> z
      f z [m:ms] -> f m (foldr f z ms)
    }
    foldr

  .unfoldr ->
    unfoldr = {
      f z ->
        match f z {
          [.some, [v, s]] -> [v] ++ unfoldr f s
                          => []
        }
    }
    unfoldr

  .find ->
    find = {
      _ []     -> data.none
      f [m:ms] ->
        match f m {
          [.none] -> find f ms
          some    -> some
        }
    }
    find

  .filter ->
    filter = {
      _ []             -> []
      f [m:ms] : (f m) -> [m] ++ filter f ms
      f [_:ms]         -> filter f ms
    }
    filter
}

is_char_class = {
  chars ->
    chars_list = std.foldr { c list -> [c] ++ list } [] chars
    {
      i : (std.foldr { c bool -> bool || c == i } .false chars_list)
    }
}

is_alpha = is_char_class "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
is_num = is_char_class "0123456789"
is_whitespace = is_char_class " \t\r"

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

tokenizer =
  scanners = [
    scan_class_t .token_identifier   is_alpha,
    scan_class_t .token_num          is_num,
    scan_class_t .token_whitespace   is_whitespace,
    scan_word_t  .token_arrow        "->",
    scan_word_t  .token_newline      "\n",
    scan_word_t  .token_equal        "=",
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

parser = tokenizer "

fib = {
  0 -> 1
  1 -> 1
  n -> fib (n - 1) + fib (n - 2)
}

print (fib 10)".next

# Just print every token value and type
tokens = std.filter {
  t : (t.type != .token_whitespace)
} (
  std.unfoldr {
    p : (p.curr.type == .end_of_file) -> data.none
    p                                 -> data.some [p.curr, (p.next)]
  } parser
)

std.map {
  t ->
    #_ = print (t.val)
    _ = print (t.type)
    t
} tokens
