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

   var n int
   fmt.Fscan(reader, &n)
   cnt := make([]int, 100001)
   var x int
   P, Q := 0, 0
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &x)
       c := cnt[x]
       P -= c / 2
       Q -= c / 4
       c++
       cnt[x] = c
       P += c / 2
       Q += c / 4
   }
   var q int
   fmt.Fscan(reader, &q)
   for i := 0; i < q; i++ {
       var op string
       fmt.Fscan(reader, &op, &x)
       c := cnt[x]
       P -= c / 2
       Q -= c / 4
       if op == "+" {
           c++
       } else {
           c--
       }
       cnt[x] = c
       P += c / 2
       Q += c / 4
       if Q >= 2 || (Q >= 1 && P >= 4) {
           fmt.Fprintln(writer, "YES")
       } else {
           fmt.Fprintln(writer, "NO")
       }
   }
}
