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

   var n, k, m, t int
   if _, err := fmt.Fscan(reader, &n, &k, &m, &t); err != nil {
       return
   }
   l := n
   for i := 0; i < t; i++ {
       var typ, pos int
       fmt.Fscan(reader, &typ, &pos)
       if typ == 1 {
           // insert at position pos
           if pos <= k {
               k++
           }
           l++
       } else {
           // break link at position pos between pos and pos+1
           if k <= pos {
               // doctor in left segment
               l = pos
           } else {
               // doctor in right segment
               l = l - pos
               k = k - pos
           }
       }
       fmt.Fprintln(writer, l, k)
   }
