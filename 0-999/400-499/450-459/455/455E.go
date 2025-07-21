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
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   a := make([]int, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   // prefix sum
   pref := make([]int64, n+1)
   for i := 1; i <= n; i++ {
       pref[i] = pref[i-1] + int64(a[i])
   }
   // prev less element
   prevLess := make([]int, n+1)
   stack := make([]int, 0, n)
   for i := 1; i <= n; i++ {
       for len(stack) > 0 && a[stack[len(stack)-1]] >= a[i] {
           stack = stack[:len(stack)-1]
       }
       if len(stack) > 0 {
           prevLess[i] = stack[len(stack)-1]
       } else {
           prevLess[i] = 0
       }
       stack = append(stack, i)
   }
   var m int
   fmt.Fscan(reader, &m)
   for qi := 0; qi < m; qi++ {
       var x, y int
       fmt.Fscan(reader, &x, &y)
       // compute f(x,y)
       L := y - x + 1
       if L < 1 {
           L = 1
       }
       best := int64(1<<62 - 1)
       pos := y
       for pos >= L {
           mval := a[pos]
           // d = y-pos+1; stays = x-d
           d := y - pos + 1
           stay := x - d
           sumSeg := pref[y] - pref[pos-1]
           cost := sumSeg + int64(stay)*int64(mval)
           if cost < best {
               best = cost
           }
           // jump to next suffix minima
           pos = prevLess[pos]
       }
       fmt.Fprintln(writer, best)
   }
}
