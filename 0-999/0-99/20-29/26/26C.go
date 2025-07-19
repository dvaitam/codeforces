package main

import (
   "fmt"
   "os"
)

func e() {
   fmt.Println("IMPOSSIBLE")
   os.Exit(0)
}

func main() {
   var n, m, a, b, c int
   if _, err := fmt.Scan(&n, &m, &a, &b, &c); err != nil {
       return
   }
   // both odd dimensions impossible
   if n%2 == 1 && m%2 == 1 {
       e()
   }
   // fill odd row with horizontal planks
   if n%2 == 1 {
       a -= m / 2
   }
   // fill odd column with vertical planks
   if m%2 == 1 {
       b -= n / 2
   }
   if a < 0 || b < 0 {
       e()
   }
   // initialize grid with 2x2 blocks labeled 'a' and 'b'
   s := make([][]byte, n)
   for i := 0; i < n; i++ {
       s[i] = make([]byte, m)
       for j := 0; j < m; j++ {
           if (i/2)%2 == (j/2)%2 {
               s[i][j] = 'a'
           } else {
               s[i][j] = 'b'
           }
       }
   }
   // determine excess 2x2 planks to break
   c -= (n / 2) * (m / 2)
   // break extra 2x2 blocks into smaller planks
   for i := 0; i < n-n%2; i += 2 {
       for j := 0; j < m-m%2; j += 2 {
           if c < 0 {
               if a > 1 {
                   a -= 2
                   c++
                   // two horizontal 1x2 planks on this 2x2 block (top and bottom rows)
                   s[i][j] += 2
                   s[i][j+1] += 2
               } else if b > 1 {
                   b -= 2
                   c++
                   // two vertical 2x1 planks on this 2x2 block (left and right columns)
                   s[i][j] += 2
                   s[i+1][j] += 2
               } else {
                   e()
               }
           }
       }
   }
   // output grid
   for i := 0; i < n; i++ {
       for j := 0; j < m; j++ {
           fmt.Printf("%c", s[i][j])
       }
       fmt.Println()
   }
}
