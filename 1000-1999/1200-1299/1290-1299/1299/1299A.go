package main

import (
   "bufio"
   "fmt"
   "os"
   "math/bits"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n int
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   a := make([]uint64, n)
   var ma uint64
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &a[i])
       if a[i] > ma {
           ma = a[i]
       }
   }
   // handle single element
   if n == 1 {
       fmt.Fprintln(out, a[0])
       return
   }
   // find smallest ind with 2^ind > ma
   ind := bits.Len64(ma)
   xo := (uint64(1) << ind) - 1
   p := make([]uint64, n)
   s := make([]uint64, n)
   p[0] = xo ^ a[0]
   for i := 1; i < n; i++ {
       p[i] = p[i-1] & (xo ^ a[i])
   }
   s[n-1] = xo ^ a[n-1]
   for i := n - 2; i >= 0; i-- {
       s[i] = s[i+1] & (xo ^ a[i])
   }
   // find element with max unique bits
   var fi, st uint64
   // first element
   fi = s[1] & a[0]
   st = a[0]
   for i := 1; i < n; i++ {
       var x uint64
       if i == n-1 {
           x = p[i-1] & a[i]
       } else {
           x = p[i-1] & s[i+1] & a[i]
       }
       if x > fi {
           fi = x
           st = a[i]
       }
   }
   // output: move first occurrence of st to front
   fmt.Fprint(out, st)
   first := true
   for i := 0; i < n; i++ {
       if a[i] == st && first {
           first = false
           continue
       }
       fmt.Fprint(out, " ", a[i])
   }
   fmt.Fprintln(out)
}
