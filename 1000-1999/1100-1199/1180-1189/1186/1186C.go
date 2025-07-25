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

   var a, b string
   if _, err := fmt.Fscan(reader, &a); err != nil {
       return
   }
   if _, err := fmt.Fscan(reader, &b); err != nil {
       return
   }
   n := len(a)
   m := len(b)
   // target parity: sum of bits in b mod 2
   var targetParity int
   for i := 0; i < m; i++ {
       if b[i] == '1' {
           targetParity ^= 1
       }
   }
   // compute initial window parity for a[0..m)
   var winParity int
   for i := 0; i < m && i < n; i++ {
       if a[i] == '1' {
           winParity ^= 1
       }
   }
   var count int64
   // slide window over a
   for i := 0; i <= n-m; i++ {
       if winParity == targetParity {
           count++
       }
       // slide: remove a[i], add a[i+m]
       if i < n-m {
           if a[i] == '1' {
               winParity ^= 1
           }
           if a[i+m] == '1' {
               winParity ^= 1
           }
       }
   }
   fmt.Fprint(writer, count)
}
