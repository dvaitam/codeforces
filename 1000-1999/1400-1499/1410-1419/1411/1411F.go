package main

import (
   "bufio"
   "fmt"
   "os"
)

const mod = 1000000007

// fast exponentiation a^e % mod
func modPow(a int64, e int) int64 {
   res := int64(1)
   base := a % mod
   for e > 0 {
       if e&1 == 1 {
           res = (res * base) % mod
       }
       base = (base * base) % mod
       e >>= 1
   }
   return res
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   var t int
   fmt.Fscan(in, &t)
   for ; t > 0; t-- {
       var n int
       fmt.Fscan(in, &n)
       p := make([]int, n+1)
       for i := 1; i <= n; i++ {
           fmt.Fscan(in, &p[i])
       }
       // count initial cycles
       visited := make([]bool, n+1)
       k0 := 0
       for i := 1; i <= n; i++ {
           if !visited[i] {
               k0++
               for j := i; !visited[j]; j = p[j] {
                   visited[j] = true
               }
           }
       }
       // compute maximum number of days
       rem := n % 3
       var days int64
       switch rem {
       case 0:
           days = modPow(3, n/3)
       case 1:
           // use one less 3 to make a 4
           if n >= 4 {
               days = modPow(3, n/3-1) * 4 % mod
           } else {
               days = int64(n)
           }
       case 2:
           days = modPow(3, n/3) * 2 % mod
       }
       // target cycles count: floor((n+2)/3)
       kt := (n + 2) / 3
       swaps := k0 - kt
       if swaps < 0 {
           swaps = -swaps
       }
       fmt.Fprintln(out, days, swaps)
   }
}
