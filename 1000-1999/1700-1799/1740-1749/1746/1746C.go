package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   var t int
   fmt.Fscan(in, &t)
   for t > 0 {
       t--
       var n int
       fmt.Fscan(in, &n)
       a := make([]int, n)
       for i := range a {
           fmt.Fscan(in, &a[i])
       }
       type pair struct{ diff, idx int }
       b := make([]pair, max(0, n-1))
       for i := 0; i < n-1; i++ {
           b[i] = pair{a[i] - a[i+1], i + 2}
       }
       sort.Slice(b, func(i, j int) bool {
           return b[i].diff < b[j].diff
       })
       fmt.Fprint(out, "1")
       for _, p := range b {
           fmt.Fprint(out, " ", p.idx)
       }
       fmt.Fprintln(out)
   }
}

func max(a, b int) int {
   if a > b {
       return a
   }
   return b
}
