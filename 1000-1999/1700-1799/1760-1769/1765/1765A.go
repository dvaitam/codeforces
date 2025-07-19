package main

import (
   "bufio"
   "fmt"
   "os"
)

// Bitset represents a simple bitset using uint64 words
type Bitset struct { data []uint64 }

// NewBitset creates a bitset for n bits
func NewBitset(n int) Bitset {
   return Bitset{data: make([]uint64, (n+63)/64)}
}

// Set sets the bit at position pos
func (b *Bitset) Set(pos int) {
   b.data[pos>>6] |= 1 << (pos & 63)
}

// Equals checks if two bitsets are equal
func (b *Bitset) Equals(o *Bitset) bool {
   for i := range b.data {
      if b.data[i] != o.data[i] {
         return false
      }
   }
   return true
}

// IsSuperset checks if b is a superset of other (i.e., b & other == other)
func (b *Bitset) IsSuperset(o *Bitset) bool {
   for i := range b.data {
      if b.data[i]&o.data[i] != o.data[i] {
         return false
      }
   }
   return true
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, m int
   fmt.Fscan(reader, &n, &m)
   v := make([]string, n)
   for i := 0; i < n; i++ {
      fmt.Fscan(reader, &v[i])
   }

   // Build bitsets for each column
   bs := make([]Bitset, m)
   for j := 0; j < m; j++ {
      bs[j] = NewBitset(n)
   }
   for i := 0; i < n; i++ {
      for j := 0; j < m; j++ {
         if v[i][j] == '1' {
            bs[j].Set(i)
         }
      }
   }

   // Build graph: edge i->j if column j's ones are subset of i's ones (or equal and i<j)
   g := make([][]int, m)
   for i := 0; i < m; i++ {
      for j := 0; j < m; j++ {
         if i == j {
            continue
         }
         if bs[i].Equals(&bs[j]) {
            if i < j {
               g[i] = append(g[i], j)
            }
         } else if bs[i].IsSuperset(&bs[j]) {
            g[i] = append(g[i], j)
         }
      }
   }

   // Maximum matching (Kuhn's algorithm) on the graph
   mt := make([]int, m)
   for i := range mt {
      mt[i] = -1
   }
   var used []bool
   var dfs func(int) bool
   dfs = func(vv int) bool {
      if used[vv] {
         return false
      }
      used[vv] = true
      for _, to := range g[vv] {
         if mt[to] == -1 || dfs(mt[to]) {
            mt[to] = vv
            return true
         }
      }
      return false
   }
   for i := 0; i < m; i++ {
      used = make([]bool, m)
      dfs(i)
   }

   // Build reverse matching for path reconstruction
   revMt := make([]int, m)
   for i := range revMt {
      revMt[i] = -1
   }
   for to, left := range mt {
      if left != -1 {
         revMt[left] = to
      }
   }

   // Reconstruct paths from free right vertices
   var paths [][]int
   for i := 0; i < m; i++ {
      if mt[i] == -1 {
         path := make([]int, 0)
         cur := i
         for cur != -1 {
            path = append(path, cur)
            cur = revMt[cur]
         }
         // reverse
         for l, r := 0, len(path)-1; l < r; l, r = l+1, r-1 {
            path[l], path[r] = path[r], path[l]
         }
         paths = append(paths, path)
      }
   }
   k := len(paths)

   // Assign group indices, access levels, and developer matrix
   which := make([]int, m)
   access := make([]int, m)
   matrix := make([][]int, n)
   for i := 0; i < n; i++ {
      matrix[i] = make([]int, k)
      for j := 0; j < k; j++ {
         matrix[i][j] = 1
      }
   }
   for pi, p := range paths {
      for _, x := range p {
         which[x] = pi + 1
         zeros := 0
         for i := 0; i < n; i++ {
            if v[i][x] == '0' {
               zeros++
            }
         }
         who := 2 + zeros
         access[x] = who
         for i := 0; i < n; i++ {
            if v[i][x] == '1' && matrix[i][pi] < who {
               matrix[i][pi] = who
            }
         }
      }
   }

   // Output results
   fmt.Fprintln(writer, k)
   for i := 0; i < m; i++ {
      if i > 0 {
         fmt.Fprint(writer, " ")
      }
      fmt.Fprint(writer, which[i])
   }
   fmt.Fprintln(writer)
   for i := 0; i < m; i++ {
      if i > 0 {
         fmt.Fprint(writer, " ")
      }
      fmt.Fprint(writer, access[i])
   }
   fmt.Fprintln(writer)
   for i := 0; i < n; i++ {
      for j := 0; j < k; j++ {
         if j > 0 {
            fmt.Fprint(writer, " ")
         }
         fmt.Fprint(writer, matrix[i][j])
      }
      fmt.Fprintln(writer)
   }
