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
   adj := make([][]int, n+1)
   for i := 0; i < n-1; i++ {
       var u, v int
       fmt.Fscan(in, &u, &v)
       adj[u] = append(adj[u], v)
       adj[v] = append(adj[v], u)
   }
   vals := make([]int64, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(in, &vals[i])
   }
   var dfs func(u, p int) (inc, dec int64)
   dfs = func(u, p int) (int64, int64) {
       var inc, dec int64
       for _, v := range adj[u] {
           if v == p {
               continue
           }
           ci, cd := dfs(v, u)
           if ci > inc {
               inc = ci
           }
           if cd > dec {
               dec = cd
           }
       }
       // adjust current value
       cur := vals[u] + inc - dec
       if cur > 0 {
           // need more negative operations
           dec += cur
       } else {
           // need more positive operations
           inc += -cur
       }
       return inc, dec
   }
   inc, dec := dfs(1, 0)
   // result is sum of operations
   fmt.Println(inc + dec)
}
