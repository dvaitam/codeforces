package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   fmt.Fscan(reader, &n)
   var digits []int
   if n == 0 {
       digits = append(digits, 0)
   } else {
       for n > 0 {
           digits = append(digits, n%10)
           n /= 10
       }
   }
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   for _, d := range digits {
       // go-dama (5-bead)
       if d >= 5 {
           writer.WriteString("-O|")
       } else {
           writer.WriteString("O-|")
       }
       // ichi-damas (four 1-beads)
       x := d % 5
       for i := 0; i < x; i++ {
           writer.WriteByte('O')
       }
       writer.WriteByte('-')
       for i := 0; i < 4-x; i++ {
           writer.WriteByte('O')
       }
       writer.WriteByte('\n')
   }
}
