package main

import (
   "fmt"
)

func abs(x int) int {
   if x < 0 {
       return -x
   }
   return x
}

func main() {
   var s, x1, x2, t1, t2, p, d int
   if _, err := fmt.Scan(&s, &x1, &x2); err != nil {
       return
   }
   fmt.Scan(&t1, &t2)
   fmt.Scan(&p, &d)
   // walking time directly
   walkTime := int64(abs(x2-x1)) * int64(t2)
   // if walking is not slower than tram, just walk
   if t2 <= t1 {
       fmt.Println(walkTime)
       return
   }
   // desired tram direction: +1 if going right, -1 if left
   dir := 1
   if x2 < x1 {
       dir = -1
   }
   // simulate tram until it reaches x1 moving in desired direction
   curPos := p
   curDir := d
   var timeToBoard int64
   for {
       // can board if tram is moving towards destination and will reach x1 next
       if curDir == dir {
           if (dir == 1 && x1 >= curPos) || (dir == -1 && x1 <= curPos) {
               timeToBoard += int64(abs(x1-curPos)) * int64(t1)
               break
           }
       }
       // move tram to end and reverse
       if curDir == 1 {
           // go to s
           timeToBoard += int64(s-curPos) * int64(t1)
           curPos = s
           curDir = -1
       } else {
           // go to 0
           timeToBoard += int64(curPos) * int64(t1)
           curPos = 0
           curDir = 1
       }
   }
   // ride tram from x1 to x2
   rideTime := int64(abs(x2-x1)) * int64(t1)
   tramTotal := timeToBoard + rideTime
   // take faster option
   if tramTotal < walkTime {
       fmt.Println(tramTotal)
   } else {
       fmt.Println(walkTime)
   }
}
