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

func main() {
   in := bufio.NewReader(os.Stdin)
   var a, b, w, x, c int64
   if _, err := fmt.Fscan(in, &a, &b, &w, &x, &c); err != nil {
       return
   }
   // If Alexander already ahead
   if c <= a {
       fmt.Println(0)
       return
   }
   // Required difference of no-wrap steps
   D := c - a
   // Build wrap pattern
   g := gcd(w, x)
   L := w / g
   wrap := make([]int64, L)
   var S int64
   for i := int64(0); i < L; i++ {
       // b_i = (b - i*x) mod w
       bi := (b - (i*x)%w + w) % w
       if bi < x {
           wrap[i] = 1
           S++
       }
   }
   noWrap := L - S
   // Compute number of full cycles and remaining no-wraps needed
   // Let D = k*noWrap + rem; if rem==0, use last wrap in previous cycle
   k := D / noWrap
   rem := D - k*noWrap
   if rem == 0 {
       rem = noWrap
       k--
   }
   // Prefix sums of wrap flags for one cycle
   prefix := make([]int64, L+1)
   for i := int64(0); i < L; i++ {
       prefix[i+1] = prefix[i] + wrap[i]
   }
   // Find minimal r in [1..L] such that (r - wraps_in_r) >= rem
   var r int64
   for i := int64(1); i <= L; i++ {
       if i - prefix[i] >= rem {
           r = i
           break
       }
   }
   t := k*L + r
   fmt.Println(t)
}
