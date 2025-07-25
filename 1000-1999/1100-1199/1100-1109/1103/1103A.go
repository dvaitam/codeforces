package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var s string
   fmt.Fscan(reader, &s)
   x, y := 0, 0
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   for _, ch := range s {
       if ch == '0' {
           // vertical tile: place in column 1, alternating between rows 1 and 3
           row := 1 + 2*(x%2)
           col := 1
           fmt.Fprintf(writer, "%d %d\n", row, col)
           x++
       } else {
           // horizontal tile: place in row 4, alternating between columns 1 and 3
           row := 4
           col := 1 + 2*(y%2)
           fmt.Fprintf(writer, "%d %d\n", row, col)
           y++
       }
   }
}
