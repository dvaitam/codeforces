package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   t := make([]int64, n+1)
   h := make([]int, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(in, &t[i])
   }
   for i := 1; i <= n; i++ {
       fmt.Fscan(in, &h[i])
   }
   adj := make([][]int, n+1)
   for i := 0; i < n-1; i++ {
       var u, v int
       fmt.Fscan(in, &u, &v)
       adj[u] = append(adj[u], v)
       adj[v] = append(adj[v], u)
   }
   var ans int64
   for i := 1; i <= n; i++ {
       var l, r, s int
       for _, j := range adj[i] {
           if h[j] < h[i] {
               l++
           } else if h[j] > h[i] {
               r++
           } else {
               s++
           }
       }
       d := l + r + s
       var pairs int
       if abs(l-r) <= s {
           // can pair to maximize, pairs = floor(d/2)
           pairs = d / 2
       } else {
           // limited by smaller side plus s
           if l < r {
               pairs = l + s
           } else {
               pairs = r + s
           }
       }
       cnt := d - pairs
       // each node appears cnt times in total challenges
       ans += int64(cnt) * t[i]
   }
   fmt.Println(ans)
}

func abs(x int) int {
   if x < 0 {
       return -x
   }
   return x
}
