package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

const mod = 998244353

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   a := make([]uint64, n)
   b := make([]uint64, n)
   for i := 0; i < n; i++ {
       var x uint64
       fmt.Fscan(reader, &x)
       a[i] = x
   }
   for i := 0; i < n; i++ {
       var x uint64
       fmt.Fscan(reader, &x)
       b[i] = x
   }
   // compute weights c[i] = a[i] * i*(n-i+1)
   c := make([]uint64, n)
   nn := uint64(n)
   for i := 0; i < n; i++ {
       idx := uint64(i + 1)
       w := idx * (nn - idx + 1)
       c[i] = w * a[i]
   }
   // sort c descending, b ascending
   sort.Slice(c, func(i, j int) bool { return c[i] > c[j] })
   sort.Slice(b, func(i, j int) bool { return b[i] < b[j] })
   var ans uint64
   for i := 0; i < n; i++ {
       ci := c[i] % mod
       bi := b[i] % mod
       ans = (ans + ci*bi) % mod
   }
   fmt.Println(ans)
}
