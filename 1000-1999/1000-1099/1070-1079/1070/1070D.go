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
   var k int64
   fmt.Fscan(reader, &n, &k)
   a := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   var carry int64 = 0
   var bags int64 = 0
   for i := 0; i < n; i++ {
       t := carry + a[i]
       if i == n-1 {
           // last day: dispose all
           bags += (t + k - 1) / k
           carry = 0
       } else {
           // try to leave minimal leftover to make disposals divisible by k
           rem := t % k
           if rem <= a[i] {
               // can carry rem to next day, dispose t-rem
               bags += t / k
               carry = rem
           } else {
               // cannot carry rem, carry as much as possible (a[i])
               carry = a[i]
               // dispose the rest
               dispose := t - carry
               bags += (dispose + k - 1) / k
           }
       }
   }
   fmt.Fprintln(writer, bags)
}
