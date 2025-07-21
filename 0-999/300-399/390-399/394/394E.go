package main

import (
   "bufio"
   "fmt"
   "os"
   "strconv"
   "strings"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   // read int
   readInt := func() int {
       s, _ := reader.ReadString('\n')
       s = strings.TrimSpace(s)
       v, _ := strconv.Atoi(s)
       return v
   }
   n := readInt()
   var sumX, sumY, sumSq float64
   for i := 0; i < n; i++ {
       line, _ := reader.ReadString('\n')
       fields := strings.Fields(line)
       xi, _ := strconv.ParseFloat(fields[0], 64)
       yi, _ := strconv.ParseFloat(fields[1], 64)
       sumX += xi
       sumY += yi
       sumSq += xi*xi + yi*yi
   }
   Cx := sumX / float64(n)
   Cy := sumY / float64(n)
   m := readInt()
   pts := make([]struct{ x, y float64 }, m)
   for i := 0; i < m; i++ {
       line, _ := reader.ReadString('\n')
       f := strings.Fields(line)
       xi, _ := strconv.ParseFloat(f[0], 64)
       yi, _ := strconv.ParseFloat(f[1], 64)
       pts[i].x = xi
       pts[i].y = yi
   }
   // check if centroid inside polygon (clockwise order)
   inside := true
   for i := 0; i < m; i++ {
       j := (i + 1) % m
       // edge from pts[i] to pts[j]
       ex := pts[j].x - pts[i].x
       ey := pts[j].y - pts[i].y
       wx := Cx - pts[i].x
       wy := Cy - pts[i].y
       // cross = ex*wy - ey*wx
       if ex*wy - ey*wx > 0 {
           inside = false
           break
       }
   }
   var d2 float64
   if inside {
       d2 = 0
   } else {
       // find closest point on polygon boundary
       d2 = 1e300
       for i := 0; i < m; i++ {
           j := (i + 1) % m
           x1, y1 := pts[i].x, pts[i].y
           x2, y2 := pts[j].x, pts[j].y
           vx := x2 - x1
           vy := y2 - y1
           wx := Cx - x1
           wy := Cy - y1
           vlen2 := vx*vx + vy*vy
           var t float64
           if vlen2 > 0 {
               t = (vx*wx + vy*wy) / vlen2
           } else {
               t = 0
           }
           var px, py float64
           if t <= 0 {
               px, py = x1, y1
           } else if t >= 1 {
               px, py = x2, y2
           } else {
               px = x1 + vx*t
               py = y1 + vy*t
           }
           dx := Cx - px
           dy := Cy - py
           dist2 := dx*dx + dy*dy
           if dist2 < d2 {
               d2 = dist2
           }
       }
   }
   // result = sumSq - n*(Cx^2+Cy^2) + n*d2
   res := sumSq - float64(n)*(Cx*Cx+Cy*Cy) + float64(n)*d2
   fmt.Printf("%.10f\n", res)
}
