package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var k, l, r, t, x, y int64
   if _, err := fmt.Fscan(in, &k, &l, &r, &t, &x, &y); err != nil {
       return
   }
   // Case when adding does not increase net water
   if y <= x {
       // best to never add, check if k - t*x >= l
       if k - t*x >= l {
           fmt.Println("Yes")
       } else {
           fmt.Println("No")
       }
       return
   }
   // y > x: possible cycles
   cur := k
   tLeft := t
   seen := make(map[int64]struct{})
   for tLeft > 0 {
       rem := cur % x
       if _, ok := seen[rem]; ok {
           // cycle detected, can continue forever
           fmt.Println("Yes")
           return
       }
       seen[rem] = struct{}{}
       // determine minimal days to skip to allow addition
       var minDays int64
       if cur + y <= r {
           minDays = 0
       } else {
           // need cur - d*x + y <= r  -> d >= (cur+y-r + x-1)/x
           need := cur + y - r
           minDays = (need + x - 1) / x
       }
       // max days we can skip without dropping below l
       maxDays := (cur - l) / x
       if minDays > tLeft {
           // can skip all remaining days safely
           fmt.Println("Yes")
           return
       }
       if maxDays < minDays {
           // cannot skip enough days to add safely
           fmt.Println("No")
           return
       }
       // skip minDays days of no addition
       cur -= minDays * x
       tLeft -= minDays
       // now tLeft > 0 and cur + y <= r
       // perform one add and one usage
       cur += y
       // after addition still must be <= r
       // subtract usage
       cur -= x
       tLeft--
       if cur < l {
           fmt.Println("No")
           return
       }
   }
   // survived all days
   fmt.Println("Yes")
}
