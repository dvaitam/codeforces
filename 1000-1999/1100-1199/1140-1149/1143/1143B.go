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

   var n int64
   fmt.Fscan(reader, &n)

   // count digits
   temp := n
   x := 0
   for temp > 0 {
       x++
       temp /= 10
   }
   if x == 0 {
       x = 1
   }

   // product of digits with zero-break
   prod := func(m int64) int64 {
       ans := int64(1)
       for i := 0; i < x; i++ {
           d := m % 10
           if d == 0 {
               break
           }
           ans *= d
           m /= 10
       }
       return ans
   }

   // initial maximum
   mans := int64(1)
   if p := prod(n); p > mans {
       mans = p
   }
   if p := prod(n - 1); p > mans {
       mans = p
   }

   // try decreasing prefix
   for x1 := x - 1; x1 >= 0; x1-- {
       pow10 := int64(1)
       for i := 0; i < x1; i++ {
           pow10 *= 10
       }
       n1 := n - (n % pow10) - 1
       if n1 < 0 {
           continue
       }
       if p := prod(n1); p > mans {
           mans = p
       }
   }

   fmt.Fprint(writer, mans)
}
