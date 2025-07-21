package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func abs(x int) int {
   if x < 0 {
       return -x
   }
   return x
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m, k int
   if _, err := fmt.Fscan(reader, &n, &m, &k); err != nil {
       return
   }
   pies := make([][2]int, k)
   for i := 0; i < k; i++ {
       fmt.Fscan(reader, &pies[i][0], &pies[i][1])
   }
   // Compute minimal time to reach each exit edge
   INF := n + m + 5
   T := make([]int, 0, 2*(n+m))
   // top edges (row 0), columns 1..m
   for j := 1; j <= m; j++ {
       best := INF
       for _, pie := range pies {
           x, y := pie[0], pie[1]
           d := x + abs(y-j)
           if d < best {
               best = d
           }
       }
       T = append(T, best)
   }
   // bottom edges (row n+1)
   for j := 1; j <= m; j++ {
       best := INF
       for _, pie := range pies {
           x, y := pie[0], pie[1]
           d := (n-x+1) + abs(y-j)
           if d < best {
               best = d
           }
       }
       T = append(T, best)
   }
   // left edges (col 0), rows 1..n
   for i := 1; i <= n; i++ {
       best := INF
       for _, pie := range pies {
           x, y := pie[0], pie[1]
           d := y + abs(x-i)
           if d < best {
               best = d
           }
       }
       T = append(T, best)
   }
   // right edges (col m+1)
   for i := 1; i <= n; i++ {
       best := INF
       for _, pie := range pies {
           x, y := pie[0], pie[1]
           d := (m-y+1) + abs(x-i)
           if d < best {
               best = d
           }
       }
       T = append(T, best)
   }
   sort.Ints(T)
   // Check if exists exit that Volodya can reach before Vlad bans it
   ans := "NO"
   for idx, t := range T {
       // at time t, Volodya moves; Vlad has done (t-1) bans before
       // ban of this exit would occur at time (idx+1)
       if t <= idx+1 {
           ans = "YES"
           break
       }
   }
   fmt.Println(ans)
}
