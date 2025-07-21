package main

import (
   "bufio"
   "fmt"
   "os"
)

const MOD = 1000000007

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   a := make([]int, n)
   total := 0
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
       total += a[i]
   }
   // dpPrev[s] = number of ways for subarrays ending at i-1 with sum = s-offset
   offset := total
   size := 2*total + 1
   dpPrev := make([]int, size)
   dpCurr := make([]int, size)
   var ans int
   for i := 0; i < n; i++ {
       ai := a[i]
       // reset dpCurr
       for j := 0; j < size; j++ {
           dpCurr[j] = 0
       }
       // start new subarray at i
       dpCurr[offset+ai] = (dpCurr[offset+ai] + 1) % MOD
       dpCurr[offset-ai] = (dpCurr[offset-ai] + 1) % MOD
       // extend previous subarrays
       for j := 0; j < size; j++ {
           v := dpPrev[j]
           if v != 0 {
               // sum = j-offset, add ai or subtract ai
               jp := j + ai
               if jp < size {
                   dpCurr[jp] = (dpCurr[jp] + v) % MOD
               }
               jm := j - ai
               if jm >= 0 {
                   dpCurr[jm] = (dpCurr[jm] + v) % MOD
               }
           }
       }
       // add count of sum zero for subarrays ending at i
       ans = (ans + dpCurr[offset]) % MOD
       // swap dpPrev and dpCurr
       dpPrev, dpCurr = dpCurr, dpPrev
   }
   fmt.Println(ans)
}
