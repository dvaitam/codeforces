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

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   // pos[i][0]: first occurrence, pos[i][1]: second occurrence
   pos := make([][2]int, n+1)
   // initial previous positions for 0
   pos[0][0], pos[0][1] = 1, 1

   total := 2 * n
   for i := 1; i <= total; i++ {
       var x int
       fmt.Fscan(reader, &x)
       if pos[x][0] != 0 {
           pos[x][1] = i
       } else {
           pos[x][0] = i
       }
   }

   var ans int64
   for i := 1; i <= n; i++ {
       p0, p1 := pos[i][0], pos[i][1]
       q0, q1 := pos[i-1][0], pos[i-1][1]
       cost1 := abs(p1-q1) + abs(p0-q0)
       cost2 := abs(p0-q1) + abs(p1-q0)
       ans += int64(min(cost1, cost2))
   }
   fmt.Fprint(writer, ans)
}
