package main

import (
   "fmt"
)

func main() {
   var a1, b1, a2, b2 int64
   if _, err := fmt.Scanf("%d %d", &a1, &b1); err != nil {
       return
   }
   if _, err := fmt.Scanf("%d %d", &a2, &b2); err != nil {
       return
   }
   // Copy originals for simulation
   A1, B1 := a1, b1
   A2, B2 := a2, b2
   // Factor counts of 2 and 3 for each bar
   p1a, q1a, x1 := countFactors(a1)
   p1b, q1b, y1 := countFactors(b1)
   p2a, q2a, x2 := countFactors(a2)
   p2b, q2b, y2 := countFactors(b2)
   p1 := p1a + p1b
   q1 := q1a + q1b
   p2 := p2a + p2b
   q2 := q2a + q2b
   // Base (remaining primes other than 2,3)
   base1 := x1 * y1
   base2 := x2 * y2
   if base1 != base2 {
       fmt.Println(-1)
       return
   }
   // Target q is min(q1,q2)
   qTarget := q1
   if q2 < qTarget {
       qTarget = q2
   }
   // Third-chip operations
   d3_1 := q1 - qTarget
   d3_2 := q2 - qTarget
   // New p values after third-chips
   np1 := p1 + d3_1
   np2 := p2 + d3_2
   // Target p is min(np1,np2)
   pTarget := np1
   if np2 < pTarget {
       pTarget = np2
   }
   // Half-split operations
   d2_1 := np1 - pTarget
   d2_2 := np2 - pTarget
   // Total operations
   m := d3_1 + d3_2 + d2_1 + d2_2
   // Simulate operations on dimensions
   // Bar1
   for i := 0; i < d3_1; i++ {
       if A1%3 == 0 {
           A1 = A1/3 * 2
       } else {
           B1 = B1/3 * 2
       }
   }
   for i := 0; i < d2_1; i++ {
       if A1%2 == 0 {
           A1 /= 2
       } else {
           B1 /= 2
       }
   }
   // Bar2
   for i := 0; i < d3_2; i++ {
       if A2%3 == 0 {
           A2 = A2/3 * 2
       } else {
           B2 = B2/3 * 2
       }
   }
   for i := 0; i < d2_2; i++ {
       if A2%2 == 0 {
           A2 /= 2
       } else {
           B2 /= 2
       }
   }
   // Output
   fmt.Println(m)
   fmt.Printf("%d %d\n", A1, B1)
   fmt.Printf("%d %d\n", A2, B2)
}

// countFactors returns counts of 2s and 3s, and the remaining integer
func countFactors(n int64) (cnt2 int, cnt3 int, rem int64) {
   rem = n
   for rem%2 == 0 {
       rem /= 2
       cnt2++
   }
   for rem%3 == 0 {
       rem /= 3
       cnt3++
   }
   return cnt2, cnt3, rem
}
