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
   var T int
   fmt.Fscan(reader, &T)
   for T > 0 {
       T--
       var s, t string
       fmt.Fscan(reader, &s, &t)
       n := len(s)
       m := len(t)
       ok := false
       // try splitting t into prefix [0:i) and suffix [i:m)
       for i := 0; i <= m; i++ {
           used := make([]bool, n)
           // match prefix of length i
           p := 0
           for j := 0; j < n && p < i; j++ {
               if s[j] == t[p] {
                   used[j] = true
                   p++
               }
           }
           if p < i {
               continue
           }
           // match suffix
           q := i
           for j := 0; j < n && q < m; j++ {
               if used[j] {
                   continue
               }
               if s[j] == t[q] {
                   q++
               }
           }
           if q == m {
               ok = true
               break
           }
       }
       if ok {
           fmt.Fprintln(writer, "YES")
       } else {
           fmt.Fprintln(writer, "NO")
       }
   }
}
