package main

import (
   "bufio"
   "fmt"
   "os"
)

const MOD = 1000000007

func modPow(a, e int64) int64 {
   res := int64(1)
   a %= MOD
   for e > 0 {
       if e&1 == 1 {
           res = res * a % MOD
       }
       a = a * a % MOD
       e >>= 1
   }
   return res
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var t int
   fmt.Fscan(reader, &t)
   for t > 0 {
       t--
       var n int
       fmt.Fscan(reader, &n)
       p := make([]int, n)
       for i := 0; i < n; i++ {
           fmt.Fscan(reader, &p[i])
           p[i]--
       }
       // count initial cycles
       visited := make([]bool, n)
       c0 := 0
       for i := 0; i < n; i++ {
           if !visited[i] {
               c0++
               for j := i; !visited[j]; j = p[j] {
                   visited[j] = true
               }
           }
       }
       // compute optimal product
       var days int64
       switch n % 3 {
       case 0:
           days = modPow(3, int64(n/3))
       case 1:
           // use one 4 instead of two 3s
           if n >= 4 {
               days = modPow(3, int64(n/3-1)) * 4 % MOD
           } else {
               days = 1
           }
       case 2:
           days = modPow(3, int64(n/3)) * 2 % MOD
       }
       // target number of cycles
       k := n/3
       if n%3 == 2 {
           k++
       }
       swaps := c0 - k
       if swaps < 0 {
           swaps = -swaps
       }
       fmt.Fprintf(writer, "%d %d\n", days, swaps)
   }
}
