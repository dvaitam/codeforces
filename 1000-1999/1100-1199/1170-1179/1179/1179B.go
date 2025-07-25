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

   var n, m int
   fmt.Fscan(reader, &n, &m)
   total := n * m
   l, r := 0, total-1
   for l < r {
       x1 := l/m + 1
       y1 := l%m + 1
       fmt.Fprintln(writer, x1, y1)

       x2 := r/m + 1
       y2 := r%m + 1
       fmt.Fprintln(writer, x2, y2)

       l++
       r--
   }
   if l == r {
       x := l/m + 1
       y := l%m + 1
       fmt.Fprintln(writer, x, y)
   }
}
