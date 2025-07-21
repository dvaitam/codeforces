package main

import (
   "fmt"
)

func main() {
   var a [3][3]int
   // Read the 3x3 grid; zeros are on the main diagonal
   for i := 0; i < 3; i++ {
       for j := 0; j < 3; j++ {
           fmt.Scan(&a[i][j])
       }
   }
   // Extract known entries
   b := a[0][1]
   c := a[0][2]
   d := a[1][0]
   f := a[1][2]
   g := a[2][0]
   h := a[2][1]
   // Compute magic sum M = (b+c + d+f + g+h) / 2
   sum := b + c + d + f + g + h
   M := sum / 2
   // Fill the main diagonal
   a[0][0] = M - b - c
   a[1][1] = M - d - f
   a[2][2] = M - g - h
   // Output the magic square
   for i := 0; i < 3; i++ {
       fmt.Println(a[i][0], a[i][1], a[i][2])
   }
}
