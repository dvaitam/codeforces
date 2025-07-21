package main

import (
   "bufio"
   "fmt"
   "os"
   "math/bits"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, t uint64
   if _, err := fmt.Fscan(in, &n, &t); err != nil {
       return
   }
   // check t is power of two
   if t == 0 || (t & (t - 1)) != 0 {
       fmt.Println(0)
       return
   }
   // t = 2^k, number of zeros required in binary representation
   k := bits.TrailingZeros64(t)
   // precompute combinations up to 64
   const MAX = 64
   var comb [MAX][MAX]uint64
   for i := 0; i < MAX; i++ {
       comb[i][0] = 1
       for j := 1; j <= i; j++ {
           comb[i][j] = comb[i-1][j-1] + comb[i-1][j]
       }
   }
   // bit-length of n
   Ln := 64 - bits.LeadingZeros64(n)
   var ans uint64
   // count full ranges for lengths less than Ln
   for L := k + 1; L < Ln; L++ {
       // choose positions of zeros among lower L-1 bits
       ans += comb[L-1][k]
   }
   // handle numbers of bit-length Ln and <= n
   zeros := 0
   // iterate bits from Ln-2 down to 0
   for i := Ln - 2; i >= 0; i-- {
       if ((n >> uint(i)) & 1) == 1 {
           // place zero here, prefix becomes smaller
           remZeros := int(k) - (zeros + 1)
           remBits := int(i)
           if remZeros >= 0 && remZeros <= remBits {
               ans += comb[remBits][remZeros]
           }
           // or place one and continue equal
       }
       // if bit is zero (or we chose one above), track zeros for equal prefix
       if ((n >> uint(i)) & 1) == 0 {
           zeros++
       }
   }
   // include n itself if it matches
   if zeros == int(k) {
       ans++
   }
   fmt.Println(ans)
}
