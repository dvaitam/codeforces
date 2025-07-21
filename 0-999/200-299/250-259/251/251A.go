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

   var n int
   var d int64
   fmt.Fscan(in, &n, &d)
   x := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &x[i])
   }
   var res int64
   r := 0
   for l := 0; l < n; l++ {
       if r < l {
           r = l
       }
       for r+1 < n && x[r+1]-x[l] <= d {
           r++
       }
       // number of points between l and r is r-l
       cnt := int64(r - l)
       if cnt >= 2 {
           // choose cnt points two among them: C(cnt,2)
           res += cnt * (cnt - 1) / 2
       }
   }
   fmt.Fprintln(out, res)
}
