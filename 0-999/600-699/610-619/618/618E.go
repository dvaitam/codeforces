package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
)

var (
   n, m int
   cirCos, cirSin [360]float64
   S   []int
   vx  []float64
   vy  []float64
   Ang []int
   Len []int
   out *bufio.Writer
)

func update(o int) {
   lc, rc := o<<1, o<<1|1
   S[o] = (S[lc] + S[rc]) % 360
   // rotate right child vector by S[lc]
   c := cirCos[S[lc]]
   s := cirSin[S[lc]]
   x := vx[lc] + vx[rc]*c - vy[rc]*s
   y := vy[lc] + vx[rc]*s + vy[rc]*c
   vx[o], vy[o] = x, y
}

func modify(o, l, r, k int) {
   if l == r {
       ang := Ang[k]
       S[o] = ang
       vx[o] = float64(Len[k]) * cirCos[ang]
       vy[o] = float64(Len[k]) * cirSin[ang]
       return
   }
   mid := (l + r) >> 1
   if k <= mid {
       modify(o<<1, l, mid, k)
   } else {
       modify(o<<1|1, mid+1, r, k)
   }
   update(o)
}

func build(o, l, r int) {
   if l == r {
       S[o] = Ang[l] // zero
       vx[o], vy[o] = 1.0, 0.0
   } else {
       mid := (l + r) >> 1
       build(o<<1, l, mid)
       build(o<<1|1, mid+1, r)
       update(o)
   }
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   out = bufio.NewWriter(os.Stdout)
   defer out.Flush()
   fmt.Fscan(reader, &n, &m)
   // precompute cos and sin
   const pi = math.Pi
   for i := 0; i < 360; i++ {
       rad := float64(i) / 360.0 * 2 * pi
       cirCos[i] = math.Cos(rad)
       cirSin[i] = math.Sin(rad)
   }
   // init arrays
   size := 4 * n
   S = make([]int, size+5)
   vx = make([]float64, size+5)
   vy = make([]float64, size+5)
   Ang = make([]int, n+1)
   Len = make([]int, n+1)
   for i := 1; i <= n; i++ {
       Len[i] = 1
   }
   build(1, 1, n)
   for i := 0; i < m; i++ {
       var typ, y, z int
       fmt.Fscan(reader, &typ, &y, &z)
       if typ == 1 {
           Len[y] += z
           modify(1, 1, n, y)
       } else {
           // rotate clockwise by z: subtract z
           Ang[y] = (Ang[y] - (z % 360) + 360) % 360
           modify(1, 1, n, y)
       }
       fmt.Fprintf(out, "%.5f %.5f\n", vx[1], vy[1])
   }
}
