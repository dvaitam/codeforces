package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var t int
   if _, err := fmt.Fscan(reader, &t); err != nil {
       return
   }
   // scoring points for positions 1..10
   ptsAward := []int{25, 18, 15, 12, 10, 8, 6, 4, 2, 1}
   const maxPos = 50
   ptsMap := make(map[string]int)
   posMap := make(map[string][]int)

   for i := 0; i < t; i++ {
       var n int
       fmt.Fscan(reader, &n)
       for pos := 1; pos <= n; pos++ {
           var name string
           fmt.Fscan(reader, &name)
           if _, ok := posMap[name]; !ok {
               posMap[name] = make([]int, maxPos)
               ptsMap[name] = 0
           }
           // record position count
           if pos-1 < maxPos {
               posMap[name][pos-1]++
           }
           // record points
           if pos <= len(ptsAward) {
               ptsMap[name] += ptsAward[pos-1]
           }
       }
   }
   // collect drivers
   names := make([]string, 0, len(posMap))
   for name := range posMap {
       names = append(names, name)
   }

   // original scoring: by total points, then by position counts from first, second, ...
   origNames := make([]string, len(names))
   copy(origNames, names)
   sort.Slice(origNames, func(i, j int) bool {
       a, b := origNames[i], origNames[j]
       if ptsMap[a] != ptsMap[b] {
           return ptsMap[a] > ptsMap[b]
       }
       // tie-break by most wins, seconds, etc.
       pa, pb := posMap[a], posMap[b]
       for k := 0; k < maxPos; k++ {
           if pa[k] != pb[k] {
               return pa[k] > pb[k]
           }
       }
       return false
   })

   // alternative scoring: by wins, then by points, then by 2nd, 3rd, ...
   altNames := make([]string, len(names))
   copy(altNames, names)
   sort.Slice(altNames, func(i, j int) bool {
       a, b := altNames[i], altNames[j]
       pa, pb := posMap[a], posMap[b]
       if pa[0] != pb[0] {
           return pa[0] > pb[0]
       }
       if ptsMap[a] != ptsMap[b] {
           return ptsMap[a] > ptsMap[b]
       }
       for k := 1; k < maxPos; k++ {
           if pa[k] != pb[k] {
               return pa[k] > pb[k]
           }
       }
       return false
   })

   // output champions
   if len(origNames) > 0 {
       fmt.Println(origNames[0])
   }
   if len(altNames) > 0 {
       fmt.Println(altNames[0])
   }
}
