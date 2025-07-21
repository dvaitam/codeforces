package main

import (
   "bufio"
   "fmt"
   "os"
)

const MOD = 1000000007

func add(a, b int) int {
   a += b
   if a >= MOD {
       a -= MOD
   }
   return a
}

func mul(a, b int64) int64 {
   return (a * b) % MOD
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   var k int64
   if _, err := fmt.Fscan(reader, &n, &k); err != nil {
       return
   }
   a := make([]int, n)
   for i := 0; i < n; i++ {
       var x int
       fmt.Fscan(reader, &x)
       a[i] = x % MOD
   }
   // If k == 0, output original array
   if k == 0 {
       for i, v := range a {
           if i > 0 {
               writer.WriteByte(' ')
           }
           fmt.Fprint(writer, v)
       }
       writer.WriteByte('\n')
       return
   }
   // Precompute inverses up to n
   inv := make([]int, n+1)
   inv[1] = 1
   for i := 2; i <= n; i++ {
       inv[i] = int((MOD - int64(MOD/i) * int64(inv[MOD%i]) % MOD) % MOD)
   }
   // Precompute coefficients c[d] = C(d + k - 1, d)
   c := make([]int, n)
   c[0] = 1
   for d := 1; d < n; d++ {
       // c[d] = c[d-1] * (k - 1 + d) / d
       num := ( (k - 1) + int64(d) ) % MOD
       tmp := (int64(c[d-1]) * num) % MOD
       tmp = (tmp * int64(inv[d])) % MOD
       c[d] = int(tmp)
   }
   // Compute result
   res := make([]int, n)
   for i := 0; i < n; i++ {
       var sum int64 = 0
       // sum over j=0..i: a[j] * c[i-j]
       for j := 0; j <= i; j++ {
           sum += int64(a[j]) * int64(c[i-j])
           if sum >= 8*MOD*MOD {
               sum %= MOD
           }
       }
       res[i] = int(sum % MOD)
   }
   // Output
   for i, v := range res {
       if i > 0 {
           writer.WriteByte(' ')
       }
       fmt.Fprint(writer, v)
   }
   writer.WriteByte('\n')
}
