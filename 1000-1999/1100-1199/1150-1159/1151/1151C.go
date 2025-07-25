package main

import (
   "bufio"
   "fmt"
   "os"
)

const MOD = 1000000007

// f returns the sum of the first n numbers in the sequence described, mod MOD
func f(n int64) int64 {
   if n <= 0 {
       return 0
   }
   var sumMod, oddCntMod, evenCntMod int64
   var totalPos int64
   stage := 1
   length := int64(1)
   for totalPos < n {
       rem := n - totalPos
       k := length
       if rem < length {
           k = rem
       }
       kMod := k % MOD
       if stage&1 == 1 {
           // odd numbers: values are 2*(oddCnt) + 1 + 2*t for t in [0,k)
           // sum = k*(2*oddCnt + k)
           sumStage := kMod * ((2*oddCntMod + kMod) % MOD) % MOD
           sumMod = (sumMod + sumStage) % MOD
           oddCntMod = (oddCntMod + kMod) % MOD
       } else {
           // even numbers: values are 2*(evenCnt + t) for t in [1,k]
           // sum = 2*k*evenCnt + k*(k+1)
           sumStage := (2*(kMod*evenCntMod%MOD)%MOD + (kMod*((kMod+1)%MOD)%MOD)) % MOD
           sumMod = (sumMod + sumStage) % MOD
           evenCntMod = (evenCntMod + kMod) % MOD
       }
       totalPos += k
       if rem < length {
           break
       }
       length <<= 1
       stage++
   }
   return sumMod
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var l, r int64
   if _, err := fmt.Fscan(reader, &l, &r); err != nil {
       return
   }
   res := (f(r) - f(l-1) + MOD) % MOD
   fmt.Println(res)
}
