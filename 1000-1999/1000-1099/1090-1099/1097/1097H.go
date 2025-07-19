package main

import (
   "bufio"
   "fmt"
   "os"
)

var (
   d, m, n int
   gen, b []int
   MASK uint64
   f []Value
   pw []int64
   k int
)

// rotate x right by y in m bits
func rot(x uint64, y int) uint64 {
   return ((x >> uint(y)) | (x << uint(m-y))) & MASK
}

// LBorder represents left border data
type LBorder struct {
   g   []uint64
   len int64
}

func newLBorder() LBorder {
   lb := LBorder{g: make([]uint64, n+1), len: 0}
   for i := range lb.g {
       lb.g[i] = MASK
   }
   return lb
}

func (lb LBorder) at(x int64) uint64 {
   if x < int64(n) {
       return lb.g[x]
   }
   return MASK
}

// rot returns rotated border
func (lb LBorder) rot(y int) LBorder {
   ans := newLBorder()
   ans.len = lb.len
   for i := 0; i < n; i++ {
       ans.g[i] = rot(lb.g[i], y)
   }
   return ans
}

// addL concatenates two LBorder
func addL(lhs, rhs LBorder) LBorder {
   ans := newLBorder()
   ans.len = lhs.len + rhs.len
   for i := 0; i <= n; i++ {
       ans.g[i] = lhs.at(int64(i)) & rhs.at(int64(i)+lhs.len)
   }
   return ans
}

// RBorder represents right border data
type RBorder struct {
   g   []uint64
   len int64
}

func newRBorder() RBorder {
   rb := RBorder{g: make([]uint64, n+1), len: 0}
   for i := range rb.g {
       rb.g[i] = MASK
   }
   return rb
}

func (rb RBorder) at(x int64) uint64 {
   if x > 0 && x < int64(len(rb.g)) {
       return rb.g[x]
   }
   return MASK
}

func (rb RBorder) rot(y int) RBorder {
   ans := newRBorder()
   ans.len = rb.len
   for i := 0; i < n; i++ {
       ans.g[i] = rot(rb.g[i], y)
   }
   return ans
}

func addR(lhs, rhs RBorder) RBorder {
   ans := newRBorder()
   ans.len = lhs.len + rhs.len
   for i := 0; i <= n; i++ {
       ans.g[i] = rhs.at(int64(i)) & lhs.at(int64(i)-rhs.len)
   }
   return ans
}

// Value holds the aggregated data
type Value struct {
   l LBorder
   r RBorder
   a []int64
}

func newValue() Value {
   v := Value{
       l: newLBorder(),
       r: newRBorder(),
       a: make([]int64, m),
   }
   return v
}

// solo initializes Value for a single element
func (v *Value) solo() {
   for i := 0; i < n; i++ {
       mask := (uint64(1) << uint(b[i]+1)) - 1
       v.l.g[i] = mask
       v.r.g[i+1] = mask
   }
   v.l.len = 1
   v.r.len = 1
   if n == 1 {
       for i := 0; i <= b[0]; i++ {
           v.a[i]++
       }
   }
}

func (v Value) length() int64 {
   return v.l.len
}

// rot rotates Value by y
func (v Value) rot(y int) Value {
   res := newValue()
   res.l = v.l.rot(y)
   res.r = v.r.rot(y)
   for i := 0; i < m; i++ {
       res.a[i] = v.a[(i+y)%m]
   }
   return res
}

// addV concatenates two Values
func addV(lhs, rhs Value) Value {
   ans := newValue()
   ans.l = addL(lhs.l, rhs.l)
   ans.r = addR(lhs.r, rhs.r)
   for i := 0; i < m; i++ {
       ans.a[i] = lhs.a[i] + rhs.a[i]
   }
   for i := 1; i < n; i++ {
       if int64(i) <= lhs.length() && int64(i)+rhs.length() >= int64(n) {
           cur := lhs.r.at(int64(i)) & rhs.l.at(int64(i))
           for j := 0; j < m; j++ {
               ans.a[j] += int64((cur >> uint(j)) & 1)
           }
       }
   }
   return ans
}

// precalc builds f and pw tables
func precalc(r int64) {
   f = []Value{newValue()}
   f[0].solo()
   pw = []int64{1}
   k = 0
   for r/int64(pw[k]) >= int64(d) {
       // next level
       f = append(f, newValue())
       pw = append(pw, pw[k]*int64(d))
       // build f[k+1]
       for i := 0; i < d; i++ {
           f[k+1] = addV(f[k+1], f[k].rot(gen[i]))
       }
       k++
   }
}

// eval computes the result for x
func eval(x0 int64) int64 {
   var y int
   cur := newValue()
   x := x0
   for i := k; i >= 0; i-- {
       cn := x / pw[i]
       for j := int64(0); j < cn; j++ {
           cur = addV(cur, f[i].rot((y+gen[int(j)])%m))
       }
       y = (y + gen[int(cn)]) % m
       x %= pw[i]
   }
   return cur.a[0]
}

func main() {
   in := bufio.NewReader(os.Stdin)
   fmt.Fscan(in, &d, &m)
   gen = make([]int, d)
   for i := 0; i < d; i++ {
       fmt.Fscan(in, &gen[i])
   }
   fmt.Fscan(in, &n)
   b = make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &b[i])
   }
   var l, r int64
   fmt.Fscan(in, &l, &r)
   MASK = (uint64(1) << uint(m)) - 1
   precalc(r)
   ar := eval(r)
   al := eval(l + int64(n) - 2)
   fmt.Println(ar - al)
}
