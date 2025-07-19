package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

type offer struct {
   b, a, k int
}

func minInt(a, b int) int {
   if a < b {
       return a
   }
   return b
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   A := make([]offer, n)
   for i := 0; i < n; i++ {
       var a, b, k int
       fmt.Fscan(reader, &a, &b, &k)
       A[i] = offer{b: b, a: a, k: k}
   }
   sort.Slice(A, func(i, j int) bool {
       if A[i].b != A[j].b {
           return A[i].b > A[j].b
       }
       if A[i].a != A[j].a {
           return A[i].a > A[j].a
       }
       return A[i].k > A[j].k
   })
   dp := make([][]int64, n+1)
   vis := make([][]bool, n+1)
   for i := 0; i <= n; i++ {
       dp[i] = make([]int64, n+1)
       vis[i] = make([]bool, n+1)
   }
   var f func(i, id int) int64
   f = func(i, id int) int64 {
       if i == n {
           return 0
       }
       if vis[i][id] {
           return dp[i][id]
       }
       vis[i][id] = true
       // skip this offer
       ans := f(i+1, id)
       o := A[i]
       // take with full k penalty
       t1 := int64(o.a) - int64(o.b)*int64(o.k) + f(i+1, id)
       if t1 > ans {
           ans = t1
       }
       // take with limited penalty and increment id
       t2 := int64(o.a) - int64(o.b)*int64(minInt(o.k, id)) + f(i+1, id+1)
       if t2 > ans {
           ans = t2
       }
       dp[i][id] = ans
       return ans
   }
   res := f(0, 0)
   fmt.Println(res)
}
