package main

import (
   "bufio"
   "fmt"
   "os"
)

// UnionFind implements union-find with size
type UnionFind struct {
   p []int
}

// NewUnionFind creates a new UnionFind of given size
func NewUnionFind(n int) *UnionFind {
   p := make([]int, n)
   for i := range p {
       p[i] = -1
   }
   return &UnionFind{p: p}
}

// Find returns the root of x
func (u *UnionFind) Find(x int) int {
   if u.p[x] < 0 {
       return x
   }
   u.p[x] = u.Find(u.p[x])
   return u.p[x]
}

// Union joins sets of x and y, returns true if merged
func (u *UnionFind) Union(x, y int) bool {
   x = u.Find(x)
   y = u.Find(y)
   if x == y {
       return false
   }
   // union by size: p[x] and p[y] are negative sizes
   if u.p[x] > u.p[y] {
       x, y = y, x
   }
   u.p[x] += u.p[y]
   u.p[y] = x
   return true
}

// Same reports whether x and y are in same set
func (u *UnionFind) Same(x, y int) bool {
   return u.Find(x) == u.Find(y)
}

// Size returns the size of the set containing x
func (u *UnionFind) Size(x int) int {
   r := u.Find(x)
   return -u.p[r]
}

func abs(a int) int {
   if a < 0 {
       return -a
   }
   return a
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, m int
   if _, err := fmt.Fscan(in, &n, &m); err != nil {
       return
   }
   e := make([][]bool, n)
   for i := 0; i < n; i++ {
       e[i] = make([]bool, n)
   }
   for i := 0; i < m; i++ {
       var a, b int
       fmt.Fscan(in, &a, &b)
       a--
       b--
       e[a][b] = true
       e[b][a] = true
   }
   uf := NewUnionFind(2 * n)
   for i := 0; i < n; i++ {
       for j := 0; j < n; j++ {
           if i == j || e[i][j] {
               continue
           }
           uf.Union(i, j+n)
           uf.Union(i+n, j)
       }
   }
   for i := 0; i < n; i++ {
       if uf.Same(i, i+n) {
           fmt.Println("No")
           return
       }
   }
   ans := make([]byte, n)
   for i := 0; i < n; i++ {
       ans[i] = 'b'
   }
   for i := 0; i < n; i++ {
       if uf.Size(i) == 1 {
           continue
       }
       par := uf.Find(i)
       if par >= n {
           ans[i] = 'c'
       } else {
           ans[i] = 'a'
       }
   }
   ok := true
   for i := 0; i < n; i++ {
       for j := 0; j < n; j++ {
           if i == j {
               continue
           }
           df := abs(int(ans[i]) - int(ans[j]))
           if df <= 1 && !e[i][j] {
               ok = false
           }
           if df >= 2 && e[i][j] {
               ok = false
           }
       }
   }
   if !ok {
       fmt.Println("No")
   } else {
       fmt.Println("Yes")
       fmt.Println(string(ans))
   }
}
