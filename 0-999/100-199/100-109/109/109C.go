package main

import (
   "bufio"
   "fmt"
   "os"
)

// UF is a Union-Find (disjoint set) structure
type UF struct {
   parent []int
   size   []int
}

// NewUF initializes a Union-Find for n elements (0..n-1)
func NewUF(n int) *UF {
   parent := make([]int, n)
   size := make([]int, n)
   for i := 0; i < n; i++ {
       parent[i] = i
       size[i] = 1
   }
   return &UF{parent: parent, size: size}
}

// Find returns the representative of x with path compression
func (u *UF) Find(x int) int {
   if u.parent[x] != x {
       u.parent[x] = u.Find(u.parent[x])
   }
   return u.parent[x]
}

// Union merges the sets containing a and b
func (u *UF) Union(a, b int) {
   a = u.Find(a)
   b = u.Find(b)
   if a == b {
       return
   }
   // attach smaller tree to larger
   if u.size[a] > u.size[b] {
       a, b = b, a
   }
   u.parent[a] = b
   u.size[b] += u.size[a]
   u.size[a] = 0
}

// lucky returns true if all digits of x are 4 or 7 (or x is 0)
func lucky(x int) bool {
   if x == 0 {
       return true
   }
   for x > 0 {
       d := x % 10
       if d != 4 && d != 7 {
           return false
       }
       x /= 10
   }
   return true
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   uf := NewUF(n)
   // read edges
   for i := 0; i < n-1; i++ {
       var a, b, c int
       fmt.Fscan(reader, &a, &b, &c)
       // input is 1-based; convert to 0-based
       if !lucky(c) {
           uf.Union(a-1, b-1)
       }
   }
   // count component sizes
   // uf.size at roots contains size; non-roots have size 0
   total := int64(n)
   var ans int64
   for i := 0; i < n; i++ {
       s := uf.size[i]
       if s == 0 {
           continue
       }
       rem := total - int64(s)
       // for each node in component, choose two distinct outside nodes
       ans += int64(s) * rem * (rem - 1)
   }
   fmt.Println(ans)
}
