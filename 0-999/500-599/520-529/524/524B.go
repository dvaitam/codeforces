package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   w := make([]int, n)
   h := make([]int, n)
   heights := make([]int, 0, 2*n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &w[i], &h[i])
       heights = append(heights, w[i], h[i])
   }
   sort.Ints(heights)
   // unique candidate heights
   uniq := heights[:0]
   prev := -1
   for _, v := range heights {
       if v != prev {
           uniq = append(uniq, v)
           prev = v
       }
   }
   var minArea int64 = 1<<63 - 1
   for _, H := range uniq {
       totalW := 0
       valid := true
       for i := 0; i < n; i++ {
           wi, hi := w[i], h[i]
           // consider both orientations: (w,h) and (h,w)
           minW := int(^uint(0) >> 1) // max int
           // orientation A: width=wi, height=hi
           if hi <= H {
               if wi < minW {
                   minW = wi
               }
           }
           // orientation B: width=hi, height=wi
           if wi <= H {
               if hi < minW {
                   minW = hi
               }
           }
           if minW == int(^uint(0)>>1) {
               valid = false
               break
           }
           totalW += minW
       }
       if !valid {
           continue
       }
       area := int64(totalW) * int64(H)
       if area < minArea {
           minArea = area
       }
   }
   fmt.Println(minArea)
}
