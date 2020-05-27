# The slang Programming Language

Slang is a dynamically typed pure and eager functional programming language with higher order pattern matching and simple syntax. Slang derives its name from (s)imple (lang)uage, and its dynamic nature speaks to its informal semantics. 

## Overview

Slang source files have the `.sl` extension, and are composed of a package name, imports, and a slang expression. 

```
package main

import "std"

# This is a comment
print "Hello World!"
```

Slang has these common primitive values:
 - Integer
 - Decimal
 - String
 - List
 - Labels
 - Patterns (aka lambdas with a twist)

Additionally slang supports recursive and non recursive binding.

```
# Identifiers are bound to values like so:
# <identifier> = <expression>
# <expression>

a = 1
b = 1.5
c = "Hello World!
d = [a, b, c]
e = .true
f = {
  n -> "I'm a pattern!"
}

print "This is all you need!"
```

Labels are static enum like values which start with `.`. Labels bind tightly when not followed by a space.

```
# True and false are just labels
.true == (1 == 1)

data = {
  .a -> "Data!"
}

# data.a is evaluated first because of label binding rules
print data.a
```

Patterns open up deconstruction of these values.

```
add_pair = {
  [a, b] -> a + b
}

add_list = {
  # Multiple matches are typically lined up at the arrows
  []     -> 0
  [a:as] -> a + add_list as
}

data = [1, 2]

add_pair data == add_list data
```
