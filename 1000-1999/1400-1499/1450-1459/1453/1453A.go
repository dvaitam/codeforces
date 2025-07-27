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
   fmt.Fscan(reader, &t)
   for tc := 0; tc < t; tc++ {
       var n, m int
       fmt.Fscan(reader, &n, &m)
       seen := make(map[int]bool, n)
       for i := 0; i < n; i++ {
           var x int
           fmt.Fscan(reader, &x)
           seen[x] = true
       }
       cancel := 0
       for i := 0; i < m; i++ {
           var x int
           fmt.Fscan(reader, &x)
           if seen[x] {
               cancel++
           }
       }
       fmt.Fprintln(writer, cancel)
   }
}
