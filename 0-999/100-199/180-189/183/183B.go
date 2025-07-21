package main

import (
   "bufio"
   "fmt"
   "os"
)

// void type for map values
type void struct{}

// Line holds normalized line coefficients and set of flamingo indices
type Line struct {
   A, C int64
   pts  map[int]void
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   var m int
   fmt.Fscan(in, &n, &m)
   xs := make([]int64, m)
   ys := make([]int64, m)
   for i := 0; i < m; i++ {
       fmt.Fscan(in, &xs[i], &ys[i])
   }
   // map from line key to Line
   lines := make(map[string]*Line)
   for i := 0; i < m; i++ {
       for j := i + 1; j < m; j++ {
           x1, y1 := xs[i], ys[i]
           x2, y2 := xs[j], ys[j]
           // line: A x + B y + C = 0
           A := y2 - y1
           B := x1 - x2
           C := -(A*x1 + B*y1)
           // normalize coefficients
           if A == 0 && B == 0 {
               continue
           }
           g := gcd(abs64(A), abs64(B))
           g = gcd(g, abs64(C))
           A /= g; B /= g; C /= g
           // ensure unique representation
           if A < 0 || (A == 0 && B < 0) {
               A, B, C = -A, -B, -C
           }
           key := fmt.Sprintf("%d_%d_%d", A, B, C)
           ln, ok := lines[key]
           if !ok {
               ln = &Line{A: A, C: C, pts: make(map[int]void)}
               lines[key] = ln
           }
           ln.pts[i] = void{}
           ln.pts[j] = void{}
       }
   }
   // extra sightings beyond the first flamingo for each binocular
   extra := make(map[int]int)
   for _, ln := range lines {
       k := len(ln.pts)
       if k < 2 || ln.A == 0 {
           continue
       }
       // intersection with y=0: A*x + C = 0 => x = -C/A
       if (-ln.C)%ln.A != 0 {
           continue
       }
       x0 := int((-ln.C) / ln.A)
       if x0 < 1 || x0 > n {
           continue
       }
       // each binocular initially sees one flamingo; extra is k-1
       if extra[x0] < k-1 {
           extra[x0] = k-1
       }
   }
   // compute total: n binoculars each see at least 1 flamingo
   ans := int64(n)
   for _, v := range extra {
       ans += int64(v)
   }
   fmt.Println(ans)
}

// abs64 returns absolute value of x
func abs64(x int64) int64 {
   if x < 0 {
       return -x
   }
   return x
}

// gcd returns greatest common divisor of a and b
func gcd(a, b int64) int64 {
   for b != 0 {
       a, b = b, a%b
   }
   return a
}
