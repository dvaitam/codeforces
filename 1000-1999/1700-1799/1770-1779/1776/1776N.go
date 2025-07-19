package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   var n int
   fmt.Fscan(reader, &n)
   var str string
   fmt.Fscan(reader, &str)
   // Counting '<' and others in comparisons str
   n1, n2 := 0, 0
   // str length is typically n-1
   m := len(str)
   for i := 0; i < m; i++ {
       if str[i] == '<' {
           n2++
       } else {
           n1++
       }
   }
   // Prepare boundaries
   a := make([]int, n+2)
   b := make([]int, n+2)
   x, y := n1, 0
   for i := 0; i < m; i++ {
       if a[x] < y {
           a[x] = y
       }
       if b[y] < x {
           b[y] = x
       }
       if str[i] == '<' {
           y++
       } else {
           x--
       }
   }
   // Final boundary
   if a[x] < y {
       a[x] = y
   }
   if b[y] < x {
       b[y] = x
   }
   const B = 1000
   f := make([]float64, B+1)
   f[0] = 1.0
   x, y = n1, 0
   delt := 0.0
   // Dynamic programming with scaling
   for i := 0; i < n; i++ {
       mx := 0.0
       for j := 0; j <= B; j++ {
           xx := x - j
           yy := y - j
           if xx >= 0 && yy >= 0 {
               denom := float64(a[xx] + b[yy] - xx - yy + 1)
               f[j] /= denom
               if f[j] > mx {
                   mx = f[j]
               }
           } else {
               f[j] = 0
           }
       }
       delt += math.Log2(mx)
       for j := 0; j <= B; j++ {
           f[j] /= mx
       }
       // Transition: use '<' if available, else treat as '>'
       if i < m && str[i] == '<' {
           for j := B; j >= 1; j-- {
               f[j] += f[j-1]
           }
           y++
       } else {
           for j := 1; j <= B; j++ {
               f[j-1] += f[j]
           }
           x--
       }
   }
   // Compute final result
   res := math.Log2(f[0]) + delt
   for i := 1; i <= n; i++ {
       res += math.Log2(float64(i))
   }
   fmt.Fprintf(writer, "%.10f\n", res)
}
