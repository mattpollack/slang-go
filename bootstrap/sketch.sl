
# Linked list map

map_get = {
  key [x:xs] -> match x.key == key {
    .false -> map_get key xs
    .true  -> x.val
  }
  => .nil
}
  
map_set = {
  key val [x:xs] -> match x.key == key {
    .false -> [x] ++ map_set key val xs
    .true  -> [{
      .key -> key
      .val -> val
    }] ++ xs
  }
  => [{
    .key -> key
    .val -> val
  }]
}

map_new = {
  data -> {
    .get key -> map_get key data
    .set key -> {
      val -> map_new (map_set key val data)
    }
    .unset key -> map_new (map_set key .nil data)
  }
}

map_new = map_new []

map = map_new
map = map.set "a" 10
map