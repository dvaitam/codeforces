package main

import (
   "bufio"
   "fmt"
   "os"
)

type Segment struct {
   vertical   bool   // true if vertical segment
   fixed      int64  // x if vertical, y if horizontal
   low, high  int64  // inclusive range on moving axis
   headCoord  int64  // coordinate of head on moving axis
   headDir    byte   // 'U','D','L','R'
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n int
   var b int64
   fmt.Fscan(in, &n, &b)
   vert := make([]Segment, 0, n)
   hor := make([]Segment, 0, n)
   for i := 0; i < n; i++ {
       var x0, y0, x1, y1 int64
       fmt.Fscan(in, &x0, &y0, &x1, &y1)
       if x0 == x1 {
           // vertical
           s := Segment{vertical: true, fixed: x0}
           if y0 < y1 {
               s.low, s.high = y0, y1
               s.headCoord = y1
               s.headDir = 'U'
           } else {
               s.low, s.high = y1, y0
               s.headCoord = y1 // head at y1
               s.headDir = 'D'
           }
           vert = append(vert, s)
       } else {
           // horizontal
           s := Segment{vertical: false, fixed: y0}
           if x0 < x1 {
               s.low, s.high = x0, x1
               s.headCoord = x1
               s.headDir = 'R'
           } else {
               s.low, s.high = x1, x0
               s.headCoord = x1
               s.headDir = 'L'
           }
           hor = append(hor, s)
       }
   }
   var q int
   fmt.Fscan(in, &q)
   for i := 0; i < q; i++ {
       var x, y, t int64
       var ws string
       fmt.Fscan(in, &x, &y, &ws, &t)
       dir := ws[0]
       for t > 0 {
           var dx, dy int64
           if dir == 'R' {
               dx = 1
           } else if dir == 'L' {
               dx = -1
           } else if dir == 'U' {
               dy = 1
           } else if dir == 'D' {
               dy = -1
           }
           // find nearest event
           distBound := int64(0)
           if dx > 0 {
               distBound = b - x
           } else if dx < 0 {
               distBound = x
           } else if dy > 0 {
               distBound = b - y
           } else if dy < 0 {
               distBound = y
           }
           bestDist := distBound + 1
           var seg *Segment
           if dx != 0 {
               // horizontal, look for vertical segments
               for j := range vert {
                   s := &vert[j]
                   if y < s.low || y > s.high {
                       continue
                   }
                   d := (s.fixed - x) * dx
                   if d > 0 && d < bestDist {
                       bestDist = d
                       seg = s
                   }
               }
           } else {
               // vertical, look for horizontal segments
               for j := range hor {
                   s := &hor[j]
                   if x < s.low || x > s.high {
                       continue
                   }
                   d := (s.fixed - y) * dy
                   if d > 0 && d < bestDist {
                       bestDist = d
                       seg = s
                   }
               }
           }
           // check if hit segment or boundary
           if seg == nil || bestDist > distBound {
               // boundary hit or no seg before boundary
               move := t
               if move > distBound {
                   move = distBound
               }
               x += dx * move
               y += dy * move
               break
           }
           // hit seg at bestDist
           if t < bestDist {
               // cannot reach seg
               x += dx * t
               y += dy * t
               break
           }
           // move to intersection
           x += dx * bestDist
           y += dy * bestDist
           t -= bestDist
           // traverse along the arrow segment to head
           if seg.vertical {
               // move in y towards headCoord
               dseg := seg.headCoord - y
               var sign int64 = 1
               if dseg < 0 {
                   sign = -1
               }
               lenSeg := dseg * sign
               if t < lenSeg {
                   y += sign * t
                   break
               }
               // reach head
               y = seg.headCoord
               t -= lenSeg
           } else {
               // horizontal seg, move in x
               dseg := seg.headCoord - x
               var sign int64 = 1
               if dseg < 0 {
                   sign = -1
               }
               lenSeg := dseg * sign
               if t < lenSeg {
                   x += sign * t
                   break
               }
               x = seg.headCoord
               t -= lenSeg
           }
           // update direction to seg.headDir
           dir = seg.headDir
       }
       fmt.Fprintf(out, "%d %d\n", x, y)
   }
}
