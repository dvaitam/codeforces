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
   ans := make([]int, n+2)
   parent := make([]int, n+2)
   for i := 1; i <= n+1; i++ {
       parent[i] = i
   }

   // find with path halving
   var find func(int) int
   find = func(x int) int {
       for parent[x] != x {
           parent[x] = parent[parent[x]]
           x = parent[x]
       }
       return x
   }
   // union a to b
   union := func(a, b int) {
       pa := find(a)
       pb := find(b)
       parent[pa] = pb
   }

   for i := 0; i < m; i++ {
       var l, r, x int
       fmt.Fscan(reader, &l, &r, &x)
       j := find(l)
       for j <= r {
           if j == x {
               j = find(j + 1)
               continue
           }
           ans[j] = x
           union(j, j+1)
           j = find(j)
       }
   }

   // output for knights 1..n
   for i := 1; i <= n; i++ {
       if i > 1 {
           writer.WriteByte(' ')
       }
       writer.WriteString(fmt.Sprint(ans[i]))
   }
   writer.WriteByte('\n')
}
