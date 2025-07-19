package main

import (
   "bufio"
   "fmt"
   "os"
)

func max(a, b int) int {
   if a > b {
       return a
   }
   return b
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, x int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   mx1Idx, mx1Val := 0, 0
   mx2Idx, mx2Val := 0, 0
   ans := 0
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &x)
       if mx1Val != 0 && mx1Val != x {
           ans = max(ans, i-mx1Idx)
       }
       if mx2Val != 0 && mx2Val != x {
           ans = max(ans, i-mx2Idx)
       }
       if mx1Val == 0 {
           mx1Val = x
           mx1Idx = i
           continue
       }
       if mx2Val == 0 && x != mx1Val {
           mx2Val = x
           mx2Idx = i
       }
   }
   fmt.Fprintln(writer, ans)
}
