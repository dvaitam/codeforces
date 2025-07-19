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

func min(a, b int) int {
   if a < b {
       return a
   }
   return b
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n1, m1 int
   if _, err := fmt.Fscan(reader, &n1, &m1); err != nil {
       return
   }
   mat1 := make([][]int, n1)
   for i := 0; i < n1; i++ {
       var line string
       fmt.Fscan(reader, &line)
       row := make([]int, m1)
       for j, ch := range line {
           row[j] = int(ch - '0')
       }
       mat1[i] = row
   }
   var n2, m2 int
   fmt.Fscan(reader, &n2, &m2)
   mat2 := make([][]int, n2)
   for i := 0; i < n2; i++ {
       var line string
       fmt.Fscan(reader, &line)
       row := make([]int, m2)
       for j, ch := range line {
           row[j] = int(ch - '0')
       }
       mat2[i] = row
   }

   limitX := max(n1, n2)
   limitY := max(m1, m2)
   maxSum, bestX, bestY := 0, 0, 0
   for x := -limitX; x < limitX; x++ {
       for y := -limitY; y < limitY; y++ {
           sum := 0
           rStart := max(0, -x)
           rEnd := min(n1, n2-x)
           cStart := max(0, -y)
           cEnd := min(m1, m2-y)
           for i := rStart; i < rEnd; i++ {
               for j := cStart; j < cEnd; j++ {
                   sum += mat1[i][j] * mat2[i+x][j+y]
               }
           }
           if sum > maxSum {
               maxSum = sum
               bestX = x
               bestY = y
           }
       }
   }
   fmt.Printf("%d %d", bestX, bestY)
}
