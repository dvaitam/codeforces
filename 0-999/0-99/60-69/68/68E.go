package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
)

type Point struct{ x, y float64 }

func dist2(a, b Point) float64 {
   dx := a.x - b.x
   dy := a.y - b.y
   return dx*dx + dy*dy
}

var eps = 1e-6

// isZero reports if v is approximately zero
func isZero(v float64) bool {
   return math.Abs(v) < eps
}

// isEqual reports if a and b are approximately equal
func isEqual(a, b float64) bool {
   return math.Abs(a-b) < eps
}

// isometry represented by R (2x2) and t
type Iso struct{ a11, a12, a21, a22, tx, ty float64 }

func (iso Iso) apply(p Point) Point {
   return Point{
       iso.a11*p.x + iso.a12*p.y + iso.tx,
       iso.a21*p.x + iso.a22*p.y + iso.ty,
   }
}

// solveIsoFromVec finds R,t such that R*v = w, mapping by rotation/reflection options
func solveIsoFromVec(pA, pB, cA, cB Point) []Iso {
   // v = pB - pA, w = cB - cA
   v := Point{pB.x - pA.x, pB.y - pA.y}
   w := Point{cB.x - cA.x, cB.y - cA.y}
   lv2 := v.x*v.x + v.y*v.y
   if !isEqual(lv2, w.x*w.x+w.y*w.y) {
       return nil
   }
   // build orthonormal basis
   lv := math.Sqrt(lv2)
   if isZero(lv) {
       return nil
   }
   e1 := Point{v.x / lv, v.y / lv}
   e2 := Point{-e1.y, e1.x}
   fw := Point{w.x / lv, w.y / lv}
   f2 := Point{-fw.y, fw.x}
   // rotation R: maps e1->fw, e2->f2
   // R = [fw f2] * [e1 e2]^T
   r11 := fw.x*e1.x + f2.x*e2.x
   r12 := fw.x*e1.y + f2.x*e2.y
   r21 := fw.y*e1.x + f2.y*e2.x
   r22 := fw.y*e1.y + f2.y*e2.y
   t1x := cA.x - (r11*pA.x + r12*pA.y)
   t1y := cA.y - (r21*pA.x + r22*pA.y)
   // reflection R2: maps e1->fw, e2->-f2
   rf11 := fw.x*e1.x - f2.x*e2.x
   rf12 := fw.x*e1.y - f2.x*e2.y
   rf21 := fw.y*e1.x - f2.y*e2.x
   rf22 := fw.y*e1.y - f2.y*e2.y
   t2x := cA.x - (rf11*pA.x + rf12*pA.y)
   t2y := cA.y - (rf21*pA.x + rf22*pA.y)
   return []Iso{
       {r11, r12, r21, r22, t1x, t1y},
       {rf11, rf12, rf21, rf22, t2x, t2y},
   }
}

var triangles [4][3]Point
var best int

