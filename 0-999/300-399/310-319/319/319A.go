package main

import (
   "bufio"
   "fmt"
   "os"
)

const mod = 1000000007

func main() {
   reader := bufio.NewReader(os.Stdin)
   var x string
   if _, err := fmt.Fscan(reader, &x); err != nil {
       return
   }
   n := len(x)
   // precompute powers of 2 up to 2*n
   maxExp := 2*n
   pow2 := make([]int64, maxExp+1)
   pow2[0] = 1
   for i := 1; i <= maxExp; i++ {
       pow2[i] = pow2[i-1] * 2 % mod
   }
   var ans int64 = 0
   // build f_k for suffix of length k from LSB to MSB
   for k := 1; k <= n; k++ {
       ans = (ans * 2) % mod
       // bit at position k-1 (0-based) is x[n-k]
       if x[n-k] == '1' {
           // add cross-block inversions: 2^{2k-2}
           ans = (ans + pow2[2*k-2]) % mod
       }
   }
   fmt.Println(ans)
}
