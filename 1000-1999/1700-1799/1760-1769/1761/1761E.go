package main

import (
   "bufio"
   "fmt"
   "os"
)

// UnionFind implements a disjoint set union structure
type UnionFind struct {
   parent []int
   sz     []int
}

// NewUnionFind creates a UnionFind for 1..n
func NewUnionFind(n int) *UnionFind {
   parent := make([]int, n+1)
   sz := make([]int, n+1)
   for i := 1; i <= n; i++ {
       parent[i] = i
       sz[i] = 1
   }
   return &UnionFind{parent: parent, sz: sz}
}

// Find returns the representative of x
func (u *UnionFind) Find(x int) int {
   if u.parent[x] != x {
       u.parent[x] = u.Find(u.parent[x])
   }
   return u.parent[x]
}

// Union joins the sets of x and y
func (u *UnionFind) Union(x, y int) {
   fx := u.Find(x)
   fy := u.Find(y)
   if fx == fy {
       return
   }
   if u.sz[fx] > u.sz[fy] {
       fx, fy = fy, fx
   }
   u.parent[fx] = fy
   u.sz[fy] += u.sz[fx]
}

// Size returns size of set containing x
func (u *UnionFind) Size(x int) int {
   return u.sz[u.Find(x)]
}

func solve(reader *bufio.Reader, writer *bufio.Writer) {
   var n int
   fmt.Fscan(reader, &n)
   uf := NewUnionFind(n)
   s := make([]string, n+1)
   d := make([]int, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &s[i])
   }
   for i := 1; i <= n; i++ {
       row := s[i]
       for j := 0; j < n; j++ {
           if row[j] == '1' {
               uf.Union(i, j+1)
               d[i]++
           }
       }
   }
   // if all connected
   if uf.Size(1) == n {
       fmt.Fprintln(writer, 0)
       return
   }
   // group nodes by root
   v := make([][]int, n+1)
   roots := make([]int, 0, n)
   for i := 1; i <= n; i++ {
       fi := uf.Find(i)
       if len(v[fi]) == 0 {
           roots = append(roots, fi)
       }
       v[fi] = append(v[fi], i)
       // detect singleton component
       if fi == i && uf.sz[fi] == 1 {
           fmt.Fprintln(writer, 1)
           fmt.Fprintln(writer, i)
           return
       }
   }
   // check incomplete clique in any component
   for _, root := range roots {
       comp := v[root]
       for _, u := range comp {
           if d[u] != uf.sz[root]-1 {
               // pick node with minimal degree
               mn := comp[0]
               for _, x := range comp {
                   if d[x] < d[mn] {
                       mn = x
                   }
               }
               fmt.Fprintln(writer, 1)
               fmt.Fprintln(writer, mn)
               return
           }
       }
   }
   // two components
   if len(roots) == 2 {
       // choose smaller component
       a, b := roots[0], roots[1]
       if len(v[a]) > len(v[b]) {
           a = roots[1]
       }
       comp := v[a]
       fmt.Fprintln(writer, len(comp))
       for i, x := range comp {
           if i > 0 {
               writer.WriteByte(' ')
           }
           fmt.Fprint(writer, x)
       }
       fmt.Fprintln(writer)
       return
   }
   // three or more components
   // output two nodes from first two components
   fmt.Fprintln(writer, 2)
   fmt.Fprint(writer, v[roots[0]][0], " ", v[roots[1]][0], "\n")
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   var t int
   fmt.Fscan(reader, &t)
   for t > 0 {
       solve(reader, writer)
       t--
   }
}
