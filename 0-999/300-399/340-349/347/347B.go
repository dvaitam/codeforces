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
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }

   fixed := 0
   hasNonFixed := false
   hasMutual := false
   for i := 0; i < n; i++ {
       if a[i] == i {
           fixed++
       } else {
           hasNonFixed = true
           j := a[i]
           if j >= 0 && j < n && a[j] == i {
               hasMutual = true
           }
       }
   }

   result := fixed
   if hasMutual {
       // swapping mutual pair fixes two more
       result += 2
   } else if hasNonFixed {
       // can fix one by swapping
       result += 1
   }
   // else all fixed, result == n

   fmt.Fprint(writer, result)
}
