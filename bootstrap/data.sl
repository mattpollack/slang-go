package data

none = [.none]
some = { x -> [.some, x] }

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