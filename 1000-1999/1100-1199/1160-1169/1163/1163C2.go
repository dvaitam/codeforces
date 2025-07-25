package main

import (
   "bufio"
   "fmt"
   "os"
)

func gcd(a, b int64) int64 {
   if a < 0 {
       a = -a
   }
   if b < 0 {
       b = -b
   }
   for b != 0 {
       a, b = b, a%b
   }
   return a
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   fmt.Fscan(in, &n)
   xs := make([]int64, n)
   ys := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &xs[i], &ys[i])
   }
   type Line struct{ A, B, C int64 }
   type Slope struct{ A, B int64 }
   lines := make(map[Line]struct{})
   slopes := make(map[Slope]int64)
   for i := 0; i < n; i++ {
       x1, y1 := xs[i], ys[i]
       for j := i + 1; j < n; j++ {
           x2, y2 := xs[j], ys[j]
           // line through (x1,y1) and (x2,y2): A x + B y + C = 0
           A := y2 - y1
           B := x1 - x2
           C := -(A*x1 + B*y1)
           // normalize
           g := gcd(A, B)
           g = gcd(g, C)
           if g != 0 {
               A /= g; B /= g; C /= g
           }
           // ensure unique sign: A>0 or A==0 && B>0
           if A < 0 || (A == 0 && B < 0) {
               A = -A; B = -B; C = -C
           }
           l := Line{A, B, C}
           if _, exists := lines[l]; !exists {
               lines[l] = struct{}{}
               // record slope
               s := Slope{A, B}
               slopes[s]++
           }
       }
   }
   // total distinct lines
   m := int64(len(lines))
   // total pairs
   total := m * (m - 1) / 2
   // subtract parallel pairs
   var skip int64
   for _, k := range slopes {
       skip += k * (k - 1) / 2
   }
   fmt.Println(total - skip)
}
// End of file
