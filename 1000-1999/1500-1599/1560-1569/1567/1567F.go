package main

import (
   "bufio"
   "fmt"
   "os"
)

// DSU with parity
type DSU struct {
   parent     []int
   sz         []int
   f          []bool
   components int
}

// NewDSU initializes DSU of size n
func NewDSU(n int) *DSU {
   parent := make([]int, n)
   sz := make([]int, n)
   f := make([]bool, n)
   for i := 0; i < n; i++ {
       parent[i] = i
       sz[i] = 1
       f[i] = false
   }
   return &DSU{parent: parent, sz: sz, f: f, components: n}
}

// find with path compression and parity update
func (d *DSU) find(a int) int {
   if d.parent[a] == a {
       return a
   }
   p := d.parent[a]
   root := d.find(p)
   d.f[a] = d.f[a] != d.f[p]
   d.parent[a] = root
   return root
}

// link merges b into a with parity myf
func (d *DSU) link(a, b int, myf bool) {
   d.components--
   if d.sz[a] < d.sz[b] {
       a, b = b, a
   }
   d.sz[a] += d.sz[b]
   d.parent[b] = a
   d.f[b] = myf
}

// Unite connects a and b, returns true if merged
func (d *DSU) Unite(a, b int) bool {
   pa := d.find(a)
   pb := d.find(b)
   if pa != pb {
       myf := d.f[a] == d.f[b]
       d.link(pa, pb, myf)
       return true
   }
   return false
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   grid := make([]string, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &grid[i])
   }

   // direction vectors
   dirs := [][2]int{{1, 0}, {0, 1}, {-1, 0}, {0, -1}}

   id := func(i, j int) int { return i*m + j }

   dsu := NewDSU(n * m)
   res := make([][]int, n)
   for i := 0; i < n; i++ {
       res[i] = make([]int, m)
   }

   // pair up neighbors of 'X'
   for i := 0; i < n; i++ {
       for j := 0; j < m; j++ {
           if grid[i][j] != 'X' {
               continue
           }
           var ne []int
           for _, d := range dirs {
               ni, nj := i+d[0], j+d[1]
               if ni >= 0 && ni < n && nj >= 0 && nj < m && grid[ni][nj] == '.' {
                   ne = append(ne, id(ni, nj))
               }
           }
           if len(ne)%2 == 1 {
               fmt.Fprintln(writer, "NO")
               return
           }
           for k := 1; k < len(ne); k++ {
               dsu.Unite(ne[k-1], ne[k])
           }
       }
   }

   // assign values for '.' cells
   for i := 0; i < n; i++ {
       for j := 0; j < m; j++ {
           if grid[i][j] != '.' {
               continue
           }
           idx := id(i, j)
           dsu.find(idx)
           if dsu.f[idx] {
               res[i][j] = 1
           } else {
               res[i][j] = 4
           }
       }
   }

   // compute values for 'X' cells
   for i := 0; i < n; i++ {
       for j := 0; j < m; j++ {
           if grid[i][j] != 'X' {
               continue
           }
           sum := 0
           for _, d := range dirs {
               ni, nj := i+d[0], j+d[1]
               if ni >= 0 && ni < n && nj >= 0 && nj < m && grid[ni][nj] == '.' {
                   sum += res[ni][nj]
               }
           }
           res[i][j] = sum
       }
   }

   // output
   fmt.Fprintln(writer, "YES")
   for i := 0; i < n; i++ {
       for j := 0; j < m; j++ {
           if j > 0 {
               writer.WriteByte(' ')
           }
           fmt.Fprint(writer, res[i][j])
       }
       writer.WriteByte('\n')
   }
}
