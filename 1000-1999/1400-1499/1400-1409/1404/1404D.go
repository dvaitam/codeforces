package main

import (
   "bufio"
   "fmt"
   "os"
)

// Interactive solution for CF 1404D
func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   // If n is even, play as First
   if n%2 == 0 {
       fmt.Fprintln(writer, "First")
       for i := 1; i <= n; i++ {
           fmt.Fprintf(writer, "%d %d\n", i, i+n)
       }
       return
   }
   // n is odd: play as Second
   fmt.Fprintln(writer, "Second")
   writer.Flush()
   // read n pairs
   a := make([]int, n)
   b := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i], &b[i])
   }
   // build classes and graph
   classA := make([]int, n)
   classB := make([]int, n)
   graph := make([][]struct{ to, idx int }, n)
   for i := 0; i < n; i++ {
       u := (a[i]-1) % n
       v := (b[i]-1) % n
       classA[i] = u
       classB[i] = v
       graph[u] = append(graph[u], struct{ to, idx int }{v, i})
       graph[v] = append(graph[v], struct{ to, idx int }{u, i})
   }
   dir := make([]bool, n)       // true: pick b[i], false: pick a[i]
   usedEdge := make([]bool, n)
   cycles := make([][]int, 0, n)
   // find cycles in class-graph
   for start := 0; start < n; start++ {
       if len(graph[start]) == 0 {
           continue
       }
       // degree is 2 for all nodes
       // traverse unvisited edge
       for _, e0 := range graph[start] {
           if usedEdge[e0.idx] {
               continue
           }
           // start a new cycle
           var cycle []int
           curr := start
           prevEdge := -1
           for {
               // find next edge not prevEdge
               var e struct{ to, idx int }
               for _, ee := range graph[curr] {
                   if ee.idx != prevEdge {
                       e = ee
                       break
                   }
               }
               if usedEdge[e.idx] {
                   break
               }
               usedEdge[e.idx] = true
               cycle = append(cycle, e.idx)
               // orient edge from curr -> e.to
               p := e.idx
               if classA[p] == curr && classB[p] == e.to {
                   dir[p] = true
               } else if classB[p] == curr && classA[p] == e.to {
                   dir[p] = false
               }
               prevEdge = e.idx
               curr = e.to
               if curr == start {
                   break
               }
           }
           if len(cycle) > 0 {
               cycles = append(cycles, cycle)
           }
       }
   }
   // count picks > n
   t := 0
   for i := 0; i < n; i++ {
       if dir[i] {
           if b[i] > n {
               t++
           }
       } else {
           if a[i] > n {
               t++
           }
       }
   }
   // desired parity for t: (n+1)/2 %2
   kp := ((n + 1) / 2) & 1
   if (t & 1) != kp {
       // find odd-length cycle to flip
       for _, cycle := range cycles {
           if len(cycle)%2 == 1 {
               for _, p := range cycle {
                   dir[p] = !dir[p]
               }
               break
           }
       }
   }
   // output selections
   for i := 0; i < n; i++ {
       if dir[i] {
           fmt.Fprint(writer, b[i])
       } else {
           fmt.Fprint(writer, a[i])
       }
       if i+1 < n {
           fmt.Fprint(writer, " ")
       }
   }
   fmt.Fprintln(writer)
}
