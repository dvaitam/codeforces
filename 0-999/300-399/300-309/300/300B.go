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
   // read until EOF
   for {
       if _, err := fmt.Fscan(reader, &n, &m); err != nil {
           return
       }
       // initialize DSU
       parent := make([]int, n+1)
       size := make([]int, n+1)
       for i := 1; i <= n; i++ {
           parent[i] = i
           size[i] = 1
       }
       // find and union functions
       var find func(int) int
       find = func(x int) int {
           if parent[x] != x {
               parent[x] = find(parent[x])
           }
           return parent[x]
       }
       unite := func(x, y int) {
           fx := find(x)
           fy := find(y)
           if fx != fy {
               parent[fx] = fy
               size[fy] += size[fx]
           }
       }
       // process edges
       for i := 0; i < m; i++ {
           var a, b int
           fmt.Fscan(reader, &a, &b)
           unite(a, b)
       }
       // collect components by size
       var singles, pairs []int
       bad := false
       for i := 1; i <= n; i++ {
           if parent[i] == i {
               sz := size[i]
               if sz > 3 {
                   bad = true
                   break
               }
               if sz == 1 {
                   singles = append(singles, i)
               } else if sz == 2 {
                   pairs = append(pairs, i)
               }
           }
       }
       // check feasibility
       if bad || len(singles) < len(pairs) || (len(singles) > len(pairs) && (len(singles)-len(pairs))%3 != 0) {
           fmt.Fprintln(writer, -1)
           continue
       }
       // merge to form groups of three
       if len(singles) == len(pairs) {
           for i := range pairs {
               unite(pairs[i], singles[i])
           }
       } else {
           // more singles: group some singles into triplets
           diff := len(singles) - len(pairs)
           triplets := diff / 3
           idx := 0
           for i := 0; i < triplets; i++ {
               // merge three singles
               unite(singles[idx], singles[idx+1])
               unite(singles[idx+1], singles[idx+2])
               idx += 3
           }
           // merge each pair with one remaining single
           for i := range pairs {
               unite(pairs[i], singles[idx])
               idx++
           }
       }
       // collect final groups
       groups := make(map[int][]int, n/3)
       for i := 1; i <= n; i++ {
           root := find(i)
           groups[root] = append(groups[root], i)
       }
       // output groups
       for _, grp := range groups {
           for j, v := range grp {
               if j > 0 {
                   writer.WriteByte(' ')
               }
               fmt.Fprint(writer, v)
           }
           writer.WriteByte('\n')
       }
   }
}
