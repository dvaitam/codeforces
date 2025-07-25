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
   size := make([]int, n+1)
   for i := 1; i <= n; i++ {
       parent[i] = i
       size[i] = 1
   }

   for i := 0; i < m; i++ {
       var k int
       fmt.Fscan(reader, &k)
       if k == 0 {
           continue
       }
       users := make([]int, k)
       for j := 0; j < k; j++ {
           fmt.Fscan(reader, &users[j])
       }
       first := users[0]
       for j := 1; j < k; j++ {
           union(parent, size, first, users[j])
       }
   }

   for i := 1; i <= n; i++ {
       root := find(parent, i)
       fmt.Fprint(writer, size[root])
       if i < n {
           fmt.Fprint(writer, " ")
       }
   }
   fmt.Fprintln(writer)
}

func find(parent []int, x int) int {
   if parent[x] != x {
       parent[x] = find(parent, parent[x])
   }
   return parent[x]
}

func union(parent []int, size []int, a, b int) {
   a = find(parent, a)
   b = find(parent, b)
   if a == b {
       return
   }
   if size[a] < size[b] {
       a, b = b, a
   }
   parent[b] = a
   size[a] += size[b]
}
