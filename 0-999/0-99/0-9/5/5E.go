package main

import (
   "bufio"
   "fmt"
   "os"
)

type pair struct {
   h   int
   cnt int64
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   h := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &h[i])
   }
   // find index of a maximum height
   maxIdx := 0
   for i := 1; i < n; i++ {
       if h[i] > h[maxIdx] {
           maxIdx = i
       }
   }
   // monotonic stack
   stack := make([]pair, 0, n)
   // push first (max height)
   stack = append(stack, pair{h[maxIdx], 1})
   var res int64
   // traverse the circle from next of maxIdx
   for i := 1; i < n; i++ {
       idx := (maxIdx + i) % n
       cur := h[idx]
       cnt := int64(1)
       // pop smaller
       for len(stack) > 0 && stack[len(stack)-1].h < cur {
           res += stack[len(stack)-1].cnt
           stack = stack[:len(stack)-1]
       }
       if len(stack) > 0 && stack[len(stack)-1].h == cur {
           top := stack[len(stack)-1]
           // pairs among equals
           res += top.cnt
           cnt = top.cnt + 1
           // pop the equal
           stack = stack[:len(stack)-1]
           // one more with next higher
           if len(stack) > 0 {
               res += 1
           }
       } else if len(stack) > 0 {
           // current sees one taller
           res += 1
       }
       // push current
       stack = append(stack, pair{cur, cnt})
   }
   // account for wrap-around visibility when maximum height appears multiple times
   // last hill (before returning to start) may see start max via the other arc
   lastIdx := (maxIdx + n - 1) % n
   if h[lastIdx] < h[maxIdx] && stack[0].cnt > 1 {
       res++
   }
   fmt.Println(res)
}
