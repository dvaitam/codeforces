package main

import (
   "bufio"
   "fmt"
   "os"
)

func max(a, b int) int {
   if a > b {
       return a
   }
   return b
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var t int
   fmt.Fscan(reader, &t)
   for t > 0 {
       t--
       var n int
       fmt.Fscan(reader, &n)
       // allocate arrays of size n+2
       a := make([]int, n+2)
       L := make([]int, n+2)
       R := make([]int, n+2)
       maxL := make([]int, n+2)
       maxR := make([]int, n+2)
       for i := 1; i <= n; i++ {
           fmt.Fscan(reader, &a[i])
       }
       // Precompute L and R (next smaller to left/right)
       // Stack for indices
       stack := make([]int, 0, n+2)
       a[0] = 0
       stack = append(stack, 0)
       for i := 1; i <= n; i++ {
           for len(stack) > 0 && a[stack[len(stack)-1]] >= a[i] {
               stack = stack[:len(stack)-1]
           }
           L[i] = stack[len(stack)-1] + 1
           stack = append(stack, i)
       }
       // next smaller to right
       stack = stack[:0]
       a[n+1] = 0
       stack = append(stack, n+1)
       for i := n; i >= 1; i-- {
           for len(stack) > 0 && a[stack[len(stack)-1]] >= a[i] {
               stack = stack[:len(stack)-1]
           }
           R[i] = stack[len(stack)-1] - 1
           stack = append(stack, i)
       }
       // prefix max
       maxL[0] = 0
       for i := 1; i <= n; i++ {
           maxL[i] = max(maxL[i-1], a[i])
       }
       // suffix max
       maxR[n+1] = 0
       for i := n; i >= 1; i-- {
           maxR[i] = max(maxR[i+1], a[i])
       }
       // try to find partition
       found := false
       for i := 1; i <= n; i++ {
           u := maxL[L[i]-1]
           v := maxR[R[i]+1]
           Li, Ri := L[i], R[i]
           if u < a[i] && Li != i && a[Li] == a[i] {
               u = a[i]
               Li++
           }
           if v < a[i] && Ri != i && a[Ri] == a[i] {
               v = a[i]
               Ri--
           }
           if u == a[i] && v == a[i] {
               x := Li - 1
               y := Ri - Li + 1
               z := n - x - y
               fmt.Fprintln(writer, "YES")
               fmt.Fprintln(writer, x, y, z)
               found = true
               break
           }
       }
       if !found {
           fmt.Fprintln(writer, "NO")
       }
   }
}
