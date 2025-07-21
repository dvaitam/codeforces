package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, q int
   if _, err := fmt.Fscan(reader, &n, &q); err != nil {
       return
   }
   v := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &v[i])
   }
   c := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &c[i])
   }
   // dp per color with lazy reset via versioning
   dp := make([]int64, n+1)
   ver := make([]int, n+1)
   curVer := 0
   const minf = int64(-1e18)

   for qi := 0; qi < q; qi++ {
       var a, b int64
       fmt.Fscan(reader, &a, &b)
       curVer++
       // best1 and best2: (color, value)
       best1col, best2col := 0, 0
       best1val, best2val := int64(0), minf

       for i := 0; i < n; i++ {
           col := c[i]
           vi := v[i]
           var dpi int64
           if ver[col] == curVer {
               dpi = dp[col]
           } else {
               dpi = minf
           }
           // extend same color
           same := dpi + vi*a
           // extend from best different color (or empty)
           var diffBase int64
           if best1col != col {
               diffBase = best1val
           } else {
               diffBase = best2val
           }
           diff := diffBase + vi*b
           // choose best for this color
           mx := same
           if diff > mx {
               mx = diff
           }
           // update dp[col]
           if ver[col] != curVer {
               ver[col] = curVer
               dp[col] = mx
           } else if mx > dp[col] {
               dp[col] = mx
           }
           // update best1 and best2
           if dp[col] > best1val {
               if col != best1col {
                   best2col, best2val = best1col, best1val
               }
               best1col, best1val = col, dp[col]
           } else if col != best1col && dp[col] > best2val {
               best2col, best2val = col, dp[col]
           }
       }
       fmt.Fprintln(writer, best1val)
   }
}
