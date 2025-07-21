package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   p := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &p[i])
       // convert to 0-based value if needed, but values kept as is
   }
   // build adjacency list
   adj := make([][]int, n)
   for i := 0; i < n; i++ {
       var line string
       fmt.Fscan(reader, &line)
       for j := 0; j < n; j++ {
           if line[j] == '1' {
               adj[i] = append(adj[i], j)
           }
       }
   }
   visited := make([]bool, n)
   // process each connected component
   for i := 0; i < n; i++ {
       if visited[i] {
           continue
       }
       // BFS to collect component
       queue := []int{i}
       visited[i] = true
       compIdx := []int{i}
       for qi := 0; qi < len(queue); qi++ {
           u := queue[qi]
           for _, v := range adj[u] {
               if !visited[v] {
                   visited[v] = true
                   queue = append(queue, v)
                   compIdx = append(compIdx, v)
               }
           }
       }
       // collect values
       compVals := make([]int, len(compIdx))
       for k, idx := range compIdx {
           compVals[k] = p[idx]
       }
       // sort indices and values
       sort.Ints(compIdx)
       sort.Ints(compVals)
       // assign smallest values to smallest positions
       for k, idx := range compIdx {
           p[idx] = compVals[k]
       }
   }
   // output
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   for i, v := range p {
       if i > 0 {
           writer.WriteByte(' ')
       }
       fmt.Fprint(writer, v)
   }
   writer.WriteByte('\n')
}
