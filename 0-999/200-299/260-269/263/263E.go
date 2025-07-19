package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m, k int
   if _, err := fmt.Fscan(reader, &n, &m, &k); err != nil {
       return
   }
   size := n + m + 5
   // rotated grid and prefix sums
   a := make([][]int64, size)
   for i := range a {
       a[i] = make([]int64, size)
   }
   // read and map to rotated coordinates
   for i := 1; i <= n; i++ {
       for j := 1; j <= m; j++ {
           var v int64
           fmt.Fscan(reader, &v)
           x := i + j
           y := i - j + m
           a[x][y] = v
       }
   }
   // build 2D prefix sums
   limit := n + m
   for i := 1; i <= limit; i++ {
       row := a[i]
       prevRow := a[i-1]
       for j := 1; j <= limit; j++ {
           row[j] += prevRow[j] + row[j-1] - prevRow[j-1]
       }
   }
   // helper for sum on prefix
   sum := func(l, L, r, R int) int64 {
       return a[r][R] - a[l-1][R] - a[r][L-1] + a[l-1][L-1]
   }
   var best int64 = -1
   var ansX, ansY int
   // search best center
   for i := k; i <= n-k+1; i++ {
       for j := k; j <= m-k+1; j++ {
           x := i + j
           y := i - j + m
           var s int64
           for z := 0; z < k; z++ {
               s += sum(x-z, y-z, x+z, y+z)
           }
           if s > best {
               best = s
               ansX, ansY = i, j
           }
       }
   }
   // output result
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fprintf(writer, "%d %d", ansX, ansY)
}
