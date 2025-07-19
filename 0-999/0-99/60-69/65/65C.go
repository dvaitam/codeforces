package main

import (
   "fmt"
   "math"
)

// P represents a point or vector in 3D space
type P struct {
   x, y, z float64
}

// add returns the vector sum of p and q
func (p P) add(q P) P {
   return P{p.x + q.x, p.y + q.y, p.z + q.z}
}

// sub returns the vector difference p - q
func (p P) sub(q P) P {
   return P{p.x - q.x, p.y - q.y, p.z - q.z}
}

// mul scales p by factor s
func (p P) mul(s float64) P {
   return P{p.x * s, p.y * s, p.z * s}
}

// norm returns the Euclidean norm of p
func (p P) norm() float64 {
   return math.Sqrt(p.x*p.x + p.y*p.y + p.z*p.z)
}

func main() {
   var n int
   if _, err := fmt.Scan(&n); err != nil {
       return
   }
   v := make([]P, n+1)
   for i := 0; i <= n; i++ {
       fmt.Scan(&v[i].x, &v[i].y, &v[i].z)
   }
   var vh, vs float64
   fmt.Scan(&vh, &vs)
   var h P
   fmt.Scan(&h.x, &h.y, &h.z)

   eps := 1e-11
   // can helicopter reach point a by time t?
   can := func(a P, t float64) bool {
       return (a.sub(h).norm()/vh) < t+eps
   }

   var onde P
   onde.x = -1e10
   var endIdx int
   var tAccum float64
   // traverse each segment
   for i := 0; i < n; i++ {
       segDist := v[i].sub(v[i+1]).norm()
       dt := segDist / vs
       if can(v[i+1], tAccum+dt) {
           lo, hi := 0.0, 1.0
           for it := 0; it < 100; it++ {
               mid := (lo + hi) / 2.0
               m := v[i].add(v[i+1].sub(v[i]).mul(mid))
               tPerson := tAccum + v[i].sub(m).norm()/vs
               if can(m, tPerson) {
                   hi = mid
               } else {
                   lo = mid
               }
           }
           onde = v[i].add(v[i+1].sub(v[i]).mul(lo))
           endIdx = i
           break
       }
       tAccum += dt
   }
   if onde.x <= -1e8 {
       fmt.Println("NO")
       return
   }
   tPerson := tAccum + v[endIdx].sub(onde).norm()/vs
   tHeli := onde.sub(h).norm() / vh
   tFinal := math.Max(tPerson, tHeli)
   fmt.Println("YES")
   fmt.Printf("%.10f\n", tFinal)
   fmt.Printf("%.10f %.10f %.10f\n", onde.x, onde.y, onde.z)
}
