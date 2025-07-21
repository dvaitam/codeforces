package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var a int
   var v, u int
   var n int
   if _, err := fmt.Fscan(reader, &a, &v, &u, &n); err != nil {
       return
   }
   af := float64(a)
   vf := float64(v)
   uf := float64(u)
   // alpha = 1 - (u^2)/(v^2)
   alpha := 1.0 - (uf*uf)/(vf*vf)
   eps := 1e-12
   count := 0
   for i := 0; i < n; i++ {
       var xi, yi int
       fmt.Fscan(reader, &xi, &yi)
       xif := float64(xi)
       yif := float64(yi)
       // function f(x) = alpha*x^2 - 2*xi*x + xi^2 + yi^2
       f := func(x float64) float64 {
           return alpha*x*x - 2.0*xif*x + xif*xif + yif*yif
       }
       intercept := false
       if u == v {
           // linear case alpha == 0
           if xi > 0 {
               if f(af) <= eps {
                   intercept = true
               }
           } else {
               if f(0.0) <= eps {
                   intercept = true
               }
           }
       } else if alpha > 0 {
           // convex, minimum at x0 = xi/alpha
           x0 := xif / alpha
           if x0 < 0 {
               x0 = 0
           } else if x0 > af {
               x0 = af
           }
           if f(x0) <= eps {
               intercept = true
           }
       } else {
           // concave, check endpoints
           if f(0.0) <= eps || f(af) <= eps {
               intercept = true
           }
       }
       if intercept {
           count++
       }
   }
   fmt.Println(count)
}
