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

   var n, m int
   fmt.Fscan(reader, &n, &m)
   price := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &price[i])
   }
   con := make([][]bool, n)
   for i := range con {
       con[i] = make([]bool, n)
   }
   for i := 0; i < m; i++ {
       var u, v int
       fmt.Fscan(reader, &u, &v)
       u--
       v--
       con[u][v] = true
       con[v][u] = true
   }
   const INF = 1 << 60
   best := INF
   for i := 0; i < n; i++ {
       for k := i + 1; k < n; k++ {
           if !con[i][k] {
               continue
           }
           for j := k + 1; j < n; j++ {
               if con[i][j] && con[k][j] {
                   sum := price[i] + price[k] + price[j]
                   if sum < best {
                       best = sum
                   }
               }
           }
       }
   }
   if best == INF {
       fmt.Fprint(writer, -1)
   } else {
       fmt.Fprint(writer, best)
   }
}
