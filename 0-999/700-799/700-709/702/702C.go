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

   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   cities := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &cities[i])
   }
   towers := make([]int, m)
   for i := 0; i < m; i++ {
       fmt.Fscan(reader, &towers[i])
   }

   var result int64
   for _, x := range cities {
       idx := sort.Search(m, func(i int) bool { return towers[i] >= x })
       // compute minimal distance to a tower
       var dist int64 = 1<<62 - 1
       if idx < m {
           d := int64(towers[idx] - x)
           if d < dist {
               dist = d
           }
       }
       if idx > 0 {
           d := int64(x - towers[idx-1])
           if d < dist {
               dist = d
           }
       }
       if dist > result {
           result = dist
       }
   }
   fmt.Fprintln(writer, result)
}
