package main

import (
   "fmt"
   "math"
)

var r, h float64
var pi = math.Pi
const eps = 1e-6

func isGreaterThan(a, b float64) bool {
   return a-b > eps
}

func isZero(a float64) bool {
   return math.Abs(a) < eps
}

// map 3D point on cone to (theta, d) on development plane
func findCoordinatesOnPlane(x, y, z float64) (theta, d float64) {
   angle := math.Atan2(y, x)
   if angle < 0 {
       angle += 2 * pi
   }
   theta = angle * r / math.Hypot(r, h)
   if !isZero(z) {
       // slant distance from apex to (x,y,z)
       d = math.Sqrt(x*x + y*y + (h-z)*(h-z))
   } else {
       // full slant height
       d = math.Hypot(r, h)
   }
   return
}

func euclideanDistance(x1, y1, x2, y2 float64) float64 {
   dx := x1 - x2
   dy := y1 - y2
   return math.Hypot(dx, dy)
}

// shortest path along cone surface between two 3D points
func findConeSurfaceDistance(x1, y1, z1, x2, y2, z2 float64) float64 {
   t1, d1 := findCoordinatesOnPlane(x1, y1, z1)
   t2, d2 := findCoordinatesOnPlane(x2, y2, z2)
   // order by theta
   var aT, bT, aD, bD float64
   if t1 > t2 || (isZero(t1-t2) && d1 > d2) {
       aT, aD = t1, d1
       bT, bD = t2, d2
   } else {
       aT, aD = t2, d2
       bT, bD = t1, d1
   }
   // maximum theta on development = alpha
   alpha := 2 * pi * r / math.Hypot(r, h)
   // wrap if longer than half circumference
   if isGreaterThan(aT-bT, alpha/2) {
       aT = bT + (alpha - (aT - bT))
   }
   // planar euclidean distance
   xA := aD * math.Cos(aT)
   yA := aD * math.Sin(aT)
   xB := bD * math.Cos(bT)
   yB := bD * math.Sin(bT)
   return euclideanDistance(xA, yA, xB, yB)
}

// search for optimal brink point on base circle for single-brink path
func sortaBinarySearch(left, right, xHi, yHi, zHi, xLo, yLo, zLo float64) float64 {
   // precompute values at ends
   x := r * math.Cos(left)
   y := r * math.Sin(left)
   leftVal := euclideanDistance(x, y, xLo, yLo) + findConeSurfaceDistance(xHi, yHi, zHi, x, y, zLo)
   x = r * math.Cos(right)
   y = r * math.Sin(right)
   rightVal := euclideanDistance(x, y, xLo, yLo) + findConeSurfaceDistance(xHi, yHi, zHi, x, y, zLo)
   for !isZero(right-left) {
       mid := (left + right) / 2
       x = r * math.Cos(mid)
       y = r * math.Sin(mid)
       midVal := euclideanDistance(x, y, xLo, yLo) + findConeSurfaceDistance(xHi, yHi, zHi, x, y, zLo)
       if leftVal > rightVal {
           left = mid
           leftVal = midVal
       } else {
           right = mid
           rightVal = midVal
       }
   }
   if leftVal < rightVal {
       return leftVal
   }
   return rightVal
}

// optimal path when one point on base and other on cone
func findOptimalSingleBrinkPointDistance(xHi, yHi, zHi, xLo, yLo, zLo float64) float64 {
   aDeg := pi / 180
   // initial coarse search
   minVal := math.Inf(1)
   minDeg := 0.0
   // sample every 1 degree
   for deg := 0.0; deg < 2*pi; deg += aDeg {
       x := r * math.Cos(deg)
       y := r * math.Sin(deg)
       curr := euclideanDistance(x, y, xLo, yLo) + findConeSurfaceDistance(xHi, yHi, zHi, x, y, zLo)
       if curr < minVal {
           minVal = curr
           minDeg = deg
       }
   }
   // refine around minDeg
   best := sortaBinarySearch(minDeg-aDeg, minDeg+aDeg, xHi, yHi, zHi, xLo, yLo, zLo)
   // also try wider window
   cand := sortaBinarySearch(minDeg-2*aDeg, minDeg+2*aDeg, xHi, yHi, zHi, xLo, yLo, zLo)
   if cand < best {
       best = cand
   }
   return best
}

// double-brink: path touches base twice
func sortaBinarySearch2(left, right, x1, y1, z1, x2, y2, z2 float64) float64 {
   x := r * math.Cos(left)
   y := r * math.Sin(left)
   leftVal := findConeSurfaceDistance(x1, y1, z1, x, y, 0) + findOptimalSingleBrinkPointDistance(x2, y2, z2, x, y, 0)
   x = r * math.Cos(right)
   y = r * math.Sin(right)
   rightVal := findConeSurfaceDistance(x1, y1, z1, x, y, 0) + findOptimalSingleBrinkPointDistance(x2, y2, z2, x, y, 0)
   for !isZero(right-left) {
       mid := (left + right) / 2
       x = r * math.Cos(mid)
       y = r * math.Sin(mid)
       midVal := findConeSurfaceDistance(x1, y1, z1, x, y, 0) + findOptimalSingleBrinkPointDistance(x2, y2, z2, x, y, 0)
       if leftVal > rightVal {
           left = mid
           leftVal = midVal
       } else {
           right = mid
           rightVal = midVal
       }
   }
   if leftVal < rightVal {
       return leftVal
   }
   return rightVal
}

func findOptimalDoubleBrinkPointDistance(x1, y1, z1, x2, y2, z2 float64) float64 {
   aDeg := pi / 180
   baseDeg := math.Atan2(y1, x1)
   best := sortaBinarySearch2(baseDeg-aDeg, baseDeg+aDeg, x1, y1, z1, x2, y2, z2)
   // try wider windows
   for k := 2; k <= 3; k++ {
       cand := sortaBinarySearch2(baseDeg-float64(k)*aDeg, baseDeg+float64(k)*aDeg, x1, y1, z1, x2, y2, z2)
       if cand < best {
           best = cand
       }
   }
   return best
}

func main() {
   var x1, y1, z1, x2, y2, z2 float64
   // read radius, height, and two points
   fmt.Scan(&r, &h)
   fmt.Scan(&x1, &y1, &z1)
   fmt.Scan(&x2, &y2, &z2)
   pt1OnBase := isZero(z1)
   pt2OnBase := isZero(z2)
   var res float64
   if pt1OnBase {
       if pt2OnBase {
           res = euclideanDistance(x1, y1, x2, y2)
       } else {
           res = findOptimalSingleBrinkPointDistance(x2, y2, z2, x1, y1, z1)
           if isZero(x1*x1+y1*y1 - r*r) {
               tmp := findConeSurfaceDistance(x1, y1, z1, x2, y2, z2)
               if tmp < res {
                   res = tmp
               }
           }
       }
   } else if pt2OnBase {
       res = findOptimalSingleBrinkPointDistance(x1, y1, z1, x2, y2, z2)
       if isZero(x2*x2+y2*y2 - r*r) {
           tmp := findConeSurfaceDistance(x1, y1, z1, x2, y2, z2)
           if tmp < res {
               res = tmp
           }
       }
   } else {
       res = findConeSurfaceDistance(x1, y1, z1, x2, y2, z2)
       tmp := findOptimalDoubleBrinkPointDistance(x1, y1, z1, x2, y2, z2)
       if tmp < res {
           res = tmp
       }
   }
   fmt.Printf("%.9f\n", res)
}
