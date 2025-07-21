package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, m int
   var x int64
   fmt.Fscan(in, &n, &m)
   fmt.Fscan(in, &x)
   // Each initially black square at distance k from the border (k = min(r-1,c-1,n-r,m-c))
   // is painted exactly k+1 times. We need count of squares with k+1 == x, i.e., k = x-1.
   k := x - 1
   b0 := countBlack(n, m, k)
   b1 := countBlack(n, m, k+1)
   ans := b0 - b1
   fmt.Println(ans)
}

// countBlack returns the number of initially black squares in the sub-rectangle
// after removing k layers from each border (i.e., size (n-2k) x (m-2k)).
// Top-left of the board is black, and coloring alternates, so blacks = ceil(area/2).
func countBlack(n, m int, k int64) int64 {
   // Remaining dimensions
   nn := int64(n) - 2*k
   mm := int64(m) - 2*k
   if nn <= 0 || mm <= 0 {
       return 0
   }
   area := nn * mm
   return (area + 1) / 2
}
