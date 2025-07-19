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

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   deg := make([]int, n)
   xorneigh := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &deg[i], &xorneigh[i])
   }
   sum := 0
   for _, d := range deg {
       sum += d
   }
   fmt.Fprintln(writer, sum/2)

   queue := make([]int, 0, n)
   for i, d := range deg {
       if d == 1 {
           queue = append(queue, i)
       }
   }
   for head := 0; head < len(queue); head++ {
       u := queue[head]
       if deg[u] != 1 {
           continue
       }
       v := xorneigh[u]
       fmt.Fprintln(writer, u, v)
       deg[u]--
       xorneigh[u] ^= v
       deg[v]--
       xorneigh[v] ^= u
       if deg[v] == 1 {
           queue = append(queue, v)
       }
   }
}
