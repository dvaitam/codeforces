package main

import (
   "fmt"
)

func main() {
   var n, m, a, b int
   if _, err := fmt.Scan(&n, &m, &a, &b); err != nil {
       return
   }
   // Option 1: buy all rides individually
   cost1 := n * a
   // Option 2: buy as many m-ride tickets as possible, and remaining individually
   full := n / m
   rem := n % m
   cost2 := full*b + rem*a
   // Option 3: buy only m-ride tickets, possibly overbuying
   cost3 := (full + 1) * b
   // Minimum of the three
   minCost := cost1
   if cost2 < minCost {
       minCost = cost2
   }
   if cost3 < minCost {
       minCost = cost3
   }
   fmt.Println(minCost)
}
