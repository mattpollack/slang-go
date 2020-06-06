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

atoi = {
  n []     -> some n
  n [c:cs] ->
    n = n*10

    # NOTE: no better way of approaching this?
    m = match c {
      "0" -> 0
      "1" -> 1
      "2" -> 2
      "3" -> 3
      "4" -> 4
      "5" -> 5
      "6" -> 6
      "7" -> 7
      "8" -> 8
      "9" -> 9
          => none
    }

    match m {
      none -> none
           => atoi (n + m) cs
    }
}

{
  .record -> record
  .pair   -> pair
  .none   -> none
  .some   -> some
  .error  -> error
  .atoi   -> atoi 0
}