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
   adj := make([][]int, n)
   for i := 0; i < n-1; i++ {
       var u, v int
       fmt.Fscan(in, &u, &v)
       u--; v--
       adj[u] = append(adj[u], v)
       adj[v] = append(adj[v], u)
   }
   a := make([]int64, n)
   for i := 0; i < n; i++ {
       var x int64
       fmt.Fscan(in, &x)
       a[i] = x
   }
   // build parent and order
   par := make([]int, n)
   order := make([]int, 0, n)
   stack := make([]int, 0, n)
   stack = append(stack, 0)
   par[0] = -1
   for i := 0; i < len(stack); i++ {
       u := stack[i]
       order = append(order, u)
       for _, v := range adj[u] {
           if v == par[u] {
               continue
           }
           par[v] = u
           stack = append(stack, v)
       }
   }
   // dp1: best down path
   dp1 := make([]int64, n)
   sumA1 := make([]int64, n)
   // postorder: reverse order
   for i := n-1; i >= 0; i-- {
       u := order[i]
       dp1[u] = a[u]
       sumA1[u] = a[u]
       var bestDir int64 = 0
       var bestSumA int64 = 0
       for _, v := range adj[u] {
           if v == par[u] {
               continue
           }
           // direction to child v
           dir := dp1[v] + sumA1[v]
           if dir > bestDir {
               bestDir = dir
               bestSumA = sumA1[v]
           }
       }
       if bestDir > 0 {
           dp1[u] = a[u] + bestDir
           sumA1[u] = a[u] + bestSumA
       }
   }
   // dp2: best up path
   dp2 := make([]int64, n)
   sumA2 := make([]int64, n)
   dp2[0], sumA2[0] = 0, 0
   const inf = int64(4e18)
   // reroot in order
   for _, u := range order {
       // collect best two directions by val = dp_dir + sumA_dir
       var max1Val, max2Val int64 = -inf, -inf
       var max1SumA, max2SumA int64
       var max1Id int = -1
       // parent direction
       if par[u] != -1 {
           d := dp2[u]
           s := sumA2[u]
           val := d + s
           if val > max1Val {
               max2Val, max2SumA = max1Val, max1SumA
               max1Val, max1SumA, max1Id = val, s, par[u]
           } else if val > max2Val {
               max2Val, max2SumA = val, s
           }
       }
       // children directions
       for _, v := range adj[u] {
           if v == par[u] {
               continue
           }
           s := sumA1[v]
           d := dp1[v] + s
           val := d + s
           if val > max1Val {
               max2Val, max2SumA = max1Val, max1SumA
               max1Val, max1SumA, max1Id = val, s, v
           } else if val > max2Val {
               max2Val, max2SumA = val, s
           }
       }
       // propagate to children
       for _, v := range adj[u] {
           if v == par[u] {
               continue
           }
           var bestVal, bestSumA int64
           if max1Id != v {
               bestVal = max1Val
               bestSumA = max1SumA
           } else {
               bestVal = max2Val
               bestSumA = max2SumA
           }
           if bestVal > 0 {
               dp2[v] = 2*a[u] + bestVal
               sumA2[v] = a[u] + bestSumA
           } else {
               dp2[v] = 0
               sumA2[v] = 0
           }
       }
   }
   // answer
   var ans int64 = 0
   for u := 0; u < n; u++ {
       // f[u]
       f1 := dp1[u]
       f2 := a[u] + dp2[u]
       if f1 > ans {
           ans = f1
       }
       if f2 > ans {
           ans = f2
       }
   }
   // output
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   fmt.Fprint(out, ans)
}
