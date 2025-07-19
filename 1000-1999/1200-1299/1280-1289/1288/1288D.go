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
   a := make([][]int, n)
   for i := 0; i < n; i++ {
       a[i] = make([]int, m)
       for j := 0; j < m; j++ {
           fmt.Fscan(reader, &a[i][j])
       }
   }

   fullMask := (1 << m) - 1
   mrk := make([]int, 1<<m)
   var ans1, ans2 int

   check := func(x int) bool {
       // reset markers
       for i := range mrk {
           mrk[i] = 0
       }
       // mark rows by mask
       for i := 0; i < n; i++ {
           mask := 0
           row := a[i]
           for j := 0; j < m; j++ {
               if row[j] >= x {
                   mask |= 1 << j
               }
           }
           mrk[mask] = i + 1
       }
       // find two masks covering all bits
       for i := 0; i < len(mrk); i++ {
           if mrk[i] == 0 {
               continue
           }
           for j := 0; j < len(mrk); j++ {
               if mrk[j] == 0 {
                   continue
               }
               if (i | j) == fullMask {
                   ans1 = mrk[i]
                   ans2 = mrk[j]
                   return true
               }
           }
       }
       return false
   }

   l, r := 0, 1000000003
   for r-l > 1 {
       mid := (l + r) / 2
       if check(mid) {
           l = mid
       } else {
           r = mid
       }
   }
   check(l)
   fmt.Fprintln(writer, ans1, ans2)
}
