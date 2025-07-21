package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var x, y, z, k int64
   if _, err := fmt.Fscan(in, &x, &y, &z, &k); err != nil {
       return
   }
   caps := []int64{x - 1, y - 1, z - 1}
   sort.Slice(caps, func(i, j int) bool { return caps[i] < caps[j] })
   sumcaps := caps[0] + caps[1] + caps[2]
   // if we can use all possible cuts
   if k >= sumcaps {
       fmt.Fprintln(os.Stdout, x*y*z)
       return
   }
   T := k
   // binary search fill level L
   low := int64(0)
   high := caps[2]
   for low < high {
       mid := (low + high) / 2
       var s int64
       for i := 0; i < 3; i++ {
           if caps[i] < mid {
               s += caps[i]
           } else {
               s += mid
           }
       }
       if s >= T {
           high = mid
       } else {
           low = mid + 1
       }
   }
   L := low
   xcuts := make([]int64, 3)
   var sumX int64
   for i := 0; i < 3; i++ {
       if caps[i] < L {
           xcuts[i] = caps[i]
       } else {
           xcuts[i] = L
       }
       sumX += xcuts[i]
   }
   // remove excess allocations
   E := sumX - T
   for i := 2; i >= 0 && E > 0; i-- {
       if xcuts[i] == L {
           xcuts[i]--
           E--
       }
   }
   // compute result as product of (cuts+1)
   var result int64 = 1
   for i := 0; i < 3; i++ {
       result *= (xcuts[i] + 1)
   }
   fmt.Fprintln(os.Stdout, result)
}
