package main

import (
   "io"
   "os"
   "sort"
   "strconv"
)

// Point represents a 2D point with integer coordinates.
type Point struct { x, y int64 }

func main() {
   data, _ := io.ReadAll(os.Stdin)
   var idx int
   // readInt parses the next integer from data.
   readInt := func() int64 {
       // skip non-digit, non-minus bytes
       for idx < len(data) && (data[idx] < '0' || data[idx] > '9') && data[idx] != '-' {
           idx++
       }
       neg := false
       if idx < len(data) && data[idx] == '-' {
           neg = true
           idx++
       }
       var x int64
       for idx < len(data) && data[idx] >= '0' && data[idx] <= '9' {
           x = x*10 + int64(data[idx]-'0')
           idx++
       }
       if neg {
           x = -x
       }
       return x
   }

   n := int(readInt())
   P := make([]Point, n)
   for i := 0; i < n; i++ {
       xi := readInt()
       yi := readInt()
       P[i].x = xi
       // transform y coordinate as in original solution
       P[i].y = yi - xi*xi
   }
   // sort by x asc, and by y desc for equal x
   sort.Slice(P, func(i, j int) bool {
       if P[i].x != P[j].x {
           return P[i].x < P[j].x
       }
       return P[i].y > P[j].y
   })
   // build lower convex hull
   C := make([]Point, 0, n)
   for i := 0; i < n; i++ {
       if i > 0 && P[i].x == P[i-1].x {
           continue
       }
       for len(C) > 1 && cross(sub(C[len(C)-1], C[len(C)-2]), sub(P[i], C[len(C)-1])) >= 0 {
           C = C[:len(C)-1]
       }
       C = append(C, P[i])
   }
   res := 0
   if len(C) > 0 {
       // number of segments is hull points minus one
       res = len(C) - 1
   }
   // output result
   os.Stdout.WriteString(strconv.Itoa(res))
}

// sub returns the vector a - b.
func sub(a, b Point) Point {
   return Point{a.x - b.x, a.y - b.y}
}

// cross returns the cross product of vectors a and b.
func cross(a, b Point) int64 {
   return a.x*b.y - a.y*b.x
}
