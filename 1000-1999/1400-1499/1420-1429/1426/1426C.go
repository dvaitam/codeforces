package main

import (
   "bufio"
   "fmt"
   "math"
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
   for i := 0; i < t; i++ {
       var n int64
       fmt.Fscan(reader, &n)
       // minimal moves to reach sum >= n
       // operations: k copies, then m0-1 increments, where m0 = ceil(n/(k+1))
       // total ops = k + m0 - 1
       // iterate k from 0 to sqrt(n)
       if n <= 1 {
           fmt.Fprintln(writer, 0)
           continue
       }
       ans := n - 1 // case k=0: just increments
       // limit for k
       // search optimal number of copy operations k around sqrt(n)
       limit := int(math.Sqrt(float64(n))) + 2
       for k := 1; k <= limit; k++ {
           // number of copies = k, elements = k+1
           // needed base value m0
           parts := int64(k) + 1
           // ceil(n/parts)
           m0 := (n + parts - 1) / parts
           ops := int64(k) + m0 - 1
           if ops < ans {
               ans = ops
           }
       }
       fmt.Fprintln(writer, ans)
   }
}
