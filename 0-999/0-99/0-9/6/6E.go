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
   var k int
   fmt.Fscan(reader, &n, &k)
   h := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &h[i])
   }
   maxDeque := make([]int, 0, n)
   minDeque := make([]int, 0, n)
   l := 0
   best := 0
   type interval struct{ l, r int }
   res := make([]interval, 0)
   for r := 0; r < n; r++ {
       // add new element to max deque
       for len(maxDeque) > 0 && h[maxDeque[len(maxDeque)-1]] <= h[r] {
           maxDeque = maxDeque[:len(maxDeque)-1]
       }
       maxDeque = append(maxDeque, r)
       // add new element to min deque
       for len(minDeque) > 0 && h[minDeque[len(minDeque)-1]] >= h[r] {
           minDeque = minDeque[:len(minDeque)-1]
       }
       minDeque = append(minDeque, r)
       // shrink window if invalid
       for l <= r && h[maxDeque[0]]-h[minDeque[0]] > k {
           if maxDeque[0] == l {
               maxDeque = maxDeque[1:]
           }
           if minDeque[0] == l {
               minDeque = minDeque[1:]
           }
           l++
       }
       // current window [l..r]
       length := r - l + 1
       if length > best {
           best = length
           res = res[:0]
           res = append(res, interval{l, r})
       } else if length == best {
           res = append(res, interval{l, r})
       }
   }
   // output result
   fmt.Fprintln(writer, best, len(res))
   for _, iv := range res {
       fmt.Fprintln(writer, iv.l+1, iv.r+1)
   }
}
