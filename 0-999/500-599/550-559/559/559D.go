package main

import (
   "bufio"
   "fmt"
   "os"
)

const (
   magic = 49
   D     = 105
)

var half [D]float64
var two  [D]float64

// gcd returns the greatest common divisor of a and b.
func gcd(a, b int) int {
   for b != 0 {
       a, b = b, a%b
   }
   if a < 0 {
       return -a
   }
   return a
}

func min(a, b int) int {
   if a < b {
       return a
   }
   return b
}

// Point represents a 2D integer point.
type Point struct { x, y int }

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   fmt.Fscan(in, &n)
   p := make([]Point, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &p[i].x, &p[i].y)
   }
   // Precompute powers of two and halves
   half[0], two[0] = 1.0, 1.0
   for i := 1; i < D; i++ {
       half[i] = half[i-1] / 2.0
       two[i] = two[i-1] * 2.0
   }
   area, border := 0.0, 0.0
   // Precompute denominator for n <= 50 case
   var z float64
   if n <= 50 {
       z = two[n] - 1.0 - float64(n) - float64(n*(n-1)/2)
   }
   // Main computation
   for i := 0; i < n; i++ {
       for j0 := i + 1; j0 <= min(i+n-2, i+magic); j0++ {
           var now float64
           if n > 50 {
               now = half[j0 - i + 1]
           } else {
               now = two[n-(j0 - i + 1)] - 1.0
               now = now / z
           }
           jj := j0 % n
           dx := p[i].x - p[jj].x
           if dx < 0 {
               dx = -dx
           }
           dy := p[i].y - p[jj].y
           if dy < 0 {
               dy = -dy
           }
           line := gcd(dx, dy)
           border += float64(line) * now
           cross := float64(p[i].x)*float64(p[jj].y) - float64(p[i].y)*float64(p[jj].x)
           area += 0.5 * now * cross
       }
   }
   result := area - border/2.0 + 1.0
   fmt.Printf("%.10f\n", result)
}
