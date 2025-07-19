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
   // map team id to list of values
   teamVals := make(map[int][]int)
   teams := make([]int, 0)
   for i := 0; i < n; i++ {
       var m, k int
       fmt.Fscan(reader, &m, &k)
       if _, ok := teamVals[m]; !ok {
           teams = append(teams, m)
       }
       teamVals[m] = append(teamVals[m], k)
   }

   var bestAvg float64
   var bestIds []int
   // try selecting t teams, t from 1 to 20
   for t := 1; t <= 20; t++ {
       // compute sum for each team
       type tsEntry struct{ sum, id int }
       entries := make([]tsEntry, 0, len(teams))
       for _, id := range teams {
           vals := teamVals[id]
           s := 0
           for _, v := range vals {
               if v >= t {
                   s += t
               } else {
                   s += v
               }
           }
           if s > 0 {
               entries = append(entries, tsEntry{s, id})
           }
       }
       if len(entries) < t {
           continue
       }
       sort.Slice(entries, func(i, j int) bool {
           return entries[i].sum > entries[j].sum
       })
       total := 0
       ids := make([]int, t)
       for i := 0; i < t; i++ {
           total += entries[i].sum
           ids[i] = entries[i].id
       }
       avg := float64(total) / float64(t)
       if avg > bestAvg {
           bestAvg = avg
           bestIds = ids
       }
   }

   fmt.Fprintln(writer, len(bestIds))
   for _, id := range bestIds {
       fmt.Fprint(writer, id, " ")
   }
   fmt.Fprintln(writer)
}
