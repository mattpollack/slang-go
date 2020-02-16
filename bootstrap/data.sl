package data

# Arbitrary key value pairs are great but are better when they pattern match 
struct = {
  members ->
    next = {
      []     state v -> state v
      [m:ms] state v -> next ms {
        m -> v
        n -> state n
      }
    }
    next members { _ -> .no_record }
}

pair = struct [.is_pair, .a, .b] .true

none = [.none]
some = { x -> [.some, x] }

{
  .struct -> struct
  .pair   -> pair
  .none   -> none
  .some   -> some
}