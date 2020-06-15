package sketch

scanners = module {
  or = {
    id -> id
  }
  and = {
    id -> id
  }
  many = {
    id -> id
  }
}

_ = print scanners

print (scanners.or "derp")