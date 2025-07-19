package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m int
   var k int64
   fmt.Fscan(reader, &n, &m, &k)
   a := make([]int, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   adj := make([][]int, n+1)
   edges := make([][2]int, m)
   for i := 0; i < m; i++ {
       u, v := 0, 0
       fmt.Fscan(reader, &u, &v)
       edges[i] = [2]int{u, v}
       adj[u] = append(adj[u], v)
       adj[v] = append(adj[v], u)
   }
   comp := make([]int, n+1)
   // find white components (a[i]==0)
   W := 0
   for i := 1; i <= n; i++ {
       if a[i] == 0 && comp[i] == 0 {
           W++
           // BFS
           queue := []int{i}
           comp[i] = W
           for len(queue) > 0 {
               u := queue[0]
               queue = queue[1:]
               for _, v := range adj[u] {
                   if a[v] == 0 && comp[v] == 0 {
                       comp[v] = W
                       queue = append(queue, v)
                   }
               }
           }
       }
   }
   // index black nodes
   B := 0
   for i := 1; i <= n; i++ {
       if a[i] == 1 {
           B++
           comp[i] = B
       }
   }
   // count edges between black and white, and black-black
   cnt := make([][]int, B)
   for i := 0; i < B; i++ {
       cnt[i] = make([]int, W)
   }
   Cnt := make([][]int, B)
   for i := 0; i < B; i++ {
       Cnt[i] = make([]int, B)
   }
   for _, e := range edges {
       u, v := e[0], e[1]
       if a[u] == 1 && a[v] == 1 {
           ui, vi := comp[u]-1, comp[v]-1
           Cnt[ui][vi]++
           Cnt[vi][ui]++
       } else if a[u] == 1 && a[v] == 0 {
           ui, vj := comp[u]-1, comp[v]-1
           cnt[ui][vj]++
       } else if a[u] == 0 && a[v] == 1 {
           ui, vj := comp[v]-1, comp[u]-1
           cnt[ui][vj]++
       }
   }
   // degrees
   s1 := make([]int, W)
   for j := 0; j < W; j++ {
       sum := 0
       for i := 0; i < B; i++ {
           sum += cnt[i][j]
       }
       s1[j] = sum
   }
   s2 := make([]int, B)
   for i := 0; i < B; i++ {
       sum := 0
       for j := 0; j < B; j++ {
           sum += Cnt[j][i]
       }
       for j := 0; j < W; j++ {
           sum += cnt[i][j]
       }
       s2[i] = sum
   }
   // build transition matrix for blacks
   p := NewMatrix(B)
   for i := 0; i < B; i++ {
       if s2[i] == 0 {
           continue
       }
       for j := 0; j < B; j++ {
           // direct black-black
           p.n[i][j] = float64(Cnt[i][j]) / float64(s2[i])
           // via white
           extra := 0.0
           for k2 := 0; k2 < W; k2++ {
               if s1[k2] == 0 {
                   continue
               }
               extra += (float64(cnt[i][k2]) / float64(s2[i])) * (float64(cnt[j][k2]) / float64(s1[k2]))
           }
           p.n[i][j] += extra
       }
   }
   // exponentiate
   exp := k - 2
   pPow := p.Pow(exp)
   // initial component of node 1
   j0 := comp[1] - 1
   ans := 0.0
   if j0 >= 0 && j0 < W {
       for i := 0; i < B; i++ {
           if s1[j0] > 0 {
               ans += (float64(cnt[i][j0]) / float64(s1[j0])) * pPow.n[j0][B-1]
           }
       }
   }
   fmt.Printf("%.6f", ans)
}

// Matrix represents a square matrix of size sz
type Matrix struct {
   n  [][]float64
   sz int
}

// NewMatrix creates a sz x sz zero matrix
func NewMatrix(sz int) *Matrix {
   n := make([][]float64, sz)
   for i := range n {
       n[i] = make([]float64, sz)
   }
   return &Matrix{n: n, sz: sz}
}

// Multiply returns the product of m and o
func (m *Matrix) Multiply(o *Matrix) *Matrix {
   sz := m.sz
   res := NewMatrix(sz)
   for i := 0; i < sz; i++ {
       for k := 0; k < sz; k++ {
           mik := m.n[i][k]
           if mik != 0 {
               for j := 0; j < sz; j++ {
                   res.n[i][j] += mik * o.n[k][j]
               }
           }
       }
   }
   return res
}

// Pow raises m to the power exp (exp >= 0)
func (m *Matrix) Pow(exp int64) *Matrix {
   sz := m.sz
   // identity
   res := NewMatrix(sz)
   for i := 0; i < sz; i++ {
       res.n[i][i] = 1
   }
   base := m
   for exp > 0 {
       if exp&1 == 1 {
           res = res.Multiply(base)
       }
       base = base.Multiply(base)
       exp >>= 1
   }
   return res
}
