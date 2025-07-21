package main

import (
   "bufio"
   "fmt"
   "os"
)

// xor0to returns XOR of all integers from 0 to x inclusive.
func xor0to(x int64) int64 {
   switch x & 3 {
   case 0:
       return x
   case 1:
       return 1
   case 2:
       return x + 1
   default:
       return 0
   }
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int64
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   p := make([]int64, n)
   for i := int64(0); i < n; i++ {
       fmt.Fscan(in, &p[i])
   }
   // xor of p[0] to p[n-2]
   var pxor int64
   for i := int64(0); i+1 < n; i++ {
       pxor ^= p[i]
   }
   // compute T = xor over k=2..n of C(k)
   m := n - 1
   var T int64
   for k := int64(2); k <= n; k++ {
       full := m / k
       rem := m % k
       if full&1 == 1 {
           // each full block contributes xor of 0..k-1
           T ^= xor0to(k - 1)
       }
       // remaining rem values contribute xor of 1..rem = xor0to(rem)
       T ^= xor0to(rem)
   }
   // result Q = pxor xor T
   Q := pxor ^ T
   fmt.Println(Q)
}
