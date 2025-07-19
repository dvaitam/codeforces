package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
)

type pt struct{ x, y float64 }

const (
   pi  = math.Acos(-1)
   eps = 1e-8
)

func sub(p, q pt) pt { return pt{p.x - q.x, p.y - q.y} }
func dot(p, q pt) float64 { return p.x*q.x + p.y*q.y }
func length(p pt) float64 { return math.Sqrt(dot(p, p)) }
func angle(p, q, r pt) float64 {
   v1 := sub(p, q)
   v2 := sub(r, q)
   return math.Acos(dot(v1, v2) / (length(v1) * length(v2)))
}

func area(a float64) float64 {
   if math.Abs(a-pi/2) < eps {
       return 3*pi/32 - (a-pi/2)/4
   }
   if math.Abs(a-pi) < 1e-4 {
       return (pi - a) / 6
   }
   return (math.Tan(pi/2-a) + (pi-a)*(2+1/(math.Sin(a)*math.Sin(a))) - 4*math.Sin(2*a))/16 + math.Sin(2*a)/4
}

func xx(a, t float64) float64 {
   return (math.Cos(t)*math.Sin(a+t)*math.Sin(a+t) - math.Cos(a)*math.Sin(t)*math.Sin(t)*math.Cos(a+t)) / (math.Sin(a) * math.Sin(a))
}

func yy(a, t float64) float64 {
   return -math.Sin(t)*math.Sin(t)*math.Cos(a+t) / math.Sin(a)
}

// solve xx(a, t) = x by binary search
func calc(a, x float64) float64 {
   lo, hi := 0.0, pi-a
   for hi-lo > eps {
       mi := (lo + hi) / 2
       if xx(a, mi) > x {
           lo = mi
       } else {
           hi = mi
       }
   }
   return lo
}

func integrate(a, t float64) float64 {
   r := -8*t + 2*math.Sin(4*t) + math.Cos(2*a)*(4*t - 3*math.Sin(2*t) - math.Sin(4*t) + math.Sin(6*t)) + 8*(3+2*math.Cos(2*t))*math.Pow(math.Sin(t), 4)*math.Sin(2*a)
   return -r / (64 * math.Sin(a) * math.Sin(a))
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   a := make([]pt, n)
   for i := 0; i < n; i++ {
       var xi, yi int
       fmt.Fscan(in, &xi, &yi)
       a[i] = pt{float64(xi), float64(yi)}
   }
   // extend for cyclic access
   a = append(a, a[0], a[1], a[2])

   if n == 4 {
       for i := 0; i < 2; i++ {
           if a[i].x == a[i+1].x && a[i+2].x == a[i+3].x && (math.Abs(a[i+2].x-a[i+1].x) == 1 || math.Abs(a[i].y-a[i+1].y) == 1) {
               dx := math.Abs(a[i+2].x - a[i+1].x)
               dy := math.Abs(a[i].y - a[i+1].y)
               fmt.Println(dx * dy)
               return
           }
       }
   }

   ans := 0.0
   for i := 1; i <= n; i++ {
       al := angle(a[i-1], a[i], a[i+1])
       ans += area(al)
   }
   for i := 1; i <= n; i++ {
       d := length(sub(a[i], a[i+1]))
       if d >= 2-eps {
           continue
       }
       al := angle(a[i-1], a[i], a[i+1])
       be := angle(a[i], a[i+1], a[i+2])
       lo, hi := 0.0, pi-al
       for hi-lo > eps {
           mi := (lo + hi) / 2
           if d-xx(al, mi) <= 1 && yy(al, mi) < yy(be, calc(be, d-xx(al, mi))) {
               lo = mi
           } else {
               hi = mi
           }
       }
       ans -= integrate(al, lo)
       ans -= integrate(be, calc(be, d-xx(al, lo)))
   }
   fmt.Printf("%.15f\n", ans)
}
