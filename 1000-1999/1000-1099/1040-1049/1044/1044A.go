package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, m int
   if _, err := fmt.Fscan(in, &n, &m); err != nil {
       return
   }
   // dhMap: change in horizontal cover count; dvMap: change in vertical count
   dhMap := make(map[int]int)
   dvMap := make(map[int]int)
   const MAXX = 1000000000
   // vertical spells
   for i := 0; i < n; i++ {
       var x int
       fmt.Fscan(in, &x)
       // at x+1, number of verticals < x increases by 1
       if x+1 <= MAXX {
           dvMap[x+1]++
       }
       // x+1 > MAXX, ignore since beyond domain
   }
   // horizontal spells
   for i := 0; i < m; i++ {
       var x1, x2, y int
       fmt.Fscan(in, &x1, &x2, &y)
       // cover from x1 to x2 inclusive
       dhMap[x1]++
       if x2+1 <= MAXX {
           dhMap[x2+1]--
       }
   }
   // collect all event keys and ensure x=1 included
   keys := make([]int, 0, len(dhMap)+len(dvMap)+1)
   seen := make(map[int]struct{}, len(dhMap)+len(dvMap)+1)
   addKey := func(x int) {
       if x < 1 || x > MAXX {
           return
       }
       if _, ok := seen[x]; !ok {
           seen[x] = struct{}{}
           keys = append(keys, x)
       }
   }
   addKey(1)
   for x := range dhMap {
       addKey(x)
   }
   for x := range dvMap {
       addKey(x)
   }
   sort.Ints(keys)
   // sweep
   h := 0
   v := 0
   ans := n + m
   for _, x := range keys {
       if delta, ok := dhMap[x]; ok {
           h += delta
       }
       if delta, ok := dvMap[x]; ok {
           v += delta
       }
       // cost at column x: horizontal spells covering x + vertical spells crossed
       cost := h + v
       if cost < ans {
           ans = cost
       }
   }
   // output result
   w := bufio.NewWriter(os.Stdout)
   defer w.Flush()
   fmt.Fprintln(w, ans)
}
