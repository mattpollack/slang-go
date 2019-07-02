option = {
  .none -> { .none -> .true }
  .some -> { n .some -> n }
}

match (option.some 10) {
  v : (v.none) -> (print "none")
  v : (v.some) -> (print (v.some))
}