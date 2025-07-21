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

   var n, k int
   if _, err := fmt.Fscan(reader, &n, &k); err != nil {
       return
   }
   d := make([]int, n+1)
   levels := make([][]int, n)
   maxD := 0
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &d[i])
       if d[i] < 0 || d[i] >= n {
           fmt.Fprintln(writer, -1)
           return
       }
       levels[d[i]] = append(levels[d[i]], i)
       if d[i] > maxD {
           maxD = d[i]
       }
   }
   // must have exactly one root at distance 0
   if len(levels[0]) != 1 {
       fmt.Fprintln(writer, -1)
       return
   }
   // check degree constraints per level
   for dist := 1; dist <= maxD; dist++ {
       cnt := len(levels[dist])
       if cnt == 0 {
           // skip empty level, but if further levels non-empty can't connect
           continue
       }
       parentCnt := len(levels[dist-1])
       var capacity int
       if dist == 1 {
           capacity = parentCnt * k
       } else {
           capacity = parentCnt * (k - 1)
       }
       if cnt > capacity {
           fmt.Fprintln(writer, -1)
           return
       }
   }
   // build edges
   type edge struct{ u, v int }
   edges := make([]edge, 0, n-1)
   // connect level 1 to root
   if maxD >= 1 {
       root := levels[0][0]
       for _, v := range levels[1] {
           edges = append(edges, edge{root, v})
       }
   }
   // connect deeper levels
   for dist := 2; dist <= maxD; dist++ {
       parents := levels[dist-1]
       children := levels[dist]
       if len(children) == 0 {
           continue
       }
       childCount := make([]int, len(parents))
       parentIdx := 0
       // each parent can have up to k-1 children at this level
       for _, v := range children {
           // find next available parent
           for parentIdx < len(parents) && childCount[parentIdx] >= (k-1) {
               parentIdx++
           }
           if parentIdx >= len(parents) {
               fmt.Fprintln(writer, -1)
               return
           }
           u := parents[parentIdx]
           edges = append(edges, edge{u, v})
           childCount[parentIdx]++
       }
   }
   // output
   m := len(edges)
   fmt.Fprintln(writer, m)
   for _, e := range edges {
       fmt.Fprintf(writer, "%d %d\n", e.u, e.v)
   }
}
