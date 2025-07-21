package main

import (
   "bufio"
   "fmt"
   "math/bits"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var s string
   if _, err := fmt.Fscan(reader, &s); err != nil {
       return
   }
   nDigits := []byte(s)
   L := len(nDigits)
   // dp[pos][tight][started][mask] = count
   dp := make([][2][2][1024]int64, L+1)
   dp[0][1][0][0] = 1
   for pos := 0; pos < L; pos++ {
       for tight := 0; tight < 2; tight++ {
           for started := 0; started < 2; started++ {
               for mask := 0; mask < 1024; mask++ {
                   cur := dp[pos][tight][started][mask]
                   if cur == 0 {
                       continue
                   }
                   limit := 9
                   if tight == 1 {
                       limit = int(nDigits[pos] - '0')
                   }
                   for d := 0; d <= limit; d++ {
                       nt := 0
                       if tight == 1 && d == limit {
                           nt = 1
                       }
                       ns := started
                       nm := mask
                       if started == 0 {
                           if d != 0 {
                               ns = 1
                               nm = 1 << d
                           }
                       } else {
                           nm = mask | (1 << d)
                       }
                       if ns == 1 && bits.OnesCount(uint(nm)) > 2 {
                           continue
                       }
                       dp[pos+1][nt][ns][nm] += cur
                   }
               }
           }
       }
   }
   var res int64
   for tight := 0; tight < 2; tight++ {
       for mask := 0; mask < 1024; mask++ {
           if bits.OnesCount(uint(mask)) <= 2 {
               res += dp[L][tight][1][mask]
           }
       }
   }
   fmt.Fprintln(writer, res)
}
