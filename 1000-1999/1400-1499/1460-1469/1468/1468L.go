package main

import (
   "bufio"
   "fmt"
   "math"
   "math/bits"
   "os"
   "strconv"
)

// modMul computes (a * b) % mod avoiding overflow using 128-bit emulation
func modMul(a, b, mod uint64) uint64 {
   hi, lo := bits.Mul64(a, b)
   // hi<<64 + lo mod mod = (hi * 2^64 mod mod + lo mod mod) mod mod
   // compute 2^64 % mod via (2^32)^2 mod mod
   p := (uint64(1) << 32) % mod
   two64 := (p * p) % mod
   return (two64*(hi%mod)%mod + lo%mod) % mod
}

// modPow computes (base^exp) % mod
func modPow(base, exp, mod uint64) uint64 {
   result := uint64(1)
   base %= mod
   for exp > 0 {
       if exp&1 == 1 {
           result = modMul(result, base, mod)
       }
       base = modMul(base, base, mod)
       exp >>= 1
   }
   return result
}

// isPrime returns true if n is probably prime
func isPrime(n uint64) bool {
   if n < 4 {
       return n == 2 || n == 3
   }
   if n%2 == 0 {
       return false
   }
   // write n-1 as d*2^s
   d := n - 1
   s := 0
   for d&1 == 0 {
       d >>= 1
       s++
   }
   // deterministic bases for 64-bit
   bases := []uint64{2, 325, 9375, 28178, 450775, 9780504, 1795265022}
   for _, a := range bases {
       if a%n == 0 {
           continue
       }
       x := modPow(a, d, n)
       if x == 1 || x == n-1 {
           continue
       }
       composite := true
       for r := 1; r < s; r++ {
           x = modMul(x, x, n)
           if x == n-1 {
               composite = false
               break
           }
       }
       if composite {
           return false
       }
   }
   return true
}

// myprimePower returns base c if n == c^k and c is prime, else 0
func myprimePower(n int64, k int) int64 {
   if n < 2 {
       return 0
   }
   // approximate k-th root
   c0 := int64(math.Round(math.Pow(float64(n), 1.0/float64(k))))
   if c0 < 2 {
       return 0
   }
   var pow uint64
   c := uint64(c0)
   // adjust c to find exact root
   for {
       pow = 1
       for i := 0; i < k; i++ {
           pow *= c
           if pow > uint64(n) {
               break
           }
       }
       if pow == uint64(n) {
           if isPrime(c) {
               return int64(c)
           }
           return 0
       }
       if pow < uint64(n) {
           c++
       } else {
           if c == 0 {
               break
           }
           c--
       }
       if c < 2 {
           break
       }
   }
   return 0
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, k int
   if _, err := fmt.Fscan(reader, &n, &k); err != nil {
       return
   }
   composite := make([]int64, 0, n)
   ps := make(map[int64][]int64)
   for i := 0; i < n; i++ {
       var a int64
       fmt.Fscan(reader, &a)
       pp := false
       for j := 1; j <= 60; j++ {
           if c := myprimePower(a, j); c > 0 {
               ps[c] = append(ps[c], a)
               pp = true
               break
           }
       }
       if !pp {
           composite = append(composite, a)
       }
   }
   // collect prime power groups with at least two elements
   type pair struct{ p int64; v []int64 }
   ppairs := make([]pair, 0, len(ps))
   for p, v := range ps {
       if len(v) > 1 {
           ppairs = append(ppairs, pair{p: p, v: v})
       }
   }
   m := len(composite)
   // build adjacency and indegree
   indeg := make([]int, m)
   red := make([]int64, m)
   copy(red, composite)
   adj := make(map[int64][]int)
   for _, pr := range ppairs {
       p := pr.p
       for i := 0; i < m; i++ {
           if red[i]%p == 0 {
               for red[i]%p == 0 {
                   red[i] /= p
               }
               adj[p] = append(adj[p], i)
               indeg[i]++
           }
       }
   }
   const INF = int(1e9)
   for i := 0; i < m; i++ {
       if red[i] != 1 {
           indeg[i] = INF
       }
   }
   // prioritize groups connected to minimal indegree composite
   if m > 0 {
       b, min := 0, indeg[0]
       for i := 1; i < m; i++ {
           if indeg[i] < min {
               min = indeg[i]; b = i
           }
       }
       first := make([]pair, 0, len(ppairs))
       second := make([]pair, 0, len(ppairs))
       for _, pr := range ppairs {
           found := false
           for _, idx := range adj[pr.p] {
               if idx == b {
                   found = true; break
               }
           }
           if found {
               first = append(first, pr)
           } else {
               second = append(second, pr)
           }
       }
       ppairs = append(first, second...)
   }
   ans := make([]int64, 0, k)
   extra := make([]int64, 0)
   idxp := 0
   for k >= 2 && idxp < len(ppairs) {
       pr := ppairs[idxp]
       v := pr.v
       ans = append(ans, v[0], v[1])
       if len(v) > 2 {
           extra = append(extra, v[2:]...)
       }
       for _, id := range adj[pr.p] {
           indeg[id]--
       }
       k -= 2
       idxp++
   }
   for k > 0 && len(extra) > 0 {
       last := extra[len(extra)-1]
       extra = extra[:len(extra)-1]
       ans = append(ans, last)
       k--
   }
   idxc := 0
   for k > 0 && idxc < m {
       if indeg[idxc] == 0 {
           ans = append(ans, composite[idxc])
           k--
       }
       idxc++
   }
   if k == 0 {
       for i, v := range ans {
           if i > 0 {
               writer.WriteByte(' ')
           }
           writer.WriteString(strconv.FormatInt(v, 10))
       }
       writer.WriteByte('\n')
   } else {
       writer.WriteString("0\n")
   }
}
