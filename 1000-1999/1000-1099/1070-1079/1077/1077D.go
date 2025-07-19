package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

// Node holds a number and its frequency
type Node struct {
   x, w int
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, k int
   if _, err := fmt.Fscan(reader, &n, &k); err != nil {
       return
   }
   values := make([]int, n)
   maxVal := 0
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &values[i])
       if values[i] > maxVal {
           maxVal = values[i]
       }
   }
   // count frequencies
   freq := make([]int, maxVal+1)
   for _, v := range values {
       freq[v]++
   }
   // build nodes for non-zero frequencies
   nodes := make([]Node, 0, len(freq))
   for i, w := range freq {
       if i > 0 && w > 0 {
           nodes = append(nodes, Node{x: i, w: w})
       }
   }
   // sort descending by frequency
   sort.Slice(nodes, func(i, j int) bool {
       return nodes[i].w > nodes[j].w
   })
   // binary search maximum mid such that sum(w/mid) >= k
   l, r := 1, 0
   if len(nodes) > 0 {
       r = nodes[0].w
   }
   for l <= r {
       mid := (l + r) / 2
       sum := 0
       for _, node := range nodes {
           sum += node.w / mid
           if sum >= k {
               break
           }
       }
       if sum >= k {
           l = mid + 1
       } else {
           r = mid - 1
       }
   }
   // output k numbers, each x repeated floor(w/r) times
   cnt := 0
   for _, node := range nodes {
       times := node.w / r
       for j := 0; j < times && cnt < k; j++ {
           fmt.Fprint(writer, node.x, " ")
           cnt++
       }
       if cnt >= k {
           break
       }
   }
   fmt.Fprintln(writer)
}
