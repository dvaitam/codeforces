package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
)

func min(a, b int64) int64 {
   if a < b {
       return a
   }
   return b
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var C, Hr, Hb, Wr, Wb int64
   if _, err := fmt.Fscan(reader, &C, &Hr, &Hb, &Wr, &Wb); err != nil {
       return
   }
   // Ensure w1 <= w2
   w1, h1 := Wr, Hr
   w2, h2 := Wb, Hb
   if w1 > w2 {
       w1, w2 = w2, w1
       h1, h2 = h2, h1
   }
   // Bound iteration by sqrt(C)
   limit := int64(math.Sqrt(float64(C))) + 1
   var maxJoy int64
   // iterate count of item1
   maxX := min(C/w1, limit)
   for x := int64(0); x <= maxX; x++ {
       rem := C - x*w1
       y := rem / w2
       joy := x*h1 + y*h2
       if joy > maxJoy {
           maxJoy = joy
       }
   }
   // iterate count of item2
   maxY := min(C/w2, limit)
   for y := int64(0); y <= maxY; y++ {
       rem := C - y*w2
       x := rem / w1
       joy := x*h1 + y*h2
       if joy > maxJoy {
           maxJoy = joy
       }
   }
   fmt.Println(maxJoy)
}
