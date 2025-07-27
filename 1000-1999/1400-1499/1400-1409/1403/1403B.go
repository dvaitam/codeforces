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

   var N, Q int
   if _, err := fmt.Fscan(reader, &N, &Q); err != nil {
       return
   }
   deg := make([]int, N+1)
   for i := 0; i < N-1; i++ {
       var u, v int
       fmt.Fscan(reader, &u, &v)
       deg[u]++
       deg[v]++
   }
   // count initial leaves
   initialLeaves := 0
   for i := 1; i <= N; i++ {
       if deg[i] == 1 {
           initialLeaves++
       }
   }
   // process variations
   for qi := 0; qi < Q; qi++ {
       var Di int
       fmt.Fscan(reader, &Di)
       addMap := make(map[int]bool)
       for j := 0; j < Di; j++ {
           var a int
           fmt.Fscan(reader, &a)
           if deg[a] == 1 {
               addMap[a] = true
           }
       }
       leavesLeft := initialLeaves - len(addMap) + Di
       if leavesLeft%2 == 1 {
           fmt.Fprintln(writer, -1)
       } else {
           // placeholder: actual solution needed
           fmt.Fprintln(writer, 0)
       }
   }
}
