package main

import (
   "bufio"
   "fmt"
   "os"
)

// UnionFind with path compression and union by size
type UnionFind struct {
   parent []int
   size   []int
}

func NewUnionFind(n int) *UnionFind {
   uf := &UnionFind{parent: make([]int, n+1), size: make([]int, n+1)}
   for i := 1; i <= n; i++ {
       uf.parent[i] = i
       uf.size[i] = 1
   }
   return uf
}

func (uf *UnionFind) Find(x int) int {
   if uf.parent[x] != x {
       uf.parent[x] = uf.Find(uf.parent[x])
   }
   return uf.parent[x]
}

func (uf *UnionFind) Union(x, y int) {
   rx := uf.Find(x)
   ry := uf.Find(y)
   if rx == ry {
       return
   }
   // attach smaller to larger
   if uf.size[rx] < uf.size[ry] {
       rx, ry = ry, rx
   }
   uf.parent[ry] = rx
   uf.size[rx] += uf.size[ry]
}

// powMod calculates (base^exp) mod mod
func powMod(base int64, exp int, mod int64) int64 {
   res := int64(1)
   b := base % mod
   for exp > 0 {
       if exp&1 == 1 {
           res = (res * b) % mod
       }
       b = (b * b) % mod
       exp >>= 1
   }
   return res
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m int
   var k int64
   if _, err := fmt.Fscan(reader, &n, &m, &k); err != nil {
       return
   }
   uf := NewUnionFind(n)
   for i := 0; i < m; i++ {
       var a, b int
       fmt.Fscan(reader, &a, &b)
       uf.Union(a, b)
   }
   // collect component sizes
   comp := make(map[int]int)
   for i := 1; i <= n; i++ {
       r := uf.Find(i)
       comp[r]++
   }
   C := len(comp)
   // If already connected, no edges needed: one way
   if C <= 1 {
       fmt.Println(1 % k)
       return
   }
   // product of sizes
   var prod int64 = 1
   for _, sz := range comp {
       prod = (prod * int64(sz)) % k
   }
   // n^(C-2) mod k
   pow := powMod(int64(n), C-2, k)
   ans := (prod * pow) % k
   fmt.Println(ans)
}
