package main

import (
   "bufio"
   "fmt"
   "os"
)

type car struct {
   id int
   x  int
   l  int
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var L, b, f int
   if _, err := fmt.Fscan(reader, &L, &b, &f); err != nil {
       return
   }
   var n int
   fmt.Fscan(reader, &n)
   parked := make([]car, 0, n)
   // map request id -> parked car info
   info := make(map[int]car)

   for req := 1; req <= n; req++ {
       var t, v int
       fmt.Fscan(reader, &t, &v)
       if t == 1 {
           length := v
           placed := -1
           // try before first
           if len(parked) == 0 {
               if length <= L {
                   placed = 0
               }
           } else {
               // gap before first car
               first := parked[0]
               maxX := first.x - f - length
               if maxX >= 0 {
                   placed = 0
               }
           }
           // gaps between cars
           if placed < 0 {
               for i := 0; i+1 < len(parked); i++ {
                   cur := parked[i]
                   next := parked[i+1]
                   low := cur.x + cur.l + b
                   high := next.x - f - length
                   if high >= low {
                       placed = low
                       break
                   }
               }
           }
           // after last
           if placed < 0 {
               if len(parked) == 0 {
                   // already handled
               } else {
                   last := parked[len(parked)-1]
                   low := last.x + last.l + b
                   if low+length <= L {
                       placed = low
                   }
               }
           }
           // if no cars and not placed (length > L)
           if placed < 0 {
               fmt.Fprintln(writer, -1)
               // record as unplaced
               info[req] = car{id: req, x: -1, l: length}
               continue
           }
           // place car at placed
           c := car{id: req, x: placed, l: length}
           // insert into parked keeping sorted by x
           idx := 0
           for idx < len(parked) && parked[idx].x < placed {
               idx++
           }
           parked = append(parked, car{})
           copy(parked[idx+1:], parked[idx:])
           parked[idx] = c
           info[req] = c
           fmt.Fprintln(writer, placed)
       } else if t == 2 {
           remID := v
           c, ok := info[remID]
           if !ok || c.x < 0 {
               // nothing
               continue
           }
           // remove from parked
           for i, pc := range parked {
               if pc.id == remID {
                   parked = append(parked[:i], parked[i+1:]...)
                   break
               }
           }
           // mark removed
           info[remID] = car{id: remID, x: -1, l: c.l}
       }
   }
}
