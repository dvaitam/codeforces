package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

// Point represents a point with x (index) and y (value)
type Point struct {
   x, y int
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &a[i])
   }
   // Prepare points sorted by x (index)
   pts := make([]Point, n)
   for i := 0; i < n; i++ {
       pts[i] = Point{x: i, y: a[i]}
   }
   // Sliding window set sorted by y
   var cand []Point
   left := 0
   const INF = int64(8e18)
   d := INF
   for i := 0; i < n; i++ {
       p := pts[i]
       // Remove points too far in x
       for left < i {
           dx := int64(p.x - pts[left].x)
           if dx*dx <= d {
               break
           }
           // remove pts[left] from cand
           // find by y then x
           for j := range cand {
               if cand[j].x == pts[left].x && cand[j].y == pts[left].y {
                   cand = append(cand[:j], cand[j+1:]...)
                   break
               }
           }
           left++
       }
       // Find insertion position by y then x
       pos := sort.Search(len(cand), func(i int) bool {
           if cand[i].y < p.y {
               return false
           }
           if cand[i].y > p.y {
               return true
           }
           return cand[i].x >= p.x
       })
       // Check neighbors in cand within y-range
       for j := pos - 1; j >= 0; j-- {
           dy := int64(p.y - cand[j].y)
           if dy*dy > d {
               break
           }
           dx := int64(p.x - cand[j].x)
           dist := dx*dx + dy*dy
           if dist < d {
               d = dist
           }
       }
       for j := pos; j < len(cand); j++ {
           dy := int64(cand[j].y - p.y)
           if dy*dy > d {
               break
           }
           dx := int64(p.x - cand[j].x)
           dist := dx*dx + dy*dy
           if dist < d {
               d = dist
           }
       }
       // Insert p into cand at pos
       cand = append(cand, Point{})
       copy(cand[pos+1:], cand[pos:])
       cand[pos] = p
   }
   // Output result
   // If no pair found (n<2), d remains INF, but n>=2 by constraints
   fmt.Println(d)
}
