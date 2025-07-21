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

   var n, k int
   fmt.Fscan(reader, &n, &k)
   a := make([]int, n)
   maxA := 0
   minA := int(1e9)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
       if a[i] > maxA {
           maxA = a[i]
       }
       if a[i] < minA {
           minA = a[i]
       }
   }
   // frequency and prefix sum
   cnt := make([]int, maxA+1)
   for _, v := range a {
       cnt[v]++
   }
   ps := make([]int, maxA+1)
   for i := 1; i <= maxA; i++ {
       ps[i] = ps[i-1] + cnt[i]
   }
   ans := 1
   // iterate possible d from 2 to minA
   for d := 2; d <= minA; d++ {
       ok := true
       // for each block [m, m+d), any values in (m+k, m+d) are bad
       for m := 0; m <= maxA; m += d {
           low := m + k + 1
           if low > maxA {
               // no values beyond this
               break
           }
           high := m + d - 1
           if high > maxA {
               high = maxA
           }
           if low > high {
               continue
           }
           if ps[high] - ps[low-1] > 0 {
               ok = false
               break
           }
       }
       if ok {
           ans = d
       }
   }
   fmt.Fprintln(writer, ans)
}
