package main

import (
   "bufio"
   "fmt"
   "os"
)

const mod = 1000000007

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, m int
   fmt.Fscan(reader, &n, &m)
   adj := make([][]edge, n)
   for i := 0; i < m; i++ {
       var u, v int
       var w uint64
       fmt.Fscan(reader, &u, &v, &w)
       u--
       v--
       adj[u] = append(adj[u], edge{v, w})
       adj[v] = append(adj[v], edge{u, w})
   }
   // precompute powers of 2 up to cover bit positions and basis sizes
   maxPow := m + 5
   if maxPow < 65 {
       maxPow = 65
   }
   pow2 := make([]int64, maxPow)
   pow2[0] = 1
   for i := 1; i < maxPow; i++ {
       pow2[i] = pow2[i-1] * 2 % mod
   }
   inv2 := (mod + 1) / 2

   visited := make([]bool, n)
   var result int64

   d := make([]uint64, n)
   // process each connected component
   for i := 0; i < n; i++ {
       if visited[i] {
           continue
       }
       // initialize for component
       stack := []pair{{i, -1}}
       visited[i] = true
       d[i] = 0
       nodes := []int{i}
       // linear basis
       var basis [61]uint64
       for len(stack) > 0 {
           u := stack[len(stack)-1].u
           p := stack[len(stack)-1].p
           stack = stack[:len(stack)-1]
           for _, e := range adj[u] {
               v := e.to
               w := e.w
               if !visited[v] {
                   visited[v] = true
                   d[v] = d[u] ^ w
                   stack = append(stack, pair{v, u})
                   nodes = append(nodes, v)
               } else if v != p {
                   // found cycle
                   x := d[u] ^ d[v] ^ w
                   // insert into basis
                   for b := 60; b >= 0; b-- {
                       if (x>>b)&1 == 0 {
                           continue
                       }
                       if basis[b] == 0 {
                           basis[b] = x
                           break
                       }
                       x ^= basis[b]
                   }
               }
           }
       }
       // count basis dimension k and span sum
       k := 0
       for b := 0; b <= 60; b++ {
           if basis[b] != 0 {
               k++
           }
       }
       // process node values and compute contribution per bit
       cntBits := make([]int64, 61)
       for _, u := range nodes {
           for b := 0; b <= 60; b++ {
               if (d[u]>>b)&1 == 1 {
                   cntBits[b]++
               }
           }
       }
       N := int64(len(nodes))
       // number of pairs in component
       pairsAll := N * (N - 1) % mod * int64(inv2) % mod
       // OR of basis vectors to know bits in span
       var orAll uint64
       for b := 0; b <= 60; b++ {
           orAll |= basis[b]
       }
       var compSum int64
       for b := 0; b <= 60; b++ {
           bitVal := pow2[b]
           cnt1 := cntBits[b]
           cnt0 := N - cnt1
           sumPairs := cnt1 * cnt0 % mod
           if (orAll>>b)&1 == 1 {
               // bit is in span: contributes bitVal * 2^{k-1} for each pair
               compSum = (compSum + bitVal * pow2[k-1] % mod * pairsAll % mod) % mod
           } else {
               // bit not in span: contributes bitVal * 2^k for base pairs
               compSum = (compSum + bitVal * pow2[k] % mod * sumPairs % mod) % mod
           }
       }
       result = (result + compSum) % mod
   }
   fmt.Fprintln(writer, result)
}

type edge struct {
   to int
   w  uint64
}

// pair holds node and parent in DFS
type pair struct {
   u, p int
}
