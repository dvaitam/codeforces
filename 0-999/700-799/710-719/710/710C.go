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

   var n int
   _, _ = fmt.Fscan(reader, &n)

   x := make([][]int, n)
   for i := 0; i < n; i++ {
       x[i] = make([]int, n)
   }

   oddc := 1
   evenc := 2
   flag := false

   // fill center column with odds
   for i := 0; i < n; i++ {
       x[i][n/2] = oddc
       oddc += 2
   }

   // fill center row with remaining odds
   for j := 0; j < n; j++ {
       if x[n/2][j] == 0 {
           x[n/2][j] = oddc
           oddc += 2
       }
   }

   // fill four quadrants symmetrically
   for i := 0; i < n/2; i++ {
       for j := 0; j < n/2; j++ {
           if oddc-2 == n*n {
               flag = true
               break
           }
           x[i][j] = oddc; oddc += 2
           x[n-i-1][j] = oddc; oddc += 2
           x[i][n-j-1] = oddc; oddc += 2
           x[n-i-1][n-j-1] = oddc; oddc += 2
       }
       if flag {
           break
       }
   }

   // output matrix, filling zeros with evens
   for i := 0; i < n; i++ {
       for j := 0; j < n; j++ {
           if x[i][j] == 0 {
               fmt.Fprint(writer, evenc)
               evenc += 2
           } else {
               fmt.Fprint(writer, x[i][j])
           }
           if j == n-1 {
               fmt.Fprint(writer, "\n")
           } else {
               fmt.Fprint(writer, " ")
           }
       }
   }
}
