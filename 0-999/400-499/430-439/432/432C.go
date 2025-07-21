package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n int
   fmt.Fscan(in, &n)
   a := make([]int, n+1)
   pos := make([]int, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(in, &a[i])
       pos[a[i]] = i
   }
   // sieve primes up to n+1
   isPrime := make([]bool, n+2)
   for i := 2; i <= n+1; i++ {
       isPrime[i] = true
   }
   for i := 2; i*i <= n+1; i++ {
       if isPrime[i] {
           for j := i * i; j <= n+1; j += i {
               isPrime[j] = false
           }
       }
   }
   // largest prime <= i
   lp := make([]int, n+2)
   last := 0
   for i := 0; i <= n+1; i++ {
       if i >= 2 && isPrime[i] {
           last = i
       }
       lp[i] = last
   }
   // operations
   ops := make([][2]int, 0, 5*n)
   for target := 1; target <= n; target++ {
       // bring value 'target' to position 'target'
       for pos[target] > target {
           cur := pos[target]
           dist := cur - target + 1
           p := lp[dist]
           start := cur - p + 1
           // perform swap start <-> cur
           ops = append(ops, [2]int{start, cur})
           v1 := a[start]
           v2 := a[cur]
           a[start], a[cur] = a[cur], a[start]
           pos[v1] = cur
           pos[v2] = start
       }
   }
   // output
   k := len(ops)
   fmt.Fprintln(out, k)
   for _, op := range ops {
       fmt.Fprintln(out, op[0], op[1])
   }
}
