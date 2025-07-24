package main

import (
   "fmt"
)

func main() {
   var n, m, k int64
   if _, err := fmt.Scan(&n, &m, &k); err != nil {
       return
   }
   // Zero-based index of seat
   idx := k - 1
   // Each lane has m desks, each desk has 2 seats
   seatsPerLane := 2 * m
   // Determine lane number (1-indexed)
   lane := idx/seatsPerLane + 1
   // Position inside the lane (0 to seatsPerLane-1)
   pos := idx % seatsPerLane
   // Determine desk number (1-indexed)
   desk := pos/2 + 1
   // Determine side: left if even position, right if odd
   side := 'L'
   if pos%2 == 1 {
       side = 'R'
   }
   // Output lane, desk, and side
   fmt.Printf("%d %d %c", lane, desk, side)
}
