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

   var T int
   fmt.Fscan(reader, &T)
   for ; T > 0; T-- {
       var n int
       fmt.Fscan(reader, &n)
       pos := make([]int, 0, n)
       neg := make([]int, 0, n)
       for i := 0; i < n; i++ {
           var x int
           fmt.Fscan(reader, &x)
           if x >= 0 {
               pos = append(pos, x)
           } else {
               neg = append(neg, x)
           }
       }
       sort.Slice(pos, func(i, j int) bool { return pos[i] > pos[j] })
       sort.Ints(neg)
       if len(neg) == 0 && pos[0] == 0 {
           fmt.Fprintln(writer, "NO")
           continue
       }
       fmt.Fprintln(writer, "YES")
       var now int64
       p0, p1 := 0, 0
       for i := 0; i < n; i++ {
           if now <= 0 {
               v := pos[p0]
               now += int64(v)
               fmt.Fprint(writer, v)
               p0++
           } else {
               v := neg[p1]
               now += int64(v)
               fmt.Fprint(writer, v)
               p1++
           }
           if i < n-1 {
               fmt.Fprint(writer, " ")
           }
       }
       fmt.Fprintln(writer)
   }
}
