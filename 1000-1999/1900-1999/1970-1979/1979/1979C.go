package main

import (
   "bufio"
   "fmt"
   "os"
)

func gcd(a, b int) int {
   for b != 0 {
       a, b = b, a%b
   }
   return a
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var T int
   fmt.Fscan(reader, &T)

   // Precompute L = lcm(1..20)
   L := 1
   for i := 2; i <= 20; i++ {
       L = L / gcd(L, i) * i
   }

   for T > 0 {
       T--
       var n int
       fmt.Fscan(reader, &n)
       a := make([]int, n)
       for i := 0; i < n; i++ {
           fmt.Fscan(reader, &a[i])
       }
       ans := make([]int, n)
       var sum int64
       for i, x := range a {
           v := L / x
           ans[i] = v
           sum += int64(v)
       }
       if sum >= int64(L) {
           fmt.Fprintln(writer, -1)
       } else {
           for i, v := range ans {
               if i > 0 {
                   writer.WriteByte(' ')
               }
               fmt.Fprint(writer, v)
           }
           writer.WriteByte('\n')
       }
   }
}
