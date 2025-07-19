package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   rd := bufio.NewReader(os.Stdin)
   var n int
   var S int64
   fmt.Fscan(rd, &n, &S)
   x := make([]int64, n)
   y := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(rd, &x[i], &y[i])
   }
   // Function to compute doubled area (absolute) of triangle a,b,c
   SS := func(a, b, c int) int64 {
       v := (x[a]-x[b])*(y[a]+y[b]) + (x[b]-x[c])*(y[b]+y[c]) + (x[c]-x[a])*(y[c]+y[a])
       if v < 0 {
           return -v
       }
       return v
   }
   // start with first three points
   a, b, c := 0, 1, 2
   count := 0
   for {
       cur := SS(a, b, c)
       newa, newb, newc := a, b, c
       // try improving c
       for i := 0; i < n; i++ {
           if v := SS(a, b, i); v > cur {
               cur = v
               newc = i
           }
       }
       // try improving b
       for i := 0; i < n; i++ {
           if v := SS(a, i, newc); v > cur {
               cur = v
               newb = i
           }
       }
       // try improving a
       for i := 0; i < n; i++ {
           if v := SS(i, newb, newc); v > cur {
               cur = v
               newa = i
           }
       }
       if (newa == a && newb == b && newc == c) || count > 1000 {
           break
       }
       a, b, c = newa, newb, newc
       count++
   }
   // output three points of maximum perimeter triangle
   w := bufio.NewWriter(os.Stdout)
   defer w.Flush()
   fmt.Fprintf(w, "%d %d\n", x[a]+x[b]-x[c], y[a]+y[b]-y[c])
   fmt.Fprintf(w, "%d %d\n", x[a]+x[c]-x[b], y[a]+y[c]-y[b])
   fmt.Fprintf(w, "%d %d", x[b]+x[c]-x[a], y[b]+y[c]-y[a])
}
