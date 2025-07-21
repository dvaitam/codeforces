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

   var n int
   var k int64
   if _, err := fmt.Fscan(reader, &n, &k); err != nil {
       return
   }
   a := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   // total examinations
   var total int64
   for _, v := range a {
       total += v
   }
   if total < k {
       fmt.Fprintln(writer, -1)
       return
   }
   // binary search maximum t such that sum(min(a[i], t)) <= k
   var low, high int64 = 0, 0
   for _, v := range a {
       if v > high {
           high = v
       }
   }
   for low < high {
       mid := (low + high + 1) >> 1
       var s int64
       for _, v := range a {
           if v < mid {
               s += v
           } else {
               s += mid
           }
           if s > k {
               break
           }
       }
       if s <= k {
           low = mid
       } else {
           high = mid - 1
       }
   }
   t := low
   // sum of done exams after t rounds
   var done int64
   for _, v := range a {
       if v < t {
           done += v
       } else {
           done += t
       }
   }
   rem := k - done
   // build queue of animals still needing exams
   type animal struct { idx int; rem int64 }
   queue := make([]animal, 0, n)
   for i, v := range a {
       if v > t {
           queue = append(queue, animal{idx: i + 1, rem: v - t})
       }
   }
   // simulate remaining rem examinations
   cur := 0
   for rem > 0 && cur < len(queue) {
       // one examination for queue[cur]
       queue[cur].rem--
       rem--
       cur++
   }
   // output remaining queue starting at cur
   out := make([]int, 0, len(queue))
   for i := cur; i < len(queue); i++ {
       if queue[i].rem > 0 {
           out = append(out, queue[i].idx)
       }
   }
   for i := 0; i < cur; i++ {
       if queue[i].rem > 0 {
           out = append(out, queue[i].idx)
       }
   }
   // print result
   for i, v := range out {
       if i > 0 {
           writer.WriteByte(' ')
       }
       fmt.Fprint(writer, v)
   }
   // newline
   if len(out) > 0 {
       writer.WriteByte('\n')
   }
}
