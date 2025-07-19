package main

import (
   "bufio"
   "fmt"
   "os"
)

var rdr = bufio.NewReader(os.Stdin)

// qt represents a doubled Hurwitz quaternion (s + xi + yj + zk), coords are doubled
type qt struct { s, x, y, z int64 }

// constructors
func newQt(v int64) qt { return qt{2 * v, 0, 0, 0} }
func newQtVals(s, x, y, z int64) qt { return qt{2 * s, 2 * x, 2 * y, 2 * z} }
func newQtDbl(s, x, y, z int64) qt { return qt{s, x, y, z} }

// conjugate
func (q qt) conj() qt { return newQtDbl(q.s, -q.x, -q.y, -q.z) }

// norm returns (s^2 + x^2 + y^2 + z^2)/4
func norm(q qt) int64 { return (q.s*q.s + q.x*q.x + q.y*q.y + q.z*q.z) >> 2 }

// addition, subtraction
func (a qt) add(b qt) qt { return newQtDbl(a.s+b.s, a.x+b.x, a.y+b.y, a.z+b.z) }
func (a qt) sub(b qt) qt { return newQtDbl(a.s-b.s, a.x-b.x, a.y-b.y, a.z-b.z) }

// scalar multiplication
func (q qt) mulScalar(c int64) qt { return newQtDbl(q.s*c, q.x*c, q.y*c, q.z*c) }

// quaternion multiplication
func (a qt) mul(b qt) qt {
   return newQtDbl(
       (a.s*b.s - a.x*b.x - a.y*b.y - a.z*b.z) >> 1,
       (a.s*b.x + a.x*b.s + a.y*b.z - a.z*b.y) >> 1,
       (a.s*b.y + a.y*b.s + a.z*b.x - a.x*b.z) >> 1,
       (a.s*b.z + a.z*b.s + a.x*b.y - a.y*b.x) >> 1,
   )
}

// floorDiv divides u by v, rounding toward -inf
func floorDiv(u, v int64) int64 {
   d := u / v
   r := u % v
   if (u^v) < 0 && r != 0 {
       d--
   }
   return d
}

// rightDiv computes a = b*q + r, returns q (quotient) and r (remainder)
func rightDiv(a, b qt) (qt, qt) {
   numer := b.conj().mul(a)
   den := norm(b)
   s := floorDiv(numer.s, den)
   x := floorDiv(numer.x, den)
   y := floorDiv(numer.y, den)
   z := floorDiv(numer.z, den)
   q1 := newQtDbl(s|1, x|1, y|1, z|1)
   r1 := a.sub(b.mul(q1))
   q2 := newQtDbl((s+1)&^1, (x+1)&^1, (y+1)&^1, (z+1)&^1)
   r2 := a.sub(b.mul(q2))
   if norm(r1) < norm(r2) {
       return q1, r1
   }
   return q2, r2
}

// rightGcd computes right gcd of quaternions a and b
func rightGcd(a, b qt) qt {
   for a.s != 0 || a.x != 0 || a.y != 0 || a.z != 0 {
       _, r := rightDiv(b, a)
       b, a = a, r
   }
   return b
}

// gcd of two ints
func gcd(a, b int64) int64 {
   if a < 0 {
       a = -a
   }
   if b < 0 {
       b = -b
   }
   for b != 0 {
       a, b = b, a%b
   }
   return a
}

func main() {
   var N int
   fmt.Fscan(rdr, &N)
   A := make([]qt, N)
   var k int64
   for i := 0; i < N; i++ {
       var x, y, z int64
       fmt.Fscan(rdr, &x, &y, &z)
       A[i] = newQtVals(0, x, y, z)
       k = gcd(k, x)
       k = gcd(k, y)
       k = gcd(k, z)
   }
   for i := range A {
       A[i].x /= k
       A[i].y /= k
       A[i].z /= k
   }
   var G qt
   for _, a := range A {
       G = rightGcd(G, a)
   }
   normG := norm(G)
   var r2 int64
   for _, a0 := range A {
       a := a0
       normA := norm(a)
       a = G.conj().mul(a)
       a.s /= normG; a.x /= normG; a.y /= normG; a.z /= normG
       g1 := gcd(normA, normG)
       normA /= g1
       curR2 := g1
       a = a.mul(G)
       // imaginary parts are a.x/2, a.y/2, a.z/2
       v1 := a.x >> 1; if v1<0 { v1 = -v1 }
       v2 := a.y >> 1; if v2<0 { v2 = -v2 }
       v3 := a.z >> 1; if v3<0 { v3 = -v3 }
       g2 := gcd(normA, gcd(gcd(v1, v2), v3))
       curR2 *= g2
       r2 = gcd(r2, curR2)
   }
   // compute r = max r such that r^2 divides r2
   var r int64 = 1
   n := r2
   for p := int64(2); p*p <= n; p++ {
       if n%p == 0 {
           cnt := int64(0)
           for n%p == 0 {
               cnt++
               n /= p
           }
           for i := int64(0); i < cnt/2; i++ {
               r *= p
           }
       }
   }
   // quaternion for r
   Q := rightGcd(G, newQt(r))
   area := k * r * k * r
   fmt.Println(area)
   // basis images
   e1 := newQtVals(0, 1, 0, 0)
   e2 := newQtVals(0, 0, 1, 0)
   e3 := newQtVals(0, 0, 0, 1)
   for _, e := range []qt{e1, e2, e3} {
       t := Q.mul(e).mul(Q.conj()).mulScalar(k)
       fmt.Println(t.x/2, t.y/2, t.z/2)
   }
}
