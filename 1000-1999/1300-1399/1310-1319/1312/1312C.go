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
   for tc := 0; tc < t; tc++ {
       var n int
       var k uint64
       fmt.Fscan(reader, &n, &k)
       a := make([]uint64, n)
       for i := 0; i < n; i++ {
           fmt.Fscan(reader, &a[i])
       }
       used := make([]bool, 64)
       ok := true
       for i := 0; i < n && ok; i++ {
           x := a[i]
           pos := 0
           for x > 0 {
               rem := x % k
               if rem > 1 {
                   ok = false
                   break
               }
               if rem == 1 {
                   if pos >= len(used) || used[pos] {
                       ok = false
                       break
                   }
                   used[pos] = true
               }
               x /= k
               pos++
           }
       }
       if ok {
           fmt.Fprintln(writer, "YES")
       } else {
           fmt.Fprintln(writer, "NO")
       }
   }
}
