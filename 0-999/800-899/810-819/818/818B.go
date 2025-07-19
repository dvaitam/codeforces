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
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   a := make([]int, m)
   for i := 0; i < m; i++ {
       fmt.Fscan(reader, &a[i])
   }
   // mask[i]: assigned shift for position i (1-based)
   mask := make([]int, n+1)
   used := make([]bool, n+1)
   ok := true
   for i := 0; i+1 < m; i++ {
       cur, next := a[i], a[i+1]
       dd := next - cur
       if dd <= 0 {
           dd += n
       }
       if mask[cur] == 0 {
           mask[cur] = dd
       } else if mask[cur] != dd {
           ok = false
       }
   }
   if !ok {
       fmt.Fprint(writer, -1)
       return
   }
   // Check for duplicate shifts
   for i := 1; i <= n; i++ {
       if mask[i] > 0 {
           if used[mask[i]] {
               fmt.Fprint(writer, -1)
               return
           }
           used[mask[i]] = true
       }
   }
   // Collect unused shifts
   tmp := make([]int, 0, n)
   for i := 1; i <= n; i++ {
       if !used[i] {
           tmp = append(tmp, i)
       }
   }
   // Assign remaining shifts
   for i := 1; i <= n; i++ {
       if mask[i] == 0 {
           last := tmp[len(tmp)-1]
           tmp = tmp[:len(tmp)-1]
           mask[i] = last
       }
   }
   // Output result
   for i := 1; i <= n; i++ {
       if i > 1 {
           fmt.Fprint(writer, " ")
       }
       fmt.Fprint(writer, mask[i])
   }
}
