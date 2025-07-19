package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   var n int
   if _, err := fmt.Fscan(os.Stdin, &n); err != nil {
       return
   }
   if n < 3 {
       fmt.Println(-1)
       return
   }
   // initialize grid
   grid := make([][]int, n)
   for i := range grid {
       grid[i] = make([]int, n)
   }
   // base 3x3 pattern
   uwu := [3][3]int{
       {7, 6, 9},
       {8, 2, 5},
       {1, 4, 3},
   }
   base := n*n - 9
   // place initial cells
   for i := 0; i < 3; i++ {
       for j := 0; j < 3; j++ {
           grid[i][j] = uwu[i][j] + base
       }
   }
   // fill remaining cells
   i, j, cur, mode := 3, 0, base, 0
   for cur > 0 {
       grid[i][j] = cur
       // move according to mode
       switch mode {
       case 0:
           j++
       case 1:
           j--
       case 2:
           i++
       default:
           i--
       }
       // adjust mode and bounds
       if i == j {
           mode = 3 - mode
       } else if i < 0 {
           i = 0
           j++
           mode = 2
       } else if j < 0 {
           i++
           j = 0
           mode = 0
       }
       cur--
   }
   // output
   w := bufio.NewWriter(os.Stdout)
   defer w.Flush()
   for ii := 0; ii < n; ii++ {
       for jj := 0; jj < n; jj++ {
           fmt.Fprint(w, grid[ii][jj])
           if jj+1 < n {
               w.WriteByte(' ')
           }
       }
       w.WriteByte('\n')
   }
}
