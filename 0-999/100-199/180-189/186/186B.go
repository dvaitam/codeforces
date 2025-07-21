package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, t1, t2, k int
   if _, err := fmt.Fscan(reader, &n, &t1, &t2, &k); err != nil {
       return
   }
   type result struct {
       idx int
       h   int
   }
   res := make([]result, 0, n)
   for i := 1; i <= n; i++ {
       var a, b int
       fmt.Fscan(reader, &a, &b)
       // Compute height*100 to avoid floating issues
       // option1: a before, b after
       h1 := a*t1*(100-k) + b*t2*100
       // option2: b before, a after
       h2 := b*t1*(100-k) + a*t2*100
       h := h1
       if h2 > h1 {
           h = h2
       }
       res = append(res, result{idx: i, h: h})
   }
   // Sort by descending height, then ascending index
   sort.Slice(res, func(i, j int) bool {
       if res[i].h != res[j].h {
           return res[i].h > res[j].h
       }
       return res[i].idx < res[j].idx
   })
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   for _, r := range res {
       // print height with two decimals
       whole := r.h / 100
       frac := r.h % 100
       fmt.Fprintf(writer, "%d %d.%02d\n", r.idx, whole, frac)
   }
}
