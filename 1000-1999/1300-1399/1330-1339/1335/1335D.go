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

   var T int
   if _, err := fmt.Fscan(reader, &T); err != nil {
       return
   }
   // positions to change
   coords := [][2]int{{0, 0}, {1, 3}, {2, 6}, {3, 1}, {4, 4}, {5, 7}, {6, 2}, {7, 5}, {8, 8}}
   for tc := 0; tc < T; tc++ {
       s := make([][]byte, 9)
       for i := 0; i < 9; i++ {
           var line string
           fmt.Fscan(reader, &line)
           s[i] = []byte(line)
       }
       // apply changes
       for _, c := range coords {
           x, y := c[0], c[1]
           if s[x][y] == '9' {
               s[x][y] = '1'
           } else {
               s[x][y]++
           }
       }
       // output
       for i := 0; i < 9; i++ {
           writer.Write(s[i])
           writer.WriteByte('\n')
       }
   }
}
