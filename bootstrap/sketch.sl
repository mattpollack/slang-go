# This shouldn't cause a fail to match arguments
# Create a shortcut...
let derp = {
  s : (s.derp) -> .false
  s : (s.herp) -> .true
}

(print_ast (derp { .herp -> .true }))