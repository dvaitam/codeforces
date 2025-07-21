package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m int
   var mod int64
   if _, err := fmt.Fscan(reader, &n, &m, &mod); err != nil {
       return
   }
   // read first m rows and count ones per column
   a := make([]int, n)
   for i := 0; i < m; i++ {
       var s string
       fmt.Fscan(reader, &s)
       for j, ch := range s {
           if ch == '1' {
               a[j]++
               if a[j] > 2 {
                   fmt.Println(0)
                   return
               }
           }
       }
   }
   // remaining rows
   R := n - m
   // capacities per column
   c := make([]int, n)
   for j := 0; j < n; j++ {
       c[j] = 2 - a[j]
       if c[j] < 0 {
           fmt.Println(0)
           return
       }
   }
   // dp[k]: number of ways with k open singles
   dp := make([]int64, R+3)
   dp[0] = 1
   S := 0 // total assigned slots so far
   for _, cj := range c {
       ndp := make([]int64, R+3)
       switch cj {
       case 0:
           // no change
           for k := 0; k <= R; k++ {
               ndp[k] = dp[k]
           }
       case 1:
           for k := 0; k <= R; k++ {
               v := dp[k]
               if v == 0 {
                   continue
               }
               // compute zero rows: z = R - (S + k)/2
               z := R - (S + k)/2
               // assign to open row -> close one
               if k > 0 {
                   ndp[k-1] = (ndp[k-1] + v*int64(k)) % mod
               }
               // assign to zero row -> become open
               if z > 0 {
                   ndp[k+1] = (ndp[k+1] + v*int64(z)) % mod
               }
           }
       case 2:
           for k := 0; k <= R; k++ {
               v := dp[k]
               if v == 0 {
                   continue
               }
               z := R - (S + k)/2
               // both to open rows: k choose 2
               if k >= 2 {
                   ways := int64(k*(k-1)/2)
                   ndp[k-2] = (ndp[k-2] + v*ways) % mod
               }
               // one to open, one to zero
               if k > 0 && z > 0 {
                   ways := int64(k * z)
                   ndp[k] = (ndp[k] + v*ways) % mod
               }
               // both to zero rows: z choose 2
               if z >= 2 {
                   ways := int64(z*(z-1)/2)
                   ndp[k+2] = (ndp[k+2] + v*ways) % mod
               }
           }
       }
       S += cj
       dp = ndp
   }
   // result: all rows closed => 0 open singles
   fmt.Println(dp[0] % mod)
}
