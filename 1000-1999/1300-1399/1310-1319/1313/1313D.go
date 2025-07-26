package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

type event struct {
   x   int
   typ int // 0 = remove, 1 = add
   id  int
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   var m, k int
   if _, err := fmt.Fscan(in, &n, &m, &k); err != nil {
       return
   }
   L := make([]int, n)
   R := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &L[i], &R[i])
   }
   events := make([]event, 0, 2*n)
   for i := 0; i < n; i++ {
       events = append(events, event{L[i], 1, i})
       events = append(events, event{R[i] + 1, 0, i})
   }
   sort.Slice(events, func(i, j int) bool {
       if events[i].x != events[j].x {
           return events[i].x < events[j].x
       }
       return events[i].typ < events[j].typ
   })
   // precompute parity for masks up to 1<<k
   maxMask := 1 << k
   parity := make([]int, maxMask)
   for mask := 1; mask < maxMask; mask++ {
       parity[mask] = parity[mask>>1] ^ (mask & 1)
   }
   // DP state
   f := make([]int64, 1)
   // active intervals
   active := make([]int, 0, k)
   id2pos := make([]int, n)
   prevX := events[0].x
   nev := len(events)
   i := 0
   for i < nev {
       x := events[i].x
       // segment [prevX, x)
       segLen := x - prevX
       if segLen > 0 {
           d := int64(segLen)
           for mask := range f {
               if parity[mask] == 1 {
                   f[mask] += d
               }
           }
       }
       // process events at x
       j := i
       for j < nev && events[j].x == x {
           j++
       }
       // removals first
       for t := i; t < j; t++ {
           if events[t].typ == 0 {
               id := events[t].id
               // remove interval id
               pos := id2pos[id]
               sz := len(active)
               // merge DP states by removing bit at pos
               newSz := sz - 1
               newF := make([]int64, 1<<newSz)
               for mask, val := range f {
                   // new mask without bit pos
                   low := mask & ((1 << pos) - 1)
                   high := (mask >> 1) & ^((1 << pos) - 1)
                   newMask := low | high
                   if val > newF[newMask] {
                       newF[newMask] = val
                   }
               }
               f = newF
               // update active list and positions
               active = append(active[:pos], active[pos+1:]...)
               for idx := pos; idx < len(active); idx++ {
                   id2pos[active[idx]] = idx
               }
           }
       }
       // additions
       for t := i; t < j; t++ {
           if events[t].typ == 1 {
               id := events[t].id
               // add interval id
               sz := len(active)
               // expand DP states by adding new bit at position sz
               newF := make([]int64, len(f)*2)
               for mask, val := range f {
                   newF[mask] = val
                   newF[mask|(1<<sz)] = val
               }
               f = newF
               // update active list and positions
               id2pos[id] = sz
               active = append(active, id)
           }
       }
       prevX = x
       i = j
   }
   // result
   var ans int64
   for _, v := range f {
       if v > ans {
           ans = v
       }
   }
   fmt.Println(ans)
}
