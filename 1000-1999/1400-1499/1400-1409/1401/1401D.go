package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

const MOD int64 = 1000000007

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var t int
   fmt.Fscan(in, &t)
   for t > 0 {
       t--
       var n int
       fmt.Fscan(in, &n)
       adj := make([][]int, n)
       for i := 0; i < n-1; i++ {
           var u, v int
           fmt.Fscan(in, &u, &v)
           u--
           v--
           adj[u] = append(adj[u], v)
           adj[v] = append(adj[v], u)
       }
       // compute subtree sizes and edge contributions
       parent := make([]int, n)
       order := make([]int, 0, n)
       parent[0] = -1
       stack := []int{0}
       for len(stack) > 0 {
           u := stack[len(stack)-1]
           stack = stack[:len(stack)-1]
           order = append(order, u)
           for _, v := range adj[u] {
               if v == parent[u] {
                   continue
               }
               parent[v] = u
               stack = append(stack, v)
           }
       }
       sz := make([]int64, n)
       contrib := make([]int64, 0, n-1)
       for i := n-1; i >= 0; i-- {
           u := order[i]
           sz[u] += 1
           p := parent[u]
           if p != -1 {
               // edge u-p
               c := sz[u] * (int64(n) - sz[u])
               contrib = append(contrib, c)
               sz[p] += sz[u]
           }
       }
       sort.Slice(contrib, func(i, j int) bool {
           return contrib[i] > contrib[j]
       })
       var m int
       fmt.Fscan(in, &m)
       primes := make([]int64, m)
       for i := 0; i < m; i++ {
           fmt.Fscan(in, &primes[i])
           primes[i] %= MOD
       }
       // adjust primes to length n-1
       if len(primes) < n-1 {
           // append ones
           for i := len(primes); i < n-1; i++ {
               primes = append(primes, 1)
           }
           sort.Slice(primes, func(i, j int) bool { return primes[i] > primes[j] })
       } else {
           // combine extra primes into one
           sort.Slice(primes, func(i, j int) bool { return primes[i] < primes[j] })
           // combine largest k+1 primes: indices from n-2 to end
           big := int64(1)
           // start index = (n-2)
           for i := n - 2; i < len(primes); i++ {
               big = big * primes[i] % MOD
           }
           // keep first n-2 primes, then big
           newP := make([]int64, 0, n-1)
           for i := 0; i < n-1; i++ {
               if i < n-2 {
                   newP = append(newP, primes[i])
               } else if i == n-2 {
                   newP = append(newP, big)
               }
           }
           primes = newP
           sort.Slice(primes, func(i, j int) bool { return primes[i] > primes[j] })
       }
       // compute result
       var res int64
       for i := 0; i < n-1; i++ {
           res = (res + (contrib[i]%MOD)*primes[i]%MOD) % MOD
       }
       fmt.Fprintln(out, res)
   }
}
