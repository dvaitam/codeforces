package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var t1, t2, x1, x2, t0 int64
   if _, err := fmt.Fscan(reader, &t1, &t2, &x1, &x2, &t0); err != nil {
       return
   }
   // If both water temperatures equal target, open both taps fully
   if t0 == t1 && t0 == t2 {
       fmt.Printf("%d %d\n", x1, x2)
       return
   }
   // If only hot water meets target
   if t0 == t2 {
       fmt.Printf("0 %d\n", x2)
       return
   }
   // Initialize best with using only hot tap
   bestY1, bestY2 := int64(0), x2
   bestN := (t2 - t0) * x2
   bestD := x2
   denomDiff := t2 - t0
   // Try mixing for each possible cold flow
   for y1 := int64(1); y1 <= x1; y1++ {
       // minimal hot flow to reach at least t0
       need := y1 * (t0 - t1)
       var y2 int64
       if need > 0 {
           // denomDiff > 0 here since t0 != t2
           y2 = (need + denomDiff - 1) / denomDiff
       } else {
           y2 = 0
       }
       if y2 > x2 {
           continue
       }
       // Compute deviation numerator and total flow
       sum := y1 + y2
       // deviation N = (t1 - t0)*y1 + (t2 - t0)*y2
       N := (t1 - t0)*y1 + (t2 - t0)*y2
       if N < 0 {
           continue
       }
       D := sum
       // Compare N/D with bestN/bestD: want smaller, ties by larger sum
       if N*bestD < bestN*D || (N*bestD == bestN*D && sum > bestY1+bestY2) {
           bestY1, bestY2 = y1, y2
           bestN, bestD = N, D
       }
   }
   fmt.Printf("%d %d\n", bestY1, bestY2)
}
