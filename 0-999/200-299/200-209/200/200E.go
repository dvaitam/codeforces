package main

import (
   "bufio"
   "fmt"
   "os"
)

func abs64(x int64) int64 {
   if x < 0 {
       return -x
   }
   return x
}

// exgcd returns x, y, gcd such that a*x + b*y = gcd
func exgcd(a, b int64) (int64, int64, int64) {
   if b == 0 {
       return 1, 0, a
   }
   x1, y1, d := exgcd(b, a%b)
   return y1, x1 - (a/b)*y1, d
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, s int64
   if _, err := fmt.Fscan(reader, &n, &s); err != nil {
       return
   }
   var f [3]int64
   for i := int64(0); i < n; i++ {
       var a int64
       fmt.Fscan(reader, &a)
       // a is in {3,4,5}
       f[a-3]++
   }

   const inf = int64(1e18)
   ans := inf
   var ansX, ansY, ansZ int64
   var step1, step2 int64
   var px, py, pdiv int64

   for i := int64(0); i*f[1] <= s; i++ {
       // solve f[0]*x + f[2]*z = s - i*f[1]
       x0, y0, d := exgcd(f[0], f[2])
       px, py, pdiv = x0, y0, d
       remain := s - i*f[1]
       if pdiv != 0 && remain%pdiv == 0 {
           px *= remain / pdiv
           py *= remain / pdiv
           lcm := f[0] * f[2] / pdiv
           t1 := lcm / f[0]
           t2 := lcm / f[2]
           // shift solution to minimal non-negative x
           shift := -px / t1
           px += t1 * shift
           py -= t2 * shift
           if px < 0 {
               px += t1
               py -= t2
           }
           step1 = t1
           step2 = t2

           pain := func(r2 int64) {
               start := r2 - 5
               if start < 0 {
                   start = 0
               }
               end := r2 + 5
               for r := start; r <= end; r++ {
                   c0 := px + r*step1
                   c1 := i
                   c2 := py - r*step2
                   if c0 <= c1 && c1 <= c2 {
                       cost := abs64(c0*f[0] - c1*f[1]) + abs64(c1*f[1] - c2*f[2])
                       if cost < ans {
                           ans = cost
                           ansX, ansY, ansZ = c0, c1, c2
                       }
                   }
               }
           }
           // try around candidate r2 values
           pain(0)
           pain((py - i) / step2)
           pain((i - px) / step1)
           if step1*f[0] != 0 {
               pain((i*f[1] - px*f[0]) / (step1 * f[0]))
           }
           if step2*f[2] != 0 {
               pain((i*f[1] - py*f[2]) / (step2 * f[2]))
           }
       }
   }

   if ans == inf {
       fmt.Fprint(writer, -1)
   } else {
       fmt.Fprintf(writer, "%d %d %d", ansX, ansY, ansZ)
   }
}
