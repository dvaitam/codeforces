package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var t int
   if _, err := fmt.Fscan(reader, &t); err != nil {
       return
   }
   for t > 0 {
       t--
       var n, k int
       fmt.Fscan(reader, &n, &k)
       a := make([]pair, n)
       for i := 0; i < n; i++ {
           var v int
           fmt.Fscan(reader, &v)
           a[i] = pair{v: v, idx: i + 1}
       }
       sort.Slice(a, func(i, j int) bool { return a[i].v < a[j].v })
       // initial minxor larger than any possible xor (max xor < 1<<k)
       minxor := 1 << k
       var bestI, bestJ, v1 int
       for i := 1; i < n; i++ {
           x := a[i].v ^ a[i-1].v
           if x < minxor {
               minxor = x
               v1 = a[i].v
               bestI = a[i].idx
               bestJ = a[i-1].idx
           }
       }
       mask := (1<<k) - 1
       x := mask ^ v1
       fmt.Fprintf(writer, "%d %d %d\n", bestI, bestJ, x)
   }
}

type pair struct {
   v   int
   idx int
}
