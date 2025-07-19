package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
)

const eps = 1e-7

// Point represents a 2D point or vector
type Point struct {
   x, y float64
}

// Sub returns the vector subtraction p - q
func (p Point) Sub(q Point) Point {
   return Point{p.x - q.x, p.y - q.y}
}

// Dot returns the dot product of p and q
func (p Point) Dot(q Point) float64 {
   return p.x*q.x + p.y*q.y
}

// Cross returns the cross product (scalar) of p and q
func (p Point) Cross(q Point) float64 {
   return p.x*q.y - p.y*q.x
}

// eq checks if value is approximately zero
func eq(v float64) bool {
   return math.Abs(v) < eps
}

// escua checks the square condition for points a,b,c,d
func escua(a, b, c, d Point) bool {
   return eq((b.Sub(a)).Dot(b.Sub(c))) &&
       eq((b.Sub(a)).Dot(a.Sub(d))) &&
       eq((c.Sub(a)).Dot(d.Sub(b)))
}

// esrec checks the rectangle condition for points a,b,c,d
func esrec(a, b, c, d Point) bool {
   return eq((b.Sub(a)).Dot(b.Sub(c))) &&
       eq((b.Sub(a)).Dot(a.Sub(d))) &&
       eq((d.Sub(c)).Dot(c.Sub(b)))
}

// isSquare returns true if a,b,c,d form a square in any order
func isSquare(a, b, c, d Point) bool {
   // all permutations of b,c,d
   pts := []Point{b, c, d}
   idxs := [][3]int{{0, 1, 2}, {0, 2, 1}, {1, 0, 2}, {1, 2, 0}, {2, 0, 1}, {2, 1, 0}}
   for _, id := range idxs {
       if escua(a, pts[id[0]], pts[id[1]], pts[id[2]]) {
           return true
       }
   }
   return false
}

// isRect returns true if a,b,c,d form a rectangle in any order
func isRect(a, b, c, d Point) bool {
   pts := []Point{b, c, d}
   idxs := [][3]int{{0, 1, 2}, {0, 2, 1}, {1, 0, 2}, {1, 2, 0}, {2, 0, 1}, {2, 1, 0}}
   for _, id := range idxs {
       if esrec(a, pts[id[0]], pts[id[1]], pts[id[2]]) {
           return true
       }
   }
   return false
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var P [8]Point
   for i := 0; i < 8; i++ {
       if _, err := fmt.Fscan(reader, &P[i].x, &P[i].y); err != nil {
           return
       }
   }
   // try all combinations of 4 for square
   for i := 0; i < 5; i++ {
       for j := i + 1; j < 6; j++ {
           for k := j + 1; k < 7; k++ {
               for l := k + 1; l < 8; l++ {
                   if isSquare(P[i], P[j], P[k], P[l]) {
                       // remaining indices
                       var rem []int
                       for m := 0; m < 8; m++ {
                           if m != i && m != j && m != k && m != l {
                               rem = append(rem, m)
                           }
                       }
                       if len(rem) == 4 && isRect(P[rem[0]], P[rem[1]], P[rem[2]], P[rem[3]]) {
                           fmt.Println("YES")
                           fmt.Printf("%d %d %d %d\n", i+1, j+1, k+1, l+1)
                           fmt.Printf("%d %d %d %d\n", rem[0]+1, rem[1]+1, rem[2]+1, rem[3]+1)
                           return
                       }
                   }
               }
           }
       }
   }
   fmt.Println("NO")
}
