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

   var n int
   fmt.Fscan(reader, &n)
   a := make([]string, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   if n == 0 {
       fmt.Fprint(writer, 0)
       return
   }
   k := len(a[0])
   // sum of digits at each position
   S := make([]int64, k)
   for i := 0; i < n; i++ {
       s := a[i]
       for j := 0; j < k; j++ {
           S[j] += int64(s[j] - '0')
       }
   }
   const mod = 998244353
   // precompute powers of 10
   pow10 := make([]int64, 2*k+1)
   pow10[0] = 1
   for i := 1; i <= 2*k; i++ {
       pow10[i] = pow10[i-1] * 10 % mod
   }
   var sum int64
   for idx := 0; idx < k; idx++ {
       offset := k - 1 - idx
       exp := 2 * offset
       weight := (11 * pow10[exp]) % mod
       sum = (sum + S[idx]%mod * weight) % mod
   }
   result := sum * int64(n) % mod
   fmt.Fprint(writer, result)
}
