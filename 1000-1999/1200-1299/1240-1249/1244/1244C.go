package main

import (
   "bufio"
   "fmt"
   "os"
)

func gcd(a, b int64) int64 {
   for b != 0 {
       a, b = b, a%b
   }
   return a
}

// exgcd returns gcd(a,b) and x,y such that a*x + b*y = gcd
func exgcd(a, b int64) (int64, int64, int64) {
   if b == 0 {
       return a, 1, 0
   }
   g, x1, y1 := exgcd(b, a%b)
   x := y1
   y := x1 - (a/b)*y1
   return g, x, y
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, p, w, d int64
   if _, err := fmt.Fscan(in, &n, &p, &w, &d); err != nil {
       return
   }
   g := gcd(w, d)
   if p%g != 0 {
       fmt.Println(-1)
       return
   }
   w0 := w / g
   d0 := d / g
   p0 := p / g

   // solve d0 * y ≡ p0 mod w0
   _, inv, _ := exgcd(d0, w0)
   inv = (inv%w0 + w0) % w0
   y0 := (p0 % w0 * inv) % w0
   x0 := (p0 - y0*d0) / w0

   // parameter j: x = x0 + j*d0, y = y0 - j*w0
   // ensure x>=0: j >= jMin
   jMin := int64(0)
   if x0 < 0 {
       jMin = (-x0 + d0 - 1) / d0
   }
   // ensure y>=0: j <= jMax
   jMax := y0 / w0
   if jMin > jMax {
       fmt.Println(-1)
       return
   }
   // choose j to minimize x+y => delta = j*(d0 - w0)
   var j int64
   if d0 > w0 {
       j = jMin
   } else {
       j = jMax
   }
   x := x0 + j*d0
   y := y0 - j*w0
   if x < 0 || y < 0 || x+y > n {
       fmt.Println(-1)
       return
   }
   // losses = n - x - y
   fmt.Printf("%d %d %d\n", x, y, n-x-y)
}
