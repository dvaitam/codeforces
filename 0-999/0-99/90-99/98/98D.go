package main

import (
   "bufio"
   "fmt"
   "os"
)

type pair struct { x, y int }

var (
   n int
   a []int
   b []pair
)

// dfs0 computes the minimal number of moves without recording them
func dfs0(m, x, y, z int) int {
   if m == 0 {
       return 0
   }
   // find last index i where a[i] != a[m]
   i := m
   for i > 0 && a[i] == a[m] {
       i--
   }
   if z != 0 && i+1 < m {
       t1 := dfs0(i, x, y, 0) + (m-i) + dfs0(i, y, x, 0) + (m-i) + dfs0(i, x, y, 1)
       t2 := dfs0(n-1, x, 6-x-y, 0) + 1 + dfs0(n-1, 6-x-y, y, 0)
       if t1 < t2 {
           return t1
       }
       return t2
   }
   // single operation
   return dfs0(i, x, 6-x-y, 0) + (m - i) + dfs0(i, 6-x-y, y, 0)
}

// dfs computes the minimal moves and records them in b
func dfs(m, x, y, z int) int {
   if m == 0 {
       return 0
   }
   i := m
   for i > 0 && a[i] == a[m] {
       i--
   }
   if z != 0 && i+1 < m {
       t1 := dfs0(i, x, y, 0) + (m-i) + dfs0(i, y, x, 0) + (m-i) + dfs0(i, x, y, 1)
       t2 := dfs0(n-1, x, 6-x-y, 0) + 1 + dfs0(n-1, 6-x-y, y, 0)
       if t1 < t2 {
           dfs(i, x, y, 0)
           for j := m; j > i; j-- {
               b = append(b, pair{x, 6 - x - y})
           }
           dfs(i, y, x, 0)
           for j := m; j > i; j-- {
               b = append(b, pair{6 - x - y, y})
           }
           dfs(i, x, y, 1)
           return t1
       }
       // use alternate sequence
       dfs(n-1, x, 6-x-y, 0)
       b = append(b, pair{x, y})
       dfs(n-1, 6-x-y, y, 0)
       return t2
   }
   // single split
   t := dfs(i, x, 6-x-y, 0) + (m - i)
   for j := m; j > i; j-- {
       b = append(b, pair{x, y})
   }
   t += dfs(i, 6-x-y, y, 0)
   return t
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fscan(reader, &n)
   a = make([]int, n+1)
   // read input in reverse order to match original logic
   for i := n; i >= 1; i-- {
       fmt.Fscan(reader, &a[i])
   }
   // compute moves
   cnt := dfs(n, 1, 3, 1)
   fmt.Fprintln(writer, cnt)
   for _, p := range b {
       fmt.Fprintln(writer, p.x, p.y)
   }
}
