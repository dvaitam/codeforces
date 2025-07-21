package main

import (
   "bufio"
   "fmt"
   "os"
)

// Binary Indexed Tree for int64
type BIT struct {
   n int
   t []int64
}

func newBIT(n int) *BIT {
   return &BIT{n: n, t: make([]int64, n+1)}
}

// add delta at position i
func (b *BIT) update(i int, delta int64) {
   for ; i <= b.n; i += i & -i {
       b.t[i] += delta
   }
}

// sum of [1..i]
func (b *BIT) query(i int) int64 {
   var s int64
   for ; i > 0; i -= i & -i {
       s += b.t[i]
   }
   return s
}

// sum of [l..r]
func (b *BIT) rangeQuery(l, r int) int64 {
   if l > r {
       return 0
   }
   return b.query(r) - b.query(l-1)
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   a := make([]int64, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   // Precompute sequence S[z][1..T]
   const zMin, zMax = 2, 6
   T := make([]int, zMax+1)
   S := make([][]int64, zMax+1)
   for z := zMin; z <= zMax; z++ {
       t := 2 * (z - 1)
       T[z] = t
       S[z] = make([]int64, t+1)
       for j := 1; j <= t; j++ {
           m := j % t
           if m == 0 {
               m = t
           }
           switch {
           case m == 0:
               S[z][j] = 2
           case m <= z:
               S[z][j] = int64(m)
           default:
               S[z][j] = int64(2*z - m)
           }
       }
   }
   // Initialize BITs for each z and remainder class k
   bits := make([][]*BIT, zMax+1)
   for z := zMin; z <= zMax; z++ {
       bits[z] = make([]*BIT, T[z])
       for k := 0; k < T[z]; k++ {
           bits[z][k] = newBIT(n)
       }
   }
   // Build initial
   for i := 1; i <= n; i++ {
       v := a[i]
       for z := zMin; z <= zMax; z++ {
           k := i % T[z]
           bits[z][k].update(i, v)
       }
   }
   var m int
   fmt.Fscan(reader, &m)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   for qi := 0; qi < m; qi++ {
       var ty int
       fmt.Fscan(reader, &ty)
       if ty == 1 {
           var p int
           var v int64
           fmt.Fscan(reader, &p, &v)
           delta := v - a[p]
           if delta != 0 {
               for z := zMin; z <= zMax; z++ {
                   k := p % T[z]
                   bits[z][k].update(p, delta)
               }
               a[p] = v
           }
       } else if ty == 2 {
           var l, r, z int
           fmt.Fscan(reader, &l, &r, &z)
           var res int64
           t := T[z]
           // sum over sequence positions j = 1..len
           for j := 1; j <= t; j++ {
               k := (l-1 + j) % t
               sumk := bits[z][k].rangeQuery(l, r)
               res += S[z][j] * sumk
           }
           fmt.Fprintln(writer, res)
       }
   }
}
