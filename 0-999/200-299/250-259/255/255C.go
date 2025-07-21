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

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   b := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &b[i])
   }
   // Compress values
   comp := make(map[int]int, n)
   vals := make([]int, 0, n)
   for _, v := range b {
       if _, ok := comp[v]; !ok {
           comp[v] = len(vals)
           vals = append(vals, v)
       }
   }
   m := len(vals)
   pos := make([][]int, m)
   for i, v := range b {
       id := comp[v]
       pos[id] = append(pos[id], i)
   }
   // Answer at least max freq of any single value (q=0)
   ans := 0
   for _, lst := range pos {
       if len(lst) > ans {
           ans = len(lst)
       }
   }
   // For each ordered pair of distinct values, find longest alternating subsequence
   for i := 0; i < m; i++ {
       pi0 := pos[i]
       for j := 0; j < m; j++ {
           if j == i {
               continue
           }
           pj0 := pos[j]
           // merge for pattern i,j,i,j...
           pi, pj := 0, 0
           lastPos := -1
           lastIsI := true
           cnt := 0
           for {
               if lastIsI {
                   // pick from pi0
                   for pi < len(pi0) && pi0[pi] <= lastPos {
                       pi++
                   }
                   if pi >= len(pi0) {
                       break
                   }
                   lastPos = pi0[pi]
                   pi++
                   cnt++
                   lastIsI = false
               } else {
                   // pick from pj0
                   for pj < len(pj0) && pj0[pj] <= lastPos {
                       pj++
                   }
                   if pj >= len(pj0) {
                       break
                   }
                   lastPos = pj0[pj]
                   pj++
                   cnt++
                   lastIsI = true
               }
           }
           if cnt > ans {
               ans = cnt
           }
       }
   }
   fmt.Fprintln(writer, ans)
}
