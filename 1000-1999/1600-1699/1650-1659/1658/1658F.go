package main

import (
   "bufio"
   "fmt"
   "os"
)

// gcd returns the greatest common divisor of a and b.
func gcd(a, b int) int {
   for b != 0 {
       a, b = b, a%b
   }
   return a
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var cas int
   if _, err := fmt.Fscan(reader, &cas); err != nil {
       return
   }
   for t := 0; t < cas; t++ {
       var n, m int
       var s string
       fmt.Fscan(reader, &n, &m)
       fmt.Fscan(reader, &s)
       // Count total ones
       totalOnes := 0
       for i := 0; i < n; i++ {
           if s[i] == '1' {
               totalOnes++
           }
       }
       g := gcd(n, m)
       // Check feasibility: totalOnes / n must equal b / m
       if totalOnes%(n/g) != 0 {
           fmt.Fprintln(writer, -1)
           continue
       }
       need := totalOnes/(n/g) * (m / g)
       // Sliding window of length m on circular string
       cur := 0
       for i := 0; i < m; i++ {
           if s[i] == '1' {
               cur++
           }
       }
       L, R := -1, -1
       found := false
       for i := 0; i < n; i++ {
           if cur == need {
               // non-wrap segment
               if i+m <= n {
                   fmt.Fprintln(writer, 1)
                   fmt.Fprintf(writer, "%d %d\n", i+1, i+m)
                   found = true
                   break
               }
               // record wrap positions
               L = i
               R = (i + m - 1) % n
           }
           // slide window: remove s[i], add s[(i+m)%n]
           if s[i] == '1' {
               cur--
           }
           if s[(i+m)%n] == '1' {
               cur++
           }
       }
       if found {
           continue
       }
       if L < 0 {
           fmt.Fprintln(writer, -1)
           continue
       }
       // wrap-around uses two segments
       fmt.Fprintln(writer, 2)
       fmt.Fprintf(writer, "%d %d\n", 1, R+1)
       fmt.Fprintf(writer, "%d %d\n", L+1, n)
   }
}
