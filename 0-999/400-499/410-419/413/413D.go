package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, k int
   if _, err := fmt.Fscan(reader, &n, &k); err != nil {
       return
   }
   seq := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &seq[i])
   }
   mod := int64(1000000007)
   // Precompute suffix products: number of continuations for each position
   suff := make([]int64, n)
   if n > 0 {
       suff[n-1] = 1
       for i := n - 2; i >= 0; i-- {
           mult := int64(1)
           if seq[i+1] == 0 {
               mult = 2
           }
           suff[i] = suff[i+1] * mult % mod
       }
   }
   // Masks cover powers 1..k-1, bit t-1
   maxMask := 1 << (k - 1)
   // Precompute transitions for inserting 2^h (h=1 or 2)
   maskNext := make([][]int, 3)
   willWin := make([][]bool, 3)
   for h := 1; h <= 2; h++ {
       maskNext[h] = make([]int, maxMask)
       willWin[h] = make([]bool, maxMask)
       for mask := 0; mask < maxMask; mask++ {
           t := h
           cur := mask
           for {
               if t == k {
                   willWin[h][mask] = true
                   break
               }
               bit := 1 << (t - 1)
               if cur&bit == 0 {
                   cur |= bit
                   maskNext[h][mask] = cur
                   break
               }
               // merge and continue
               cur &^= bit
               t++
           }
       }
   }
   dp := make([]int64, maxMask)
   dp[0] = 1
   var ans int64
   for i, v := range seq {
       newDp := make([]int64, maxMask)
       // determine possible h for this value
       hs := []int{}
       if v == 0 {
           hs = []int{1, 2}
       } else if v == 2 {
           hs = []int{1}
       } else if v == 4 {
           hs = []int{2}
       }
       for mask := 0; mask < maxMask; mask++ {
           ways := dp[mask]
           if ways == 0 {
               continue
           }
           for _, h := range hs {
               if willWin[h][mask] {
                   // count all continuations
                   ans = (ans + ways*suff[i]) % mod
               } else {
                   nm := maskNext[h][mask]
                   newDp[nm] = (newDp[nm] + ways) % mod
               }
           }
       }
       dp = newDp
   }
   fmt.Println(ans)
}
