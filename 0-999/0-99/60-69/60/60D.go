package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
)

// DSU with int32 indices
var parent []int32
var rankArr []uint8

func find(u int32) int32 {
   if parent[u] != u {
       parent[u] = find(parent[u])
   }
   return parent[u]
}

func union(u, v int32) {
   ru := find(u)
   rv := find(v)
   if ru == rv {
       return
   }
   // union by rank
   if rankArr[ru] < rankArr[rv] {
       parent[ru] = rv
   } else if rankArr[ru] > rankArr[rv] {
       parent[rv] = ru
   } else {
       parent[rv] = ru
       rankArr[ru]++
   }
}

func gcd(a, b int) int {
   for b != 0 {
       a, b = b, a%b
   }
   return a
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   // Read values and build map to index
   vals := make([]int, n)
   maxA := 0
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &vals[i])
       if vals[i] > maxA {
           maxA = vals[i]
       }
   }
   // map value to index, -1 if absent
   mv := make([]int32, maxA+1)
   for i := 0; i <= maxA; i++ {
       mv[i] = -1
   }
   for i, v := range vals {
       mv[v] = int32(i)
   }
   // init DSU
   parent = make([]int32, n)
   rankArr = make([]uint8, n)
   for i := 0; i < n; i++ {
       parent[i] = int32(i)
       rankArr[i] = 0
   }
   // Generate primitive Pythagorean triples
   uMax := int(math.Sqrt(float64(maxA)))
   for u := 2; u <= uMax; u++ {
       // ensure u and v have opposite parity
       vStart := 1
       if u&1 == 1 {
           vStart = 2
       }
       for v := vStart; v < u; v += 2 {
           if gcd(u, v) != 1 {
               continue
           }
           uu := u * u
           vv := v * v
           x := uu - vv
           y := 2 * u * v
           z := uu + vv
           // connect x and y
           if x <= maxA && y <= maxA {
               ix := mv[x]
               iy := mv[y]
               if ix >= 0 && iy >= 0 {
                   union(ix, iy)
               }
           }
           // connect x and z
           if x <= maxA && z <= maxA {
               ix := mv[x]
               iz := mv[z]
               if ix >= 0 && iz >= 0 {
                   union(ix, iz)
               }
           }
           // connect y and z
           if y <= maxA && z <= maxA {
               iy := mv[y]
               iz := mv[z]
               if iy >= 0 && iz >= 0 {
                   union(iy, iz)
               }
           }
       }
   }
   // count components
   seen := make([]bool, n)
   var comps int
   for i := 0; i < n; i++ {
       ri := find(int32(i))
       if !seen[ri] {
           seen[ri] = true
           comps++
       }
   }
   // output result
   fmt.Println(comps)
}
