package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, k int
   fmt.Fscan(reader, &n, &k)
   type team struct{ problems, time int }
   teams := make([]team, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &teams[i].problems, &teams[i].time)
   }
   sort.Slice(teams, func(i, j int) bool {
       if teams[i].problems != teams[j].problems {
           return teams[i].problems > teams[j].problems
       }
       return teams[i].time < teams[j].time
   })
   // Find score of k-th place (1-indexed after sorting)
   tp := teams[k-1].problems
   tt := teams[k-1].time
   count := 0
   for _, t := range teams {
       if t.problems == tp && t.time == tt {
           count++
       }
   }
   fmt.Println(count)
}
