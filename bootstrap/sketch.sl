package main

match [1, 2, [1, 2, 3]]
{
  [1, 2, [1: x]] -> print x
}