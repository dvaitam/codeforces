package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   quals := make([]int, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &quals[i])
   }
   var m int
   fmt.Fscan(reader, &m)
   const INF = int(1e18)
   best := make([]int, n+1)
   // initialize best costs to INF
   for i := 1; i <= n; i++ {
       best[i] = int(INF)
   }
   for j := 0; j < m; j++ {
       var a, b int
       var c int
       fmt.Fscan(reader, &a, &b, &c)
       // a can supervise b
       if c < best[b] {
           best[b] = c
       }
   }
   // find root: unique max qualification
   maxq := -1
   root := -1
   cntMax := 0
   for i := 1; i <= n; i++ {
       if quals[i] > maxq {
           maxq = quals[i]
           root = i
           cntMax = 1
       } else if quals[i] == maxq {
           cntMax++
       }
   }
   if cntMax != 1 {
       fmt.Println(-1)
       return
   }
   // sum minimal incoming costs for all except root
   var sum int64 = 0
   for i := 1; i <= n; i++ {
       if i == root {
           continue
       }
       if best[i] == int(INF) {
           fmt.Println(-1)
           return
       }
       sum += int64(best[i])
   }
   fmt.Println(sum)
}
