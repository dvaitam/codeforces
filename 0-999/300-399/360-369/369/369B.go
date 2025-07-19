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

   var n, k, l, r, sAll, sK int
   fmt.Fscan(reader, &n, &k, &l, &r, &sAll, &sK)

   arr := make([]int, n)
   // Distribute sK among first k elements
   rem := sK % k
   base := sK / k
   for i := 0; i < k; i++ {
       v := base
       if rem > 0 {
           v++
           rem--
       }
       arr[i] = v
   }
   // Distribute remaining sum among remaining elements
   if n > k {
       rem2 := (sAll - sK) % (n - k)
       base2 := (sAll - sK) / (n - k)
       for i := k; i < n; i++ {
           v := base2
           if rem2 > 0 {
               v++
               rem2--
           }
           arr[i] = v
       }
   }
   // Output
   for i, v := range arr {
       if i > 0 {
           fmt.Fprint(writer, " ")
       }
       fmt.Fprint(writer, v)
   }
   fmt.Fprint(writer, "\n")
}
