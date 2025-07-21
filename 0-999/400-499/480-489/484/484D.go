package main

import (
   "bufio"
   "fmt"
   "os"
)

// node stores a value for monotonic stack processing
type node struct {
   val  int64 // the charisma value
   idx  int   // the earliest index this value applies to
   best int64 // dp value prior to this segment
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   fmt.Fscan(reader, &n)
   a := make([]int64, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &a[i])
   }

   dp := make([]int64, n+1)
   stMax := make([]node, 0, n)
   stMin := make([]node, 0, n)

   for i := 1; i <= n; i++ {
       x := a[i]
       // start with no new group: extend previous grouping
       dpi := dp[i-1]

       // process max-stack: increasing stack for maxima
       pre := i
       for len(stMax) > 0 && stMax[len(stMax)-1].val <= x {
           top := stMax[len(stMax)-1]
           stMax = stMax[:len(stMax)-1]
           // consider segment ending here with this max
           if top.best + x - top.val > dpi {
               dpi = top.best + x - top.val
           }
           pre = top.idx
       }
       // push new candidate for max
       stMax = append(stMax, node{val: x, idx: pre, best: dp[pre-1]})

       // process min-stack: decreasing stack for minima
       pre2 := i
       for len(stMin) > 0 && stMin[len(stMin)-1].val >= x {
           top := stMin[len(stMin)-1]
           stMin = stMin[:len(stMin)-1]
           if top.best + top.val - x > dpi {
               dpi = top.best + top.val - x
           }
           pre2 = top.idx
       }
       // push new candidate for min
       stMin = append(stMin, node{val: x, idx: pre2, best: dp[pre2-1]})

       dp[i] = dpi
   }

   fmt.Fprint(writer, dp[n])
}
