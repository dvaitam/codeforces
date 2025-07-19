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

   var t, n int
   if _, err := fmt.Fscan(reader, &t); err != nil {
       return
   }
   for ; t > 0; t-- {
       fmt.Fscan(reader, &n)
       // initialize
       a := make([][]int, n)
       for i := 0; i < n; i++ {
           a[i] = make([]int, n)
       }
       // fill diagonals
       left := 2
       for i := 0; i < 2*n; i++ {
           if i == 0 || i == 2*(n-1) {
               continue
           }
           for j := 0; j < n; j++ {
               x := i - j
               y := j
               if x >= 0 && x < n && y >= 0 && y < n {
                   a[x][y] = left
                   left++
               }
           }
       }
       // special placements
       if n > 0 {
           a[0][0] = 0
       }
       if n >= 2 {
           a[n-1][n-1] = 1
       }
       // increment all
       for i := 0; i < n; i++ {
           for j := 0; j < n; j++ {
               a[i][j]++
           }
       }
       // check validity
       if !check(a, n) {
           fmt.Fprintln(writer, -1)
           continue
       }
       // output
       for i := 0; i < n; i++ {
           for j := 0; j < n; j++ {
               if j > 0 {
                   writer.WriteByte(' ')
               }
               fmt.Fprint(writer, a[i][j])
           }
           writer.WriteByte('\n')
       }
   }
}

func check(a [][]int, n int) bool {
   dirs := [][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
   for i := 0; i < n; i++ {
       for j := 0; j < n; j++ {
           for _, d := range dirs {
               ni := i + d[0]
               nj := j + d[1]
               if ni < 0 || ni >= n || nj < 0 || nj >= n {
                   continue
               }
               if abs(a[ni][nj]-a[i][j]) == 1 {
                   return false
               }
           }
       }
   }
   return true
}

func abs(x int) int {
   if x < 0 {
       return -x
   }
   return x
}
