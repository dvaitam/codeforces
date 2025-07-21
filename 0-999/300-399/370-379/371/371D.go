package main

import (
   "bufio"
   "fmt"
   "os"
)

var (
   n, m int
   capv []int64
   curr []int64
   parent []int
   reader *bufio.Reader
   writer *bufio.Writer
)

func find(x int) int {
   if parent[x] != x {
       parent[x] = find(parent[x])
   }
   return parent[x]
}

func union(x, y int) {
   rx := find(x)
   ry := find(y)
   parent[rx] = ry
}

func main() {
   reader = bufio.NewReader(os.Stdin)
   writer = bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   // Read n
   fmt.Fscan(reader, &n)
   capv = make([]int64, n+2)
   curr = make([]int64, n+2)
   parent = make([]int, n+2)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &capv[i])
       parent[i] = i
   }
   parent[n+1] = n + 1

   // Read queries
   fmt.Fscan(reader, &m)
   for i := 0; i < m; i++ {
       var tp int
       fmt.Fscan(reader, &tp)
       if tp == 1 {
           var p int
           var x int64
           fmt.Fscan(reader, &p, &x)
           // pour x liters starting at p
           idx := find(p)
           for idx <= n && x > 0 {
               // fill vessel idx
               space := capv[idx] - curr[idx]
               if space > x {
                   curr[idx] += x
                   x = 0
                   break
               }
               // fill and overflow
               curr[idx] = capv[idx]
               x -= space
               // skip this vessel
               union(idx, idx+1)
               idx = find(idx)
           }
       } else if tp == 2 {
           var k int
           fmt.Fscan(reader, &k)
           fmt.Fprintln(writer, curr[k])
       }
   }
}
