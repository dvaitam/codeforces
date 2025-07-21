package main

import (
   "bufio"
   "fmt"
   "os"
)

// Point represents a 2D point with integer coordinates.
type Point struct {
   x, y int
}

// isRight checks if the triangle formed by p1, p2, p3 is right-angled and nondegenerate.
func isRight(p1, p2, p3 Point) bool {
   // squared lengths of sides
   dx := p2.x - p1.x
   dy := p2.y - p1.y
   a := dx*dx + dy*dy
   dx = p3.x - p2.x
   dy = p3.y - p2.y
   b := dx*dx + dy*dy
   dx = p1.x - p3.x
   dy = p1.y - p3.y
   c := dx*dx + dy*dy
   // check nondegenerate: no zero-length side
   if a == 0 || b == 0 || c == 0 {
       return false
   }
   // check Pythagorean theorem
   if a+b == c || a+c == b || b+c == a {
       return true
   }
   return false
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var p [3]Point
   for i := 0; i < 3; i++ {
       fmt.Fscan(reader, &p[i].x, &p[i].y)
   }
   if isRight(p[0], p[1], p[2]) {
       fmt.Println("RIGHT")
       return
   }
   // try moving each point by distance 1 in four directions
   dirs := []Point{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
   for i := 0; i < 3; i++ {
       orig := p[i]
       for _, d := range dirs {
           p[i].x = orig.x + d.x
           p[i].y = orig.y + d.y
           if isRight(p[0], p[1], p[2]) {
               fmt.Println("ALMOST")
               return
           }
       }
       p[i] = orig
   }
   fmt.Println("NEITHER")
}
