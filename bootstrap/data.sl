package data

# Arbitrary key value pairs are great but are better when they pattern match 
record = {
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

none = [.none]
some = { x -> [.some, x] }

error = {
  .fail    v -> [.fail, v]
  .success v -> [.success, v]
}

{
  .record -> record
  .pair   -> pair
  .none   -> none
  .some   -> some
  .error  -> error
}