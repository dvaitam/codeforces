package main

import (
   "bufio"
   "fmt"
   "os"
)

func maxInt64(a, b int64) int64 {
   if a > b {
       return a
   }
   return b
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   var x int64
   fmt.Fscan(reader, &n, &x)

   const INF = int64(4e18)
   var dp0, dp1, dp2, ans int64
   dp1 = -INF
   dp2 = -INF
   ans = 0

   for i := 0; i < n; i++ {
       var a int64
       fmt.Fscan(reader, &a)
       // previous states
       noPrev := dp0
       mulPrev := dp1
       aftPrev := dp2

       // state without multiply
       dp0 = maxInt64(noPrev + a, 0)
       // state inside multiply segment
       dp1 = maxInt64(mulPrev + a*x, noPrev + a*x)
       // state after multiply segment
       dp2 = maxInt64(aftPrev + a, mulPrev + a)

       // update answer
       ans = maxInt64(ans, dp0)
       ans = maxInt64(ans, dp1)
       ans = maxInt64(ans, dp2)
   }

   fmt.Fprintln(writer, ans)
}
