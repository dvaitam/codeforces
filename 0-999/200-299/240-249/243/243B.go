package main

import (
   "bufio"
   "fmt"
   "os"
)

var (
   n, m, h, t int
   a [][]int
   v []bool
   out *bufio.Writer
)

func check(x, y, dx, dy, get int) bool {
   // try assigning heads in two configurations
   nx := dx - get - t
   ny := dy - get - h
   if -nx-ny <= get && dx >= t && dy >= h {
       out.WriteString("YES\n")
       out.WriteString(fmt.Sprintf("%d %d\n", y, x))
       // pick h neighbors of y
       p := 0
       for _, u := range a[y] {
           if !v[u] && u != x {
               out.WriteString(fmt.Sprintf("%d ", u))
               p++
               if p == h {
                   break
               }
           }
       }
       if p < h {
           for _, u := range a[y] {
               if v[u] && u != x {
                   out.WriteString(fmt.Sprintf("%d ", u))
                   p++
                   v[u] = false
                   if p == h {
                       break
                   }
               }
           }
       }
       out.WriteString("\n")
       // pick t neighbors of x
       p = 0
       for _, u := range a[x] {
           if v[u] && u != y {
               out.WriteString(fmt.Sprintf("%d ", u))
               p++
               if p == t {
                   break
               }
           }
       }
       out.WriteString("\n")
       out.Flush()
       os.Exit(0)
   }
   nx = dx - get - h
   ny = dy - get - t
   if -nx-ny <= get && dx >= h && dy >= t {
       out.WriteString("YES\n")
       out.WriteString(fmt.Sprintf("%d %d\n", x, y))
       // pick t neighbors of y into g
       g := make([]int, 0, t)
       p := 0
       for _, u := range a[y] {
           if !v[u] && u != x {
               g = append(g, u)
               p++
               if p == t {
                   break
               }
           }
       }
       if p < t {
           for _, u := range a[y] {
               if v[u] && u != x {
                   g = append(g, u)
                   v[u] = false
                   p++
                   if p == t {
                       break
                   }
               }
           }
       }
       // pick h neighbors of x
       p = 0
       for _, u := range a[x] {
           if v[u] && u != y {
               out.WriteString(fmt.Sprintf("%d ", u))
               p++
               if p == h {
                   break
               }
           }
       }
       out.WriteString("\n")
       for i := 0; i < t; i++ {
           out.WriteString(fmt.Sprintf("%d ", g[i]))
       }
       out.WriteString("\n")
       out.Flush()
       os.Exit(0)
   }
   return false
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out = bufio.NewWriter(os.Stdout)
   defer out.Flush()
   fmt.Fscan(in, &n, &m, &h, &t)
   a = make([][]int, n+1)
   for i := 0; i < m; i++ {
       var x, y int
       fmt.Fscan(in, &x, &y)
       a[x] = append(a[x], y)
       a[y] = append(a[y], x)
   }
   v = make([]bool, n+1)
   for i := 1; i <= n; i++ {
       dx := len(a[i]) - 1
       // mark neighbors of i
       for _, u := range a[i] {
           v[u] = true
       }
       for _, y := range a[i] {
           dy := len(a[y]) - 1
           if dx >= h+t && dy >= h+t {
               out.WriteString("YES\n")
               out.WriteString(fmt.Sprintf("%d %d\n", i, y))
               // pick h neighbors of i except y
               cnt := 0
               for _, u := range a[i] {
                   if u != y {
                       out.WriteString(fmt.Sprintf("%d ", u))
                       v[u] = true
                       cnt++
                       if cnt == h {
                           break
                       }
                   }
               }
               out.WriteString("\n")
               // pick t neighbors of y avoiding used
               cnt = 0
               for _, u := range a[y] {
                   if u != i && !v[u] {
                       out.WriteString(fmt.Sprintf("%d ", u))
                       cnt++
                       if cnt == t {
                           break
                       }
                   }
               }
               if cnt < t {
                   for _, u := range a[y] {
                       if u != i && v[u] {
                           out.WriteString(fmt.Sprintf("%d ", u))
                           cnt++
                           if cnt == t {
                               break
                           }
                       }
                   }
               }
               out.WriteString("\n")
               out.Flush()
               return
           }
           if dy < h+t {
               get := 0
               for _, u := range a[y] {
                   if u != i && v[u] {
                       get++
                   }
               }
               if check(i, y, dx, dy, get) {
                   return
               }
           }
       }
       // clear marks
       for _, u := range a[i] {
           v[u] = false
       }
   }
   out.WriteString("NO\n")
}
