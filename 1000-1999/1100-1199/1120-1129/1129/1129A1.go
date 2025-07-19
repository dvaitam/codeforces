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
   noc := make([]int, n)
   bd := make([]int, n)
   for i := 0; i < n; i++ {
       bd[i] = n
   }
   var mc int
   for i := 0; i < m; i++ {
       var a, b int
       fmt.Fscan(reader, &a, &b)
       a--
       b--
       noc[a]++
       dst := b - a
       if dst < 0 {
           dst += n
       }
       if dst < bd[a] {
           bd[a] = dst
       }
       if noc[a] > mc {
           mc = noc[a]
       }
   }
   vas := (mc - 1) * n
   lv := 0
   for i := 0; i < n; i++ {
       if noc[i] == mc {
           if i+bd[i] > lv {
               lv = i + bd[i]
           }
       }
       if mc > 1 && noc[i] == mc-1 {
           if bd[i]+i-n > lv {
               lv = bd[i] + i - n
           }
       }
   }
   // output sequence
   fmt.Fprint(writer, vas+lv)
   for i := 1; i < n; i++ {
       fmt.Fprint(writer, " ")
       lv--
       j := i - 1
       if noc[j] == mc {
           if n-1+bd[j] > lv {
               lv = n - 1 + bd[j]
           }
       }
       if mc > 1 && noc[j] == mc-1 {
           if bd[j]-1 > lv {
               lv = bd[j] - 1
           }
       }
       fmt.Fprint(writer, vas+lv)
   }
   fmt.Fprintln(writer)
}
