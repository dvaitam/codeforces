package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   adj := make([][]int, n+1)
   for i := 0; i < m; i++ {
       var a, b int
       fmt.Fscan(reader, &a, &b)
       adj[a] = append(adj[a], b)
   }
   cnt := make([]int, n+1)
   var ans int64
   for a := 1; a <= n; a++ {
       // count number of distinct length-2 paths from a to each c
       for _, b := range adj[a] {
           for _, c := range adj[b] {
               cnt[c]++
           }
       }
       // for each c, add choose(cnt[c], 2)
       for c := 1; c <= n; c++ {
           k := cnt[c]
           if k > 1 {
               ans += int64(k) * int64(k-1) / 2
           }
           cnt[c] = 0
       }
   }
   fmt.Println(ans)
}
