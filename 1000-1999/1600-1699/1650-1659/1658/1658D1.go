package main

import (
   "bufio"
   "fmt"
   "os"
)

func prefixXor(x int) int {
   switch x & 3 {
   case 0:
       return x
   case 1:
       return 1
   case 2:
       return x + 1
   }
   return 0
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var t int
   if _, err := fmt.Fscan(in, &t); err != nil {
       return
   }
   for ; t > 0; t-- {
       var l, r int
       fmt.Fscan(in, &l, &r)
       n := r - l + 1
       a := make([]int, n)
       v := 0
       for i := 0; i < n; i++ {
           fmt.Fscan(in, &a[i])
           v ^= a[i]
       }
       if n&1 == 1 {
           // XOR of [l..r]
           xr := prefixXor(r) ^ prefixXor(l-1)
           fmt.Fprintln(out, xr^v)
       } else {
           ans := 0
           // for each bit position
           for x := 0; x < 17; x++ {
               tot1 := 0
               // count bits in indices
               for i := 0; i < n; i++ {
                   if (i>>x)&1 == 1 {
                       tot1++
                   }
               }
               tot0 := n - tot1
               // subtract bits from values
               for i := 0; i < n; i++ {
                   if (a[i]>>x)&1 == 1 {
                       tot0--
                   } else {
                       tot1--
                   }
               }
               if tot1 == 0 && tot0 == 0 {
                   ans |= 1 << x
               }
           }
           fmt.Fprintln(out, ans)
       }
   }
}
