package main

import (
   "bufio"
   "fmt"
   "math/big"
   "os"
)

var (
   n int
   adj [][]int
)

// dfs returns dp and size of subtree rooted at u (excluding parent p)
// dp[k] = maximum product of components in subtree with u-component size k (product excludes the final multiplication by k)
func dfs(u, p int) ([]*big.Int, int) {
   // initialize dp for u: only size 1 with product 1
   dp := make([]*big.Int, 2)
   dp[1] = big.NewInt(1)
   sz := 1
   for _, v := range adj[u] {
       if v == p {
           continue
       }
       dpv, sv := dfs(v, u)
       // new dp size up to sz+sv
       newDp := make([]*big.Int, sz+sv+1)
       for su := 1; su <= sz; su++ {
           if dp[su] == nil {
               continue
           }
           for svv := 1; svv < len(dpv); svv++ {
               if dpv[svv] == nil {
                   continue
               }
               // case: not cut edge, combine u-component sizes
               nk := su + svv
               // prod = dp[su] * dpv[svv]
               prod1 := new(big.Int).Mul(dp[su], dpv[svv])
               if newDp[nk] == nil || prod1.Cmp(newDp[nk]) > 0 {
                   newDp[nk] = prod1
               }
               // case: cut edge, component v becomes its own
               // u-component size stays su, product includes dpv * svv
               prod2 := new(big.Int).Mul(dp[su], dpv[svv])
               prod2.Mul(prod2, big.NewInt(int64(svv)))
               if newDp[su] == nil || prod2.Cmp(newDp[su]) > 0 {
                   newDp[su] = prod2
               }
           }
       }
       sz += sv
       dp = newDp
   }
   return dp, sz
}

func main() {
   in := bufio.NewReader(os.Stdin)
   fmt.Fscan(in, &n)
   adj = make([][]int, n+1)
   for i := 0; i < n-1; i++ {
       var u, v int
       fmt.Fscan(in, &u, &v)
       adj[u] = append(adj[u], v)
       adj[v] = append(adj[v], u)
   }
   dp, _ := dfs(1, 0)
   // find max of dp[k] * k
   ans := big.NewInt(0)
   for k, prod := range dp {
       if prod == nil {
           continue
       }
       // total product = prod * k
       total := new(big.Int).Mul(prod, big.NewInt(int64(k)))
       if total.Cmp(ans) > 0 {
           ans = total
       }
   }
   fmt.Println(ans)
}
