package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   var n int
   fmt.Fscan(reader, &n)
   posMap := make(map[int][]int)
   for i := 1; i <= n; i++ {
       var x int
       fmt.Fscan(reader, &x)
       posMap[x] = append(posMap[x], i)
   }
   keys := make([]int, 0, len(posMap))
   for k := range posMap {
       keys = append(keys, k)
   }
   sort.Ints(keys)
   type result struct { x, d int }
   results := make([]result, 0, len(keys))
   for _, x := range keys {
       positions := posMap[x]
       if len(positions) == 1 {
           results = append(results, result{x, 0})
       } else {
           d := positions[1] - positions[0]
           ok := true
           for i := 2; i < len(positions); i++ {
               if positions[i]-positions[i-1] != d {
                   ok = false
                   break
               }
           }
           if ok {
               results = append(results, result{x, d})
           }
       }
   }
   fmt.Fprintln(writer, len(results))
   for _, r := range results {
       fmt.Fprintln(writer, r.x, r.d)
   }
}
