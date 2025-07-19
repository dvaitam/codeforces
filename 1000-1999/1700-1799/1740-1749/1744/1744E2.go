package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   var t int
   if _, err := fmt.Fscan(in, &t); err != nil {
       return
   }
   for ; t > 0; t-- {
       var a, b, c, d int64
       fmt.Fscan(in, &a, &b, &c, &d)
       // factor counts
       cnt := make(map[int64]int)
       // factor a
       x := a
       for i := int64(2); i*i <= x; i++ {
           for x%i == 0 {
               cnt[i]++
               x /= i
           }
       }
       if x > 1 {
           cnt[x]++
       }
       // factor b
       x = b
       for i := int64(2); i*i <= x; i++ {
           for x%i == 0 {
               cnt[i]++
               x /= i
           }
       }
       if x > 1 {
           cnt[x]++
       }
       // build prime-exp pairs
       type pe struct{p, e int64}
       var q []pe
       for p, e := range cnt {
           q = append(q, pe{p, int64(e)})
       }
       // generate divisors res <= c and e/res <= d
       eProd := a * b
       var v []int64
       var dfs func(int, int64)
       dfs = func(u int, res int64) {
           if res > c {
               return
           }
           if u == len(q) {
               if eProd/res <= d {
                   v = append(v, res)
               }
               return
           }
           p := q[u].p
           exp := q[u].e
           cur := res
           for i := int64(0); i <= exp; i++ {
               dfs(u+1, cur)
               cur *= p
           }
       }
       dfs(0, 1)
       // find answer
       ansX, ansY := int64(-1), int64(-1)
       for _, dv := range v {
           y := eProd / dv
           xmul := ((a/dv) + 1) * dv
           ymul := ((b/y) + 1) * y
           if xmul <= c && ymul <= d {
               ansX, ansY = xmul, ymul
               break
           }
       }
       fmt.Fprintln(out, ansX, ansY)
   }
}
