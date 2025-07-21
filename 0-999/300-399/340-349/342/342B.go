package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, m, s, f int
   if _, err := fmt.Fscan(reader, &n, &m, &s, &f); err != nil {
       return
   }
   type event struct{ t, l, r int }
   events := make([]event, m)
   for i := 0; i < m; i++ {
       fmt.Fscan(reader, &events[i].t, &events[i].l, &events[i].r)
   }
   dir := 1
   if s > f {
       dir = -1
   }
   cur := s
   idx := 0
   t := 1
   // use a byte slice to build result
   res := make([]byte, 0, abs(f-s)+m+5)
   for cur != f {
       if idx < m && events[idx].t == t {
           L, R := events[idx].l, events[idx].r
           // if current or next spy is watched, wait
           next := cur + dir
           if (cur >= L && cur <= R) || (next >= L && next <= R) {
               res = append(res, 'X')
           } else if dir == 1 {
               cur++
               res = append(res, 'R')
           } else {
               cur--
               res = append(res, 'L')
           }
           idx++
       } else {
           if dir == 1 {
               cur++
               res = append(res, 'R')
           } else {
               cur--
               res = append(res, 'L')
           }
       }
       t++
   }
   writer.Write(res)
}

func abs(x int) int {
   if x < 0 {
       return -x
   }
   return x
}
