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

   var t int
   if _, err := fmt.Fscan(reader, &t); err != nil {
       return
   }
   for ti := 0; ti < t; ti++ {
       var n int
       fmt.Fscan(reader, &n)
       v := make([]int, n)
       for i := 0; i < n; i++ {
           fmt.Fscan(reader, &v[i])
       }
       cur := 0
       limit := 32
       if n < limit {
           limit = n
       }
       for i := 0; i < limit; i++ {
           best := cur
           bestIdx := -1
           for j := i; j < n; j++ {
               if (cur | v[j]) > best {
                   best = cur | v[j]
                   bestIdx = j
               }
           }
           if bestIdx != -1 {
               v[i], v[bestIdx] = v[bestIdx], v[i]
           }
           cur = best
       }
       // output
       for i, x := range v {
           if i > 0 {
               writer.WriteByte(' ')
           }
           fmt.Fprint(writer, x)
       }
       writer.WriteByte('\n')
   }
}
