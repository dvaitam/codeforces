package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   fmt.Fscan(reader, &n)
   m := 1 << n
   al := uint32(2*m - 1)

   // popcount for indices
   bit := make([]int, m)
   for i := 1; i < m; i++ {
       bit[i] = bit[i>>1] + (i & 1)
   }

   // precompute tables
   var to [256][256]uint32
   var up [256][256]uint32
   for i := 0; i < 256; i++ {
       for j := 0; j < 256; j++ {
           var t0, t1 uint32
           for k := 0; k < 16; k++ {
               var sum int
               for t := 0; t <= k; t++ {
                   sum += ((i >> t) & 1) * ((j >> (k - t)) & 1)
               }
               sum &= 3
               if sum&1 != 0 {
                   t0 |= 1 << k
               }
               if sum&2 != 0 {
                   t1 |= 1 << k
               }
           }
           to[i][j] = t0
           up[i][j] = t1
       }
   }

   // read input strings
   s1 := make([]byte, m)
   s2 := make([]byte, m)
   fmt.Fscan(reader, &s1)
   fmt.Fscan(reader, &s2)

   // prepare f, g
   f := make([][2]uint32, m)
   g := make([][2]uint32, m)
   for i := 0; i < m; i++ {
       c := s1[i] - '0'
       if c&1 != 0 {
           f[i][0] |= 1 << bit[i]
       }
       if c&2 != 0 {
           f[i][1] |= 1 << bit[i]
       }
       c2 := s2[i] - '0'
       if c2&1 != 0 {
           g[i][0] |= 1 << bit[i]
       }
       if c2&2 != 0 {
           g[i][1] |= 1 << bit[i]
       }
   }

   // define inc, dec
   inc := func(u *[2]uint32, v [2]uint32) {
       u[1] ^= v[1] ^ (u[0] & v[0])
       u[0] ^= v[0]
   }
   dec := func(u *[2]uint32, v [2]uint32) {
       cc0 := (^v[0]) & al
       cc1 := (^v[1]) & al
       inc(u, [2]uint32{cc0, cc1})
       inc(u, [2]uint32{al, 0})
   }

   // FWT
   fwt := func(a [][2]uint32) {
       for l := 2; l <= m; l <<= 1 {
           half := l >> 1
           for j := 0; j < m; j += l {
               for k := 0; k < half; k++ {
                   inc(&a[j+k+half], a[j+k])
               }
           }
       }
   }
   // inverse FWT
   ifwt := func(a [][2]uint32) {
       for l := 2; l <= m; l <<= 1 {
           half := l >> 1
           for j := 0; j < m; j += l {
               for k := 0; k < half; k++ {
                   dec(&a[j+k+half], a[j+k])
               }
           }
       }
   }

   // multiply in transformed domain
   mul := func(u *[2]uint32, b, c [2]uint32) {
       u[0], u[1] = 0, 0
       for i := 0; i < 3; i++ {
           for j := 0; i+j < 3; j++ {
               shiftB := uint(i * 8)
               x0 := (b[0] >> shiftB) & 0xff
               x1 := (b[1] >> shiftB) & 0xff
               shiftC := uint(j * 8)
               y0 := (c[0] >> shiftC) & 0xff
               y1 := (c[1] >> shiftC) & 0xff
               z0 := to[x0][y0]
               z1 := to[x1][y0] ^ to[x0][y1] ^ up[x0][y0]
               shift := uint((i + j) * 8)
               v0 := z0 << shift
               v1 := z1 << shift
               inc(u, [2]uint32{v0, v1})
           }
       }
       u[0] &= al
       u[1] &= al
   }

   // transform f and g
   fwt(f)
   fwt(g)
   // pointwise multiply
   h := make([][2]uint32, m)
   for i := 0; i < m; i++ {
       mul(&h[i], f[i], g[i])
   }
   // inverse transform
   ifwt(h)

   // output result
   res := make([]byte, m)
   for i := 0; i < m; i++ {
       b0 := (h[i][0] >> uint(bit[i])) & 1
       b1 := (h[i][1] >> uint(bit[i])) & 1
       res[i] = byte('0') + byte(b0) + byte(b1*2)
   }
   writer.Write(res)
   writer.WriteByte('\n')
}
