package main

import (
   "bufio"
   "fmt"
   "os"
)

// modPow computes a^b mod m.
func modPow(a, b, m int) int {
   res := 1 % m
   a %= m
   for b > 0 {
       if b&1 == 1 {
           res = (res * a) % m
       }
       a = (a * a) % m
       b >>= 1
   }
   return res
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   // lcm of 2,3,5,7 is 210
   if n < 3 {
       fmt.Fprint(writer, "-1")
       return
   }
   // compute 10^(n-1) mod 210
   rem := modPow(10, n-1, 210)
   remNeeded := (210 - rem) % 210
   // build digits
   digits := make([]int, n)
   digits[0] = 1
   // add remNeeded to the number
   t := remNeeded
   pos := n - 1
   for t > 0 && pos >= 0 {
       digits[pos] += t % 10
       t /= 10
       pos--
   }
   // propagate carries
   for i := n - 1; i > 0; i-- {
       if digits[i] >= 10 {
           carry := digits[i] / 10
           digits[i] %= 10
           digits[i-1] += carry
       }
   }
   // output result
   for i := 0; i < n; i++ {
       writer.WriteByte(byte('0' + digits[i]))
   }
}
