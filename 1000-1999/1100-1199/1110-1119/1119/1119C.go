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
   r := make([]int, n)
   c := make([]int, m)
   // Read first matrix and accumulate XORs
   for i := 0; i < n; i++ {
       for j := 0; j < m; j++ {
           var x int
           fmt.Fscan(reader, &x)
           r[i] ^= x
           c[j] ^= x
       }
   }
   // Read second matrix and accumulate XORs
   for i := 0; i < n; i++ {
       for j := 0; j < m; j++ {
           var x int
           fmt.Fscan(reader, &x)
           r[i] ^= x
           c[j] ^= x
       }
   }
   ok := true
   for i := 0; i < n; i++ {
       if r[i] != 0 {
           ok = false
           break
       }
   }
   if ok {
       for j := 0; j < m; j++ {
           if c[j] != 0 {
               ok = false
               break
           }
       }
   }
   if ok {
       fmt.Fprint(writer, "Yes")
   } else {
       fmt.Fprint(writer, "No")
   }
}
