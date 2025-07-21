package main

import "fmt"

func main() {
   var vp, vd, t, f, c int
   if _, err := fmt.Scan(&vp, &vd, &t, &f, &c); err != nil {
       return
   }
   if vd <= vp {
       fmt.Println(0)
       return
   }
   // Initial distance when dragon starts chase
   cur := float64(vp * t)
   vpF := float64(vp)
   vdF := float64(vd)
   fF := float64(f)
   cF := float64(c)
   count := 0
   for {
       // Time for dragon to catch princess
       tCatch := cur / (vdF - vpF)
       // Distance at catch moment
       catchDist := cur + vpF*tCatch
       // If catch happens at or beyond castle, no more bijous needed
       if catchDist >= cF {
           break
       }
       count++
       // Dragon returns to cave and fixes treasury
       downtime := catchDist/vdF + fF
       // Princess continues running during downtime
       cur = catchDist + vpF*downtime
   }
   fmt.Println(count)
}
