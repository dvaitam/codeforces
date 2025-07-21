package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var A, B, C int
   if _, err := fmt.Fscan(reader, &A, &B, &C); err != nil {
       return
   }
   var sols [][2]int
   // X and Y must be positive integers (>=1)
   for x := 1; A*x < C; x++ {
       rem := C - A*x
       if rem%B != 0 {
           continue
       }
       y := rem / B
       if y >= 1 {
           sols = append(sols, [2]int{x, y})
       }
   }
   // print results
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fprintln(writer, len(sols))
   for _, p := range sols {
       fmt.Fprintf(writer, "%d %d\n", p[0], p[1])
   }
}
