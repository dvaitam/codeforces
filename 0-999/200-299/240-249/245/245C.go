package main

import (
   "bufio"
   "fmt"
   "os"
)

func min64(a, b int64) int64 {
   if a < b {
       return a
   }
   return b
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   a := make([]int, n+1)
   maxA := 0
   for i := 1; i <= n; i++ {
       fmt.Fscan(in, &a[i])
       if a[i] > maxA {
           maxA = a[i]
       }
   }
   // maximum possible incoming parent ops to consider
   MAX := maxA
   // number of valid x variables: x such that 2*x+1 <= n
   m := (n - 1) / 2
   const INF = int64(1e18)
   // f[v][p] = minimal cost in subtree v given parent contribution p
   f := make([][]int64, n+2)
   // process nodes bottom-up
   for v := n; v >= 1; v-- {
       // children
       c1 := 2 * v
       c2 := 2*v + 1
       var f1, f2 []int64
       if c1 <= n {
           f1 = f[c1]
       }
       if c2 <= n {
           f2 = f[c2]
       }
       // build s[k] = cost if we choose k operations at v
       hasVar := v <= m
       s := make([]int64, MAX+1)
       for k := 0; k <= MAX; k++ {
           if !hasVar && k > 0 {
               s[k] = INF
               continue
           }
           var cost int64 = int64(k)
           if f1 != nil {
               cost += f1[k]
           }
           if f2 != nil {
               cost += f2[k]
           }
           s[k] = cost
       }
       // suffix minimum of s
       for k := MAX - 1; k >= 0; k-- {
           if s[k+1] < s[k] {
               s[k] = s[k+1]
           }
       }
       // compute f[v]
       fv := make([]int64, MAX+1)
       for p := 0; p <= MAX; p++ {
           lb := a[v] - p
           if lb < 0 {
               lb = 0
           }
           if lb > MAX {
               // impossible
               fv[p] = INF
           } else {
               fv[p] = s[lb]
           }
       }
       f[v] = fv
   }
   res := f[1][0]
   if res >= int64(1e17) {
       fmt.Println(-1)
   } else {
       fmt.Println(res)
   }
}
