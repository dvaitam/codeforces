package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, p int
   if _, err := fmt.Fscan(in, &n, &p); err != nil {
       return
   }
   xs := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &xs[i])
   }
   // DP dimensions: at most 32 small numbers
   const maxT = 32
   // f[i][xorMask][mod] = reachable
   f := make([][][]bool, maxT)
   g := make([][][]bool, maxT)
   for i := 0; i < maxT; i++ {
       f[i] = make([][]bool, maxT)
       g[i] = make([][]bool, maxT)
       for j := 0; j < maxT; j++ {
           f[i][j] = make([]bool, p)
           g[i][j] = make([]bool, p)
       }
   }
   add := func(cur, y int) int {
       v := cur * 10
       if y > 9 {
           v *= 10
       }
       return (v + y) % p
   }
   // Find previous mod
   findPrev := func(y, curMod, curXor, idx int) int {
       for t := 0; t < p; t++ {
           if f[idx][curXor][t] && add(t, y) == curMod {
               return t
           }
       }
       return 0
   }
   f[0][0][0] = true
   num := make([]int, maxT)
   pos := make([]int, maxT)
   tot := 0
   found := false
   for i, x := range xs {
       if x < maxT {
           // add this small number
           tot++
           num[tot] = x
           pos[tot] = i + 1
           for j := 0; j < maxT && !found; j++ {
               for k := 0; k < p; k++ {
                   if f[tot-1][j][k] {
                       // without taking
                       f[tot][j][k] = true
                       // taking
                       t := add(k, x)
                       nx := j ^ x
                       f[tot][nx][t] = true
                       g[tot][nx][t] = true
                       if j == x && t == 0 {
                           found = true
                           break
                       }
                   }
               }
           }
       }
       if found {
           break
       }
   }
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   if !found {
       fmt.Fprintln(out, "No")
       return
   }
   fmt.Fprintln(out, "Yes")
   // reconstruct
   curXor, curMod := 0, 0
   var ansPos []int
   for i := tot; i > 0; i-- {
       if g[i][curXor][curMod] {
           // this was taken to reach state
           ansPos = append(ansPos, pos[i])
           y := num[i]
           curXor ^= y
           curMod = findPrev(y, curMod, curXor, i-1)
       }
   }
   m := len(ansPos)
   fmt.Fprintln(out, m)
   // print in original order
   for i := m - 1; i >= 0; i-- {
       if i < m-1 {
           out.WriteByte(' ')
       }
       fmt.Fprint(out, ansPos[i])
   }
   fmt.Fprintln(out)
}
