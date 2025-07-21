package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   // Read a 3x3 grid of press counts
   var a [3][3]int
   for i := 0; i < 3; i++ {
       for j := 0; j < 3; j++ {
           if _, err := fmt.Fscan(reader, &a[i][j]); err != nil {
               return
           }
       }
   }
   // Directions: self and four side-adjacent
   dirs := [][2]int{{0, 0}, {1, 0}, {-1, 0}, {0, 1}, {0, -1}}
   // Process each cell
   var out [3][3]byte
   for i := 0; i < 3; i++ {
       for j := 0; j < 3; j++ {
           sum := 0
           for _, d := range dirs {
               ni, nj := i+d[0], j+d[1]
               if ni >= 0 && ni < 3 && nj >= 0 && nj < 3 {
                   sum += a[ni][nj]
               }
           }
           // initial state is on (1), toggled sum times
           if sum%2 == 0 {
               out[i][j] = '1'
           } else {
               out[i][j] = '0'
           }
       }
   }
   // Print result
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   for i := 0; i < 3; i++ {
       for j := 0; j < 3; j++ {
           writer.WriteByte(out[i][j])
       }
       writer.WriteByte('\n')
   }
}
