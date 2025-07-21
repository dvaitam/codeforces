package main

import (
   "bufio"
   "fmt"
   "io"
   "os"
   "strconv"
)

type Point struct { x, y int64 }

func minIndex(ps []Point) int {
   idx := 0
   for i := 1; i < len(ps); i++ {
       if ps[i].y < ps[idx].y || (ps[i].y == ps[idx].y && ps[i].x < ps[idx].x) {
           idx = i
       }
   }
   return idx
}

func rotateTo(ps []Point, idx int) []Point {
   n := len(ps)
   res := make([]Point, n)
   for i := 0; i < n; i++ {
       res[i] = ps[(idx+i)%n]
   }
   return res
}

func edgeVectors(ps []Point) []Point {
   n := len(ps)
   ev := make([]Point, n)
   for i := 0; i < n; i++ {
       j := (i + 1) % n
       ev[i] = Point{ps[j].x - ps[i].x, ps[j].y - ps[i].y}
   }
   return ev
}

func cross(a, b Point) int64 {
   return a.x*b.y - a.y*b.x
}

// Minkowski sum of two convex polygons CCW
func minkowski(a, b []Point) []Point {
   // rotate to lowest
   ia := minIndex(a)
   ib := minIndex(b)
   A := rotateTo(a, ia)
   B := rotateTo(b, ib)
   ea := edgeVectors(A)
   eb := edgeVectors(B)
   na, nb := len(ea), len(eb)
   i, j := 0, 0
   // starting point
   res := make([]Point, 0, na+nb+1)
   res = append(res, Point{A[0].x + B[0].x, A[0].y + B[0].y})
   for i < na || j < nb {
       var e Point
       if i < na && j < nb {
           if cross(ea[i], eb[j]) > 0 {
               e = ea[i]; i++
           } else {
               e = eb[j]; j++
           }
       } else if i < na {
           e = ea[i]; i++
       } else {
           e = eb[j]; j++
       }
       last := res[len(res)-1]
       res = append(res, Point{last.x + e.x, last.y + e.y})
   }
   // remove last point if equal to first
   if len(res) > 1 {
       fst := res[0]
       lst := res[len(res)-1]
       if fst.x == lst.x && fst.y == lst.y {
           res = res[:len(res)-1]
       }
   }
   return res
}

// point in convex polygon including boundary
func inConvex(poly []Point, q Point) bool {
   n := len(poly)
   if n == 0 {
       return false
   }
   // reference
   p0 := poly[0]
   if q.x == p0.x && q.y == p0.y {
       return true
   }
   // outside wedge
   if cross(Point{poly[1].x - p0.x, poly[1].y - p0.y}, Point{q.x - p0.x, q.y - p0.y}) < 0 {
       return false
   }
   if cross(Point{poly[n-1].x - p0.x, poly[n-1].y - p0.y}, Point{q.x - p0.x, q.y - p0.y}) > 0 {
       return false
   }
   // binary search
   lo, hi := 1, n-1
   for lo+1 < hi {
       mid := (lo + hi) >> 1
       if cross(Point{poly[mid].x - p0.x, poly[mid].y - p0.y}, Point{q.x - p0.x, q.y - p0.y}) >= 0 {
           lo = mid
       } else {
           hi = mid
       }
   }
   // check inside triangle p0, poly[lo], poly[lo+1]
   a := poly[lo]
   b := poly[(lo+1)%n]
   return cross(Point{a.x - q.x, a.y - q.y}, Point{b.x - q.x, b.y - q.y}) >= 0
}

// Fast input reader
func readInts(r *bufio.Reader, n int) ([]Point, error) {
   pts := make([]Point, n)
   for i := 0; i < n; i++ {
       var err error
       pts[i].x, err = readInt(r)
       if err != nil {
           return nil, err
       }
       pts[i].y, err = readInt(r)
       if err != nil {
           return nil, err
       }
   }
   return pts, nil
}

func readInt(r *bufio.Reader) (int64, error) {
   var x int64
   var neg bool
   // skip spaces
   for {
       b, err := r.ReadByte()
       if err != nil {
           return 0, err
       }
       if (b >= '0' && b <= '9') || b == '-' {
           r.UnreadByte()
           break
       }
   }
   b, err := r.ReadByte()
   if err != nil {
       return 0, err
   }
   if b == '-' {
       neg = true
   } else if b >= '0' && b <= '9' {
       x = int64(b - '0')
   } else {
       return 0, io.ErrUnexpectedEOF
   }
   for {
       b, err := r.ReadByte()
       if err != nil {
           break
       }
       if b < '0' || b > '9' {
           r.UnreadByte()
           break
       }
       x = x*10 + int64(b-'0')
   }
   if neg {
       x = -x
   }
   return x, nil
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   // read three polygons
   polys := make([][]Point, 3)
   for i := 0; i < 3; i++ {
       ni64, err := readInt(reader)
       if err != nil {
           return
       }
       ni := int(ni64)
       pts, err := readInts(reader, ni)
       if err != nil {
           return
       }
       polys[i] = pts
   }
   // Minkowski sum of three
   sum12 := minkowski(polys[0], polys[1])
   sum123 := minkowski(sum12, polys[2])
   // read m hills
   m64, err := readInt(reader)
   if err != nil {
       return
   }
   m := int(m64)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   for i := 0; i < m; i++ {
       x, err := readInt(reader)
       if err != nil {
           return
       }
       y, err := readInt(reader)
       if err != nil {
           return
       }
       // scale
       q := Point{3 * x, 3 * y}
       if inConvex(sum123, q) {
           out.WriteString("YES\n")
       } else {
           out.WriteString("NO\n")
       }
   }
}
