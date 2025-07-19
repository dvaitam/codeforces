package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

// Node represents an element with original id and its transformed value
type Node struct {
   id, v int
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, h, m, k int
   fmt.Fscan(reader, &n, &h, &m, &k)
   // halve m as per original logic
   m /= 2
   // prepare duplicated nodes
   nodes := make([]Node, 2*n)
   for i := 0; i < n; i++ {
       var hi, vi int
       fmt.Fscan(reader, &hi, &vi)
       v := vi
       if v >= m {
           v -= m
       }
       nodes[i] = Node{id: i + 1, v: v}
       nodes[i+n] = Node{id: i + 1, v: v + m}
   }
   // sort by value
   sort.Slice(nodes, func(i, j int) bool {
       return nodes[i].v < nodes[j].v
   })
   // two-pointer to find minimal window
   bestLen := 2*n + 1
   bestJ, bestI := 0, 0
   j := 0
   for i := n; i < 2*n; i++ {
       // move j to maintain nodes[j].v >= nodes[i].v - k + 1
       for j < 2*n && nodes[j].v < nodes[i].v - k + 1 {
           j++
       }
       // window [j, i) has size i-j
       currLen := i - j
       if currLen < bestLen {
           bestLen = currLen
           bestJ = j
           bestI = i
       }
   }
   // output: number of elements and the threshold
   fmt.Fprint(writer, bestLen, " ", nodes[bestI].v - m, "\n")
   // output ids in the optimal window
   for idx := bestJ; idx < bestI; idx++ {
       fmt.Fprint(writer, nodes[idx].id, " ")
   }
}
