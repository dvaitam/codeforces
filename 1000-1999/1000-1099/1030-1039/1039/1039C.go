package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

const MOD = 1000000007

// DSU find with path compression
func find(u int, parent []int) int {
   if parent[u] != u {
       parent[u] = find(parent[u], parent)
   }
   return parent[u]
}

// union returns true if merged two components
func unify(u, v int, parent []int) bool {
   ru := find(u, parent)
   rv := find(v, parent)
   if ru == rv {
       return false
   }
   parent[rv] = ru
   return true
}

// modPow computes a^b mod mod
func modPow(a, b int64, mod int64) int64 {
   res := int64(1)
   a %= mod
   for b > 0 {
       if b&1 == 1 {
           res = res * a % mod
       }
       a = a * a % mod
       b >>= 1
   }
   return res
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, m int
   var k int64
   fmt.Fscan(reader, &n, &m, &k)
   c := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &c[i])
   }
   // read edges and compute xor value
   type Edge struct{ val int64; x, y int }
   edges := make([]Edge, m)
   for i := 0; i < m; i++ {
       var x, y int
       fmt.Fscan(reader, &x, &y)
       x--
       y--
       edges[i] = Edge{c[x] ^ c[y], x, y}
   }
   sort.Slice(edges, func(i, j int) bool { return edges[i].val < edges[j].val })

   // precompute powers of two up to n
   pow2 := make([]int64, n+1)
   pow2[0] = 1
   for i := 1; i <= n; i++ {
       pow2[i] = pow2[i-1] * 2 % MOD
   }

   parent := make([]int, n)
   for i := 0; i < n; i++ {
       parent[i] = i
   }
   visited := make([]bool, n)
   touched := make([]int, 0, 16)

   var cntGroups int64
   var ans int64
   // process each group of equal xor value
   for i := 0; i < m; {
       cntGroups++
       j := i
       // reset touched list
       touched = touched[:0]
       sz := n
       // process edges with same value
       for j < m && edges[j].val == edges[i].val {
           x := edges[j].x
           y := edges[j].y
           if !visited[x] {
               visited[x] = true
               touched = append(touched, x)
           }
           if !visited[y] {
               visited[y] = true
               touched = append(touched, y)
           }
           if unify(x, y, parent) {
               sz--
           }
           j++
       }
       // add contribution for this group
       ans = (ans + pow2[sz]) % MOD
       // reset DSU for touched nodes
       for _, u := range touched {
           parent[u] = u
           visited[u] = false
       }
       i = j
   }
   // account for groups without any edges: 2^k total xor values
   totalComb := modPow(2, k, MOD)
   rem := (totalComb - cntGroups%MOD + MOD) % MOD
   ans = (ans + rem*pow2[n]%MOD) % MOD
   fmt.Fprintln(writer, ans)
}
