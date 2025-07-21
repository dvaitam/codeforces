package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   used := make(map[string]int, n)
   for i := 0; i < n; i++ {
       var name string
       var year int
       fmt.Fscan(reader, &name, &year)
       if prev, ok := used[name]; !ok || year > prev {
           used[name] = year
       }
   }
   var m int
   fmt.Fscan(reader, &m)
   var unusedBest string
   haveUnused := false
   var usedBest string
   var usedYearBest int
   haveUsed := false
   for i := 0; i < m; i++ {
       var cand string
       fmt.Fscan(reader, &cand)
       if year, ok := used[cand]; !ok {
           if !haveUnused || cand > unusedBest {
               unusedBest = cand
               haveUnused = true
           }
       } else {
           if !haveUsed || year < usedYearBest || (year == usedYearBest && cand > usedBest) {
               usedYearBest = year
               usedBest = cand
               haveUsed = true
           }
       }
   }
   if haveUnused {
       fmt.Println(unusedBest)
   } else {
       fmt.Println(usedBest)
   }
}
