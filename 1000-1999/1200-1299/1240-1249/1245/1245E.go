package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
)

var h [10][10]int
var b [10][10]int
var x [100]int
var y [100]int
var done [100][2]bool
var memo [100][2]float64

// solve returns the expected number of moves from position p with state k
func solve(p, k int) float64 {
   if done[p][k] {
       return memo[p][k]
   }
   done[p][k] = true
   if p == 0 && k == 0 {
       memo[p][k] = 0.0
       return memo[p][k]
   }
   if k == 0 {
       c := 6.0
       memo[p][k] = 6.0
       for s := 1; s <= 6; s++ {
           if p-s < 0 {
               c -= 1.0
           } else {
               memo[p][k] += solve(p-s, 1)
           }
       }
       memo[p][k] /= c
       return memo[p][k]
   }
   memo[p][k] = solve(p, 0)
   jump := h[x[p]][y[p]]
   if jump != 0 {
       newP := b[x[p]-jump][y[p]]
       memo[p][k] = math.Min(memo[p][k], solve(newP, 0))
   }
   return memo[p][k]
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   // Read the board
   for i := 0; i < 10; i++ {
       for j := 0; j < 10; j++ {
           fmt.Fscan(reader, &h[i][j])
       }
   }

   // Map board cells to linear positions
   m := 0
   for i := 0; i < 10; i++ {
       if i%2 == 1 {
           for j := 9; j >= 0; j-- {
               b[i][j] = m
               x[m] = i
               y[m] = j
               m++
           }
       } else {
           for j := 0; j < 10; j++ {
               b[i][j] = m
               x[m] = i
               y[m] = j
               m++
           }
       }
   }

   // Compute and output result
   result := solve(99, 1)
   fmt.Fprintf(writer, "%.50f\n", result)
}
