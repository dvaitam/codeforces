package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

// circle represents a target circle with its projection interval on Ox axis.
type circle struct {
   left, right int
   xc, r       int
   idx         int
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n int
   fmt.Fscan(in, &n)
   circles := make([]circle, n)
   for i := 0; i < n; i++ {
       var x, r int
       fmt.Fscan(in, &x, &r)
       circles[i] = circle{
           left:  x - r,
           right: x + r,
           xc:    x,
           r:     r,
           idx:   i,
       }
   }
   // sort by left endpoint
   sort.Slice(circles, func(i, j int) bool {
       return circles[i].left < circles[j].left
   })

   var m int
   fmt.Fscan(in, &m)
   ans := make([]int, n)
   for i := range ans {
       ans[i] = -1
   }
   hitCount := 0
   // process shots
   for shot := 1; shot <= m; shot++ {
       var x, y int
       fmt.Fscan(in, &x, &y)
       // find first circle with left >= x
       pos := sort.Search(n, func(i int) bool {
           return circles[i].left >= x
       })
       // if a circle starts at x, check it
       if pos < n && circles[pos].left == x {
           c := circles[pos]
           dx := int64(x - c.xc)
           dy := int64(y)
           if dx*dx+dy*dy <= int64(c.r)*int64(c.r) && ans[c.idx] == -1 {
               ans[c.idx] = shot
               hitCount++
           }
           // also check previous circle in case of touching
           if pos > 0 {
               c2 := circles[pos-1]
               if x <= c2.right {
                   dx2 := int64(x - c2.xc)
                   dy2 := int64(y)
                   if dx2*dx2+dy2*dy2 <= int64(c2.r)*int64(c2.r) && ans[c2.idx] == -1 {
                       ans[c2.idx] = shot
                       hitCount++
                   }
               }
           }
       } else {
           // check the previous circle whose left < x
           if pos > 0 {
               c := circles[pos-1]
               if x <= c.right {
                   dx := int64(x - c.xc)
                   dy := int64(y)
                   if dx*dx+dy*dy <= int64(c.r)*int64(c.r) && ans[c.idx] == -1 {
                       ans[c.idx] = shot
                       hitCount++
                   }
               }
           }
       }
   }
   // output results
   fmt.Fprintln(out, hitCount)
   for i, v := range ans {
       if i > 0 {
           out.WriteByte(' ')
       }
       fmt.Fprint(out, v)
   }
   fmt.Fprintln(out)
}
