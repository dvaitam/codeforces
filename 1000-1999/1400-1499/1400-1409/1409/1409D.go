package main

import (
   "bufio"
   "fmt"
   "os"
)

// digitSum returns the sum of decimal digits of n.
func digitSum(n uint64) uint64 {
   var s uint64
   for n > 0 {
       s += n % 10
       n /= 10
   }
   return s
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var t int
   if _, err := fmt.Fscan(reader, &t); err != nil {
       return
   }
   for ; t > 0; t-- {
       var n uint64
       var s uint64
       fmt.Fscan(reader, &n, &s)
       // if already within limit
       if digitSum(n) <= s {
           fmt.Fprintln(writer, 0)
           continue
       }
       var add uint64
       var pow10 uint64 = 1
       // greedily round up digits from least significant
       for i := 0; i < 19; i++ {
           digit := (n / pow10) % 10
           // amount to add to make this digit zero
           delta := (10 - digit) * pow10
           if delta > 0 {
               add += delta
               n += delta
           }
           if digitSum(n) <= s {
               break
           }
           pow10 *= 10
       }
       fmt.Fprintln(writer, add)
   }
}
