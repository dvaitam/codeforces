package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

// event represents a possible date interval for an event
type event struct {
   l, r, idx int
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   events := make([]event, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &events[i].l, &events[i].r)
       events[i].idx = i
   }
   // Sort by earliest finishing time
   sort.Slice(events, func(i, j int) bool {
       if events[i].r != events[j].r {
           return events[i].r < events[j].r
       }
       return events[i].l < events[j].l
   })
   used := make(map[int]bool)
   ans := make([]int, n)
   // Greedily assign the smallest available day in [l, r]
   for _, e := range events {
       day := e.l
       for used[day] {
           day++
       }
       used[day] = true
       ans[e.idx] = day
   }
   // Output in original order
   for i, v := range ans {
       if i > 0 {
           fmt.Fprint(writer, " ")
       }
       fmt.Fprint(writer, v)
   }
   fmt.Fprint(writer, "\n")
}
