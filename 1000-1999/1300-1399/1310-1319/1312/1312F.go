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

   var T int
   fmt.Fscan(reader, &T)
   // Precompute Grundy
   const H = 500
   // t: 0=none/mixed,1=infantry,2=cavalry
   g := make([][3]int, H+1)
   // values for moves to subtract
   var x, y, z int
   // We'll fill g once x,y,z known per test? x,y,z differ per test.
   // Instead we need to precompute per test inside loop or cache per unique x,y,z
   // For simplicity, process each test separately
   for tc := 0; tc < T; tc++ {
       var n int
       fmt.Fscan(reader, &n, &x, &y, &z)
       a := make([]int64, n)
       for i := 0; i < n; i++ {
           fmt.Fscan(reader, &a[i])
       }
       // compute g up to H
       for h := 0; h <= H; h++ {
           for t := 0; t < 3; t++ {
               g[h][t] = 0
           }
       }
       for h := 1; h <= H; h++ {
           for t := 0; t < 3; t++ {
               var used [4]bool
               // mixed
               h2 := h - x
               if h2 < 0 {
                   h2 = 0
               }
               used[g[h2][0]] = true
               // infantry
               if t != 1 {
                   h2 = h - y
                   if h2 < 0 {
                       h2 = 0
                   }
                   used[g[h2][1]] = true
               }
               // cavalry
               if t != 2 {
                   h2 = h - z
                   if h2 < 0 {
                       h2 = 0
                   }
                   used[g[h2][2]] = true
               }
               mex := 0
               for used[mex] {
                   mex++
               }
               g[h][t] = mex
           }
       }
       // detect period
       base := 200
       period := -1
       for p := 1; p <= 100; p++ {
           ok := true
           for h := base; h+p <= H; h++ {
               for t := 0; t < 3; t++ {
                   if g[h][t] != g[h+p][t] {
                       ok = false
                       break
                   }
               }
               if !ok {
                   break
               }
           }
           if ok {
               period = p
               break
           }
       }
       if period == -1 {
           period = 1
       }
       // helper to get g for any h
       getg := func(h0 int64, t int) int {
           if h0 <= int64(H) {
               return g[int(h0)][t]
           }
           if h0 < int64(base) {
               // <base but >H? shouldn't happen since H>base
               return g[int(h0)][t]
           }
           idx := base + int((h0-int64(base))%int64(period))
           return g[idx][t]
       }
       // compute initial xor
       xor := 0
       for i := 0; i < n; i++ {
           xor ^= getg(a[i], 0)
       }
       if xor == 0 {
           fmt.Fprintln(writer, 0)
           continue
       }
       // count winning moves
       cnt := 0
       for i := 0; i < n; i++ {
           ai := a[i]
           gi := getg(ai, 0)
           // mixed
           h2 := ai - int64(x)
           if h2 < 0 {
               h2 = 0
           }
           g2 := getg(h2, 0)
           if xor^gi^g2 == 0 {
               cnt++
           }
           // infantry
           h2 = ai - int64(y)
           if h2 < 0 {
               h2 = 0
           }
           g2 = getg(h2, 1)
           if xor^gi^g2 == 0 {
               cnt++
           }
           // cavalry
           h2 = ai - int64(z)
           if h2 < 0 {
               h2 = 0
           }
           g2 = getg(h2, 2)
           if xor^gi^g2 == 0 {
               cnt++
           }
       }
       fmt.Fprintln(writer, cnt)
   }
}
