package main

import (
   "bufio"
   "fmt"
   "os"
)

func gcd(a, b int64) int64 {
   for b != 0 {
       a, b = b, a%b
   }
   return a
}

func lcm(a, b int64) int64 {
   return a / gcd(a, b) * b
}

func good(v, cnt1, cnt2, x, y int64) bool {
   // numbers acceptable for first friend: not divisible by x
   s1 := v - v/x
   // for second friend: not divisible by y
   s2 := v - v/y
   // total acceptable: numbers not divisible by both x and y (i.e., divisible by lcm(x,y) are bad for both)
   s12 := v - v/lcm(x, y)
   return s1 >= cnt1 && s2 >= cnt2 && s12 >= cnt1+cnt2
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var cnt1, cnt2, x, y int64
   if _, err := fmt.Fscan(reader, &cnt1, &cnt2, &x, &y); err != nil {
       return
   }
   // binary search for minimal v
   var lo, hi int64 = 0, 1
   // find upper bound
   for !good(hi, cnt1, cnt2, x, y) {
       hi <<= 1
   }
   for lo < hi {
       mid := lo + (hi-lo)/2
       if good(mid, cnt1, cnt2, x, y) {
           hi = mid
       } else {
           lo = mid + 1
       }
   }
   fmt.Println(lo)
}
