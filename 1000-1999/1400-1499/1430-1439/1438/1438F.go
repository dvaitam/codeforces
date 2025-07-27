package main

import (
   "bufio"
   "fmt"
   "math/rand"
   "os"
   "time"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var h int
   if _, err := fmt.Fscan(reader, &h); err != nil {
       return
   }
   n := (1 << h) - 1
   rand.Seed(time.Now().UnixNano())
   counts := make([]int, n+1)

   // Sample M = n pairs to estimate LCA frequencies
   M := n
   for i := 0; i < M; i++ {
       // pick two distinct nodes u, v
       u := rand.Intn(n) + 1
       v := rand.Intn(n) + 1
       if u == v {
           i--
           continue
       }
       // query until w is not on path(u,v)
       for {
           w := rand.Intn(n) + 1
           if w == u || w == v {
               continue
           }
           fmt.Fprintf(writer, "? %d %d %d\n", u, v, w)
           writer.Flush()
           var res int
           if _, err := fmt.Fscan(reader, &res); err != nil {
               os.Exit(0)
           }
           if res == w {
               // w on path, retry
               continue
           }
           counts[res]++
           break
       }
   }
   // find label with max count
   best := 1
   for x := 2; x <= n; x++ {
       if counts[x] > counts[best] {
           best = x
       }
   }
   // output answer
   fmt.Fprintf(writer, "! %d\n", best)
   writer.Flush()
}
