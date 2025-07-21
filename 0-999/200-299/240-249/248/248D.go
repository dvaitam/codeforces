package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   var t int64
   if _, err := fmt.Fscan(reader, &n, &t); err != nil {
       return
   }
   var s string
   if _, err := fmt.Fscan(reader, &s); err != nil {
       return
   }
   houses := make([]int64, 0, n)
   shops := make([]int64, 0, n)
   for i, ch := range s {
       pos := int64(i + 1)
       switch ch {
       case 'H':
           houses = append(houses, pos)
       case 'S':
           shops = append(shops, pos)
       }
   }
   H := len(houses)
   S := len(shops)
   // If even with all initial sweets, travel time exceeds t, impossible
   if H == 0 {
       fmt.Println(0)
       return
   }
   lastPos := houses[H-1]
   if lastPos > t {
       fmt.Println(-1)
       return
   }
   // helper: check if k initial sweets is enough within time t
   var check = func(k int) bool {
       need := H - k
       if need <= 0 {
           return lastPos <= t
       }
       if need > S {
           return false
       }
       var cost int64
       for i := 0; i < need; i++ {
           // match houses[k+i] with shops[i]
           dh := shops[i] - houses[k+i]
           if dh > 0 {
               cost += dh * 2
               if cost + lastPos > t {
                   return false
               }
           }
       }
       return cost + lastPos <= t
   }
   // binary search minimal k in [0..H]
   lo, hi := 0, H
   res := -1
   for lo <= hi {
       mid := (lo + hi) / 2
       if check(mid) {
           res = mid
           hi = mid - 1
       } else {
           lo = mid + 1
       }
   }
   fmt.Println(res)
}
