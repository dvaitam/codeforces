package main

import (
   "bufio"
   "fmt"
   "os"
)

// Convex hull trick for lines of form y = m*x + b, slopes sorted increasing, queries x sorted increasing
type line struct {
   m, b int64
   t    int // start day
}

// check if l2 is unnecessary between l1 and l3
func bad(l1, l2, l3 line) bool {
   // intersection(l1,l2) >= intersection(l2,l3)
   return (l2.b-l1.b)*(l2.m-l3.m) >= (l3.b-l2.b)*(l1.m-l2.m)
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, m int
   fmt.Fscan(reader, &n, &m)
   a := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   // read queries
   kqs := make([]int64, m)
   for i := 0; i < m; i++ {
       fmt.Fscan(reader, &kqs[i])
   }
   // handle n == 1: all spoilage zero
   if n <= 1 {
       for i := 0; i < m; i++ {
           if i > 0 {
               fmt.Fprint(writer, ' ')
           }
           fmt.Fprint(writer, 0)
       }
       fmt.Fprintln(writer)
       return
   }
   // prefix sums
   p := make([]int64, n)
   p[0] = a[0]
   for i := 1; i < n; i++ {
       p[i] = p[i-1] + a[i]
   }
   totalBeforeN := p[n-2] // sum a[0..n-2]
   // build hull
   hull := make([]line, 0, n)
   for t := 1; t <= n-1; t++ {
       // day t: index t-1 in a, but t from 1..n-1 corresponds to a[t-1]
       // suffix sum from t to n-1 is totalBeforeN - p[t-2]
       var pref int64
       if t-2 >= 0 {
           pref = p[t-2]
       }
       b := totalBeforeN - pref
       mSlope := int64(-(n - t))
       ln := line{m: mSlope, b: b, t: t}
       // add ln to hull
       for len(hull) >= 2 && bad(hull[len(hull)-2], hull[len(hull)-1], ln) {
           hull = hull[:len(hull)-1]
       }
       hull = append(hull, ln)
   }
   // answer queries
   res := make([]int, m)
   ptr := 0
   for i, k := range kqs {
       // move pointer to best line
       for ptr+1 < len(hull) {
           y1 := hull[ptr].m*k + hull[ptr].b
           y2 := hull[ptr+1].m*k + hull[ptr+1].b
           if y2 >= y1 {
               ptr++
           } else {
               break
           }
       }
       y := hull[ptr].m*k + hull[ptr].b
       if y <= 0 {
           res[i] = 0
       } else {
           res[i] = n - hull[ptr].t
       }
   }
   // output
   for i, v := range res {
       if i > 0 {
           fmt.Fprint(writer, ' ')
       }
       fmt.Fprint(writer, v)
   }
   fmt.Fprintln(writer)
}
