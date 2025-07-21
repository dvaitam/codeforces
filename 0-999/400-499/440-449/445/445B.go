package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, m int
   fmt.Fscan(reader, &n, &m)
   parent := make([]int, n+1)
   for i := 1; i <= n; i++ {
      parent[i] = i
   }
   var x, y int
   for i := 0; i < m; i++ {
      fmt.Fscan(reader, &x, &y)
      union(parent, x, y)
   }
   roots := make(map[int]bool)
   for i := 1; i <= n; i++ {
      roots[find(parent, i)] = true
   }
   comps := len(roots)
   exp := n - comps
   var ans uint64 = 1
   ans <<= exp
   fmt.Fprintln(writer, ans)
}

// find returns the root of x with path compression
func find(parent []int, x int) int {
   if parent[x] != x {
      parent[x] = find(parent, parent[x])
   }
   return parent[x]
}

// union merges the sets containing a and b
func union(parent []int, a, b int) {
   ra := find(parent, a)
   rb := find(parent, b)
   if ra != rb {
      parent[rb] = ra
   }
}
