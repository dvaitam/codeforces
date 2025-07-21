package main

import (
   "bufio"
   "fmt"
   "os"
)

func abs(x int) int {
   if x < 0 {
       return -x
   }
   return x
}

func min(a, b int) int {
   if a < b {
       return a
   }
   return b
}

func max(a, b int) int {
   if a > b {
       return a
   }
   return b
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   var n, m, i, j, a, b int
   fmt.Fscan(reader, &n, &m, &i, &j, &a, &b)
   // If already at a corner, no moves needed
   if (i == 1 || i == n) && (j == 1 || j == m) {
       fmt.Fprintln(writer, 0)
       return
   }
   // If any move would go out of bounds immediately, impossible
   if a > n-1 || b > m-1 {
       fmt.Fprintln(writer, "Poor Inna and pony!")
       return
   }
   const INF = int(1e9)
   ans := INF
   // Check all four corners
   corners := [4][2]int{{1, 1}, {1, m}, {n, 1}, {n, m}}
   for _, c := range corners {
       x, y := c[0], c[1]
       dx := abs(x - i)
       dy := abs(y - j)
       if dx%a != 0 || dy%b != 0 {
           continue
       }
       u := dx / a
       v := dy / b
       // Parity of moves must match to synchronize x and y steps
       if (u%2) != (v%2) {
           continue
       }
       k := max(u, v)
       ans = min(ans, k)
   }
   if ans == INF {
       fmt.Fprintln(writer, "Poor Inna and pony!")
   } else {
       fmt.Fprintln(writer, ans)
   }
}
