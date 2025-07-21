package main

import (
   "bufio"
   "fmt"
   "math/big"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var s string
   if _, err := fmt.Fscan(reader, &s); err != nil {
       return
   }
   n := len(s)
   // dpPrev[d] is number of ways ending with digit d at current position
   dpPrev := make([]*big.Int, 10)
   dpCur := make([]*big.Int, 10)
   // Initialize for first digit: any favorite digit 0-9
   for d := 0; d < 10; d++ {
       dpPrev[d] = big.NewInt(1)
       dpCur[d] = big.NewInt(0)
   }
   // Process positions 2..n
   for i := 1; i < n; i++ {
       // reset current dp
       for d := 0; d < 10; d++ {
           dpCur[d].SetInt64(0)
       }
       ai := int(s[i] - '0')
       for prev := 0; prev < 10; prev++ {
           count := dpPrev[prev]
           if count.Sign() == 0 {
               continue
           }
           sum := ai + prev
           if sum%2 == 0 {
               next := sum / 2
               dpCur[next].Add(dpCur[next], count)
           } else {
               low := sum / 2
               high := low + 1
               dpCur[low].Add(dpCur[low], count)
               dpCur[high].Add(dpCur[high], count)
           }
       }
       // swap dpPrev and dpCur
       dpPrev, dpCur = dpCur, dpPrev
   }
   // Sum up all possibilities at last position
   result := big.NewInt(0)
   if n > 0 {
       for d := 0; d < 10; d++ {
           result.Add(result, dpPrev[d])
       }
   }
   fmt.Println(result.String())
}
