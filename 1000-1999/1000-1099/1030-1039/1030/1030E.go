package main

import (
   "bufio"
   "fmt"
   "math/bits"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n int
   fmt.Fscan(in, &n)
   c := make([]int, n+1)
   for i := 1; i <= n; i++ {
       var a uint64
       fmt.Fscan(in, &a)
       c[i] = bits.OnesCount64(a)
   }

   // prefix parity counts
   cnt := [2]int64{1, 0}
   par := 0
   for i := 1; i <= n; i++ {
       par ^= c[i] & 1
       cnt[par]++
   }
   // total subarrays with even sum of popcounts
   ans := cnt[0]*(cnt[0]-1)/2 + cnt[1]*(cnt[1]-1)/2

   // subtract bad segments: small segments where max popcount > sum/2
   const K = 128
   for r := 1; r <= n; r++ {
       sum, mx := 0, 0
       for l := r; l >= 1 && r-l < K; l-- {
           sum += c[l]
           if c[l] > mx {
               mx = c[l]
           }
           if sum%2 == 0 && mx*2 > sum {
               ans--
           }
       }
   }
   fmt.Fprintln(out, ans)
}
