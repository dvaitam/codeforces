package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
   "sort"
)

// event represents an opening, closing, or query event
type event struct {
   typ int // -1 open, 0 query, 1 close
   x   int
   val int
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   events := make([]event, 0, 2*n+ m)
   for i := 0; i < n; i++ {
       var a, h, l, r int
       fmt.Fscan(reader, &a, &h, &l, &r)
       // open segment [a-h, a-1] with threshold 100-l
       events = append(events, event{typ: -1, x: a - h, val: 100 - l})
       events = append(events, event{typ: 1, x: a - 1, val: 100 - l})
       // open segment [a+1, a+h] with threshold 100-r
       events = append(events, event{typ: -1, x: a + 1, val: 100 - r})
       events = append(events, event{typ: 1, x: a + h, val: 100 - r})
   }
   for i := 0; i < m; i++ {
       var b, z int
       fmt.Fscan(reader, &b, &z)
       events = append(events, event{typ: 0, x: b, val: z})
   }
   sort.Slice(events, func(i, j int) bool {
       if events[i].x != events[j].x {
           return events[i].x < events[j].x
       }
       return events[i].typ < events[j].typ
   })
   // prob[j] tracks count of active thresholds equal to j
   prob := make([]int, 101)
   var ans float64
   for _, ev := range events {
       switch ev.typ {
       case -1:
           if ev.val >= 0 && ev.val <= 100 {
               prob[ev.val]++
           }
       case 1:
           if ev.val >= 0 && ev.val <= 100 {
               prob[ev.val]--
           }
       case 0:
           // if any segment blocks with 100% threshold
           if prob[0] > 0 {
               continue
           }
           cur := float64(ev.val)
           for j := 1; j < 100; j++ {
               cnt := prob[j]
               if cnt > 0 {
                   cur *= math.Pow(float64(j)/100.0, float64(cnt))
               }
           }
           ans += cur
       }
   }
   // output with 12 decimal places
   fmt.Printf("%.12f", ans)
}