func dfs(idx int, classes []Point) {
   if idx == 4 {
       if len(classes) < best {
           best = len(classes)
       }
       return
   }
   // prune
   if len(classes) >= best {
       return
   }
   tri := triangles[idx]
   n := len(classes)
   // for k matched existing: from min(3,n) down to 0
   for k := min(3, n); k >= 0; k-- {
       // generate subsets of vertices size k
       var vertsComb [][]int
       chooseVerts := func(start, need int, cur []int) {
           if need == 0 {
               tmp := make([]int, len(cur)); copy(tmp, cur)
               vertsComb = append(vertsComb, tmp)
               return
           }
           for i := start; i < 3; i++ {
               chooseVerts(i+1, need-1, append(cur, i))
           }
       }
       chooseVerts(0, k, nil)
       for _, verts := range vertsComb {
           // choose k existing classes
           var clsComb [][]int
           if k == 0 {
               clsComb = append(clsComb, []int{})
           } else {
               var cur []int
               var dfsCls func(start, need int)
               dfsCls = func(start, need int) {
                   if need == 0 {
                       tmp := make([]int, len(cur)); copy(tmp, cur)
                       clsComb = append(clsComb, tmp)
                       return
                   }
                   for j := start; j < n; j++ {
                       dfsCls(j+1, need-1)
                   }
               }
               dfsCls(0, k)
           }
           for _, cls := range clsComb {
               // match cls to verts: permutations
               perm := make([]int, k)
               used := make([]bool, k)
               var dfsPerm func(pos int)
               dfsPerm = func(pos int) {
                   if pos == k {
                       // build mapping mapV
                       mapV := []int{-1, -1, -1}
                       for i, v := range verts {
                           mapV[v] = cls[perm[i]]
                       }
                       // count matched
                       m := k
                       if m >= 2 {
                           // try all pairs of matched to define iso
                           vs := verts
                           // iterate over any two indices in verts
                           for i1 := 0; i1 < len(vs); i1++ {
                               for i2 := i1 + 1; i2 < len(vs); i2++ {
                                   v1, v2 := vs[i1], vs[i2]
                                   c1, c2 := mapV[v1], mapV[v2]
                                   for _, iso := range solveIsoFromVec(tri[v1], tri[v2], classes[c1], classes[c2]) {
                                       // check all matched
                                       ok := true
                                       for _, v := range vs {
                                           c := mapV[v]
                                           p := iso.apply(tri[v])
                                           if !isEqual(p.x, classes[c].x) || !isEqual(p.y, classes[c].y) {
                                               ok = false; break
                                           }
                                       }
                                       if !ok {
                                           continue
                                       }
                                       // create new classes if any
                                       newClasses := make([]Point, len(classes))
                                       copy(newClasses, classes)
                                       valid := true
                                       for v := 0; v < 3; v++ {
                                           if mapV[v] < 0 {
                                               p := iso.apply(tri[v])
                                               // check not coincident with existing
                                               for _, q := range newClasses {
                                                   if isEqual(p.x, q.x) && isEqual(p.y, q.y) {
                                                       valid = false; break
                                                   }
                                               }
                                               if !valid {
                                                   break
                                               }
                                               newClasses = append(newClasses, p)
                                           }
                                       }
                                       if valid {
                                           dfs(idx+1, newClasses)
                                       }
                                   }
                               }
                           }
                       } else if m == 1 {
                           // one anchor
                           // R = I, t = c - p
                           v0 := verts[0]
                           c0 := mapV[v0]
                           t := Point{classes[c0].x - tri[v0].x, classes[c0].y - tri[v0].y}
                           iso := Iso{1, 0, 0, 1, t.x, t.y}
                           newClasses := make([]Point, len(classes))
                           copy(newClasses, classes)
                           valid := true
                           for v := 0; v < 3; v++ {
                               if mapV[v] < 0 {
                                   p := iso.apply(tri[v])
                                   // avoid accidental overlap
                                   for _, q := range newClasses {
                                       if isEqual(p.x, q.x) && isEqual(p.y, q.y) {
                                           valid = false; break
                                       }
                                   }
                                   if !valid {
                                       break
                                   }
                                   newClasses = append(newClasses, p)
                               }
                           }
                           if valid {
                               dfs(idx+1, newClasses)
                           }
                       } else {
                           // m == 0, place with offset to avoid collision
                           // compute max existing coords
                           mx, my := 0.0, 0.0
                           for _, q := range classes {
                               if q.x > mx {
                                   mx = q.x
                               }
                               if q.y > my {
                                   my = q.y
                               }
                           }
                           // use translation shifting by +1
                           tX, tY := mx+1, my+1
                           iso := Iso{1, 0, 0, 1, tX, tY}
                           newClasses := make([]Point, len(classes))
                           copy(newClasses, classes)
                           valid := true
                           for v := 0; v < 3; v++ {
                               p := iso.apply(tri[v])
                               for _, q := range newClasses {
                                   if isEqual(p.x, q.x) && isEqual(p.y, q.y) {
                                       valid = false; break
                                   }
                               }
                               if !valid {
                                   break
                               }
                               newClasses = append(newClasses, p)
                           }
                           if valid {
                               dfs(idx+1, newClasses)
                           }
                       }
                       return
                   }
                   for i := 0; i < k; i++ {
                       if !used[i] {
                           used[i] = true
                           perm[pos] = i
                           dfsPerm(pos + 1)
                           used[i] = false
                       }
                   }
               }
               dfsPerm(0)
           }
       }
       // if any solution found with this k, we can break to prefer larger k
       // but we don't track; instead, full search pruning by best
   }
}

func min(a, b int) int {
   if a < b {
       return a
   }
   return b
}

func main() {
   in := bufio.NewReader(os.Stdin)
   for i := 0; i < 4; i++ {
       var x1, y1, x2, y2, x3, y3 float64
       if _, err := fmt.Fscan(in, &x1, &y1, &x2, &y2, &x3, &y3); err != nil {
           return
       }
       triangles[i][0] = Point{x1, y1}
       triangles[i][1] = Point{x2, y2}
       triangles[i][2] = Point{x3, y3}
   }
   // normalize triangles by subtracting their first point
   for i := 0; i < 4; i++ {
       base := triangles[i][0]
       for j := 0; j < 3; j++ {
           triangles[i][j].x -= base.x
           triangles[i][j].y -= base.y
       }
   }
   // for first triangle, place at classes
   // compute distances
   p0, p1, p2 := triangles[0][0], triangles[0][1], triangles[0][2]
   d01 := math.Hypot(p1.x-p0.x, p1.y-p0.y)
   // place p0->(0,0), p1->(d01,0)
   var classes []Point
   classes = append(classes, Point{0, 0})
   classes = append(classes, Point{d01, 0})
   // compute p2
   l02 := math.Hypot(p2.x-p0.x, p2.y-p0.y)
   l12 := math.Hypot(p2.x-p1.x, p2.y-p1.y)
   // x = (d01^2 + l02^2 - l12^2)/(2*d01)
   x := (d01*d01 + l02*l02 - l12*l12) / (2 * d01)
   y2 := math.Sqrt(math.Max(0, l02*l02-x*x))
   classes = append(classes, Point{x, y2})
   best = 12
   dfs(1, classes)
   fmt.Println(best)
}
