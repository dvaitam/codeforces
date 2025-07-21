package main

import (
   "bufio"
   "fmt"
   "os"
)

// Bitset for subset sum up to n bits
type Bitset struct {
   data []uint64
}

// NewBitset creates a bitset with bits [0..n]
func NewBitset(n int) *Bitset {
   size := (n + 64) / 64
   return &Bitset{data: make([]uint64, size)}
}

// Set bit i to 1
func (b *Bitset) Set(i int) {
   b.data[i/64] |= 1 << (uint(i) % 64)
}

// Get returns true if bit i is 1
func (b *Bitset) Get(i int) bool {
   return (b.data[i/64] & (1 << (uint(i) % 64))) != 0
}

// ShiftOr shifts the bitset left by s and ORs into b (in-place)
func (b *Bitset) ShiftOr(s int) {
   if s <= 0 {
       return
   }
   wordShift := s / 64
   bitShift := uint(s % 64)
   n := len(b.data)
   // process from high to low to avoid overwrite issues
   for i := n - 1; i >= 0; i-- {
       var v uint64
       j := i - wordShift
       if j >= 0 {
           v = b.data[j] << bitShift
           if bitShift > 0 && j-1 >= 0 {
               v |= b.data[j-1] >> (64 - bitShift)
           }
       }
       b.data[i] |= v
   }
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   fmt.Fscan(reader, &n)
   adj := make([][]int, n+1)
   for i := 0; i < n-1; i++ {
       var x, y int
       fmt.Fscan(reader, &x, &y)
       adj[x] = append(adj[x], y)
       adj[y] = append(adj[y], x)
   }
   parent := make([]int, n+1)
   siz := make([]int, n+1)
   // DFS to compute parent and subtree sizes
   var dfs func(u, p int)
   dfs = func(u, p int) {
       parent[u] = p
       siz[u] = 1
       for _, v := range adj[u] {
           if v == p {
               continue
           }
           dfs(v, u)
           siz[u] += siz[v]
       }
   }
   dfs(1, 0)

   possible := make([]bool, n+1)
   // For each vertex v, consider removing v as empty
   for v := 1; v <= n; v++ {
       deg := len(adj[v])
       if deg < 2 {
           continue
       }
       // Collect component sizes when removing v
       sizes := make([]int, 0, deg)
       for _, u := range adj[v] {
           if parent[u] == v {
               sizes = append(sizes, siz[u])
           } else {
               // u is parent of v
               sizes = append(sizes, n - siz[v])
           }
       }
       // subset sum bitset over sizes
       bs := NewBitset(n)
       bs.Set(0)
       for _, s := range sizes {
           bs.ShiftOr(s)
       }
       // mark possible sums a from 1 to n-2
       total := n - 1
       for a := 1; a <= total-1; a++ {
           if bs.Get(a) {
               possible[a] = true
           }
       }
   }
   // collect answers
   var ans [][2]int
   total := n - 1
   for a := 1; a <= total-1; a++ {
       if possible[a] {
           ans = append(ans, [2]int{a, total - a})
       }
   }
   // output
   fmt.Fprintln(writer, len(ans))
   for _, p := range ans {
       fmt.Fprintf(writer, "%d %d\n", p[0], p[1])
   }
}
