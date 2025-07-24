package main

import (
   "bufio"
   "fmt"
   "os"
)

const MOD = 1000000007

var n, k int
var g [][]int

func add(a, b int) int {
   a += b
   if a >= MOD {
       a -= MOD
   }
   return a
}

func mul(a, b int) int {
   return int((int64(a) * int64(b)) % MOD)
}

// dpSub[v][d]: number of ways in subtree v with minimal distance from v to painted in subtree = d (0..k+1)
var dpSub [][]int

func dfs(v, p int) {
   dpSub[v] = make([]int, k+2)
   // base dp0: combine children assuming v not painted
   dp0 := make([]int, k+2)
   dp0[k+1] = 1
   prodAll := 1
   for _, u := range g[v] {
       if u == p {
           continue
       }
       dfs(u, v)
       // update dp0
       new0 := make([]int, k+2)
       for d0 := 0; d0 <= k+1; d0++ {
           if dp0[d0] == 0 {
               continue
           }
           for du := 0; du <= k+1; du++ {
               if dpSub[u][du] == 0 {
                   continue
               }
               nd := du + 1
               if nd > k+1 {
                   nd = k+1
               }
               d := d0
               if nd < d {
                   d = nd
               }
               new0[d] = (new0[d] + dp0[d0]*dpSub[u][du]) % MOD
           }
       }
       dp0 = new0
       // total prod
       s := 0
       for du := 0; du <= k+1; du++ {
           s = (s + dpSub[u][du]) % MOD
       }
       prodAll = prodAll * s % MOD
   }
   // v painted
   dpSub[v][0] = prodAll
   // v not painted
   for d := 0; d <= k+1; d++ {
       dpSub[v][d] = add(dpSub[v][d], dp0[d])
   }
}

func main() {
   in := bufio.NewReader(os.Stdin)
   fmt.Fscan(in, &n, &k)
   g = make([][]int, n)
   for i := 0; i < n-1; i++ {
       u, v := 0, 0
       fmt.Fscan(in, &u, &v)
       u--
       v--
       g[u] = append(g[u], v)
       g[v] = append(g[v], u)
   }
   dpSub = make([][]int, n)
   dfs(0, -1)
   res := 0
   for d := 0; d <= k; d++ {
       res = add(res, dpSub[0][d])
   }
   fmt.Println(res)
}
