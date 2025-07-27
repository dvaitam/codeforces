package main

import (
   "bufio"
   "fmt"
   "os"
)

func query(writer *bufio.Writer, reader *bufio.Reader, l, r int) int {
   fmt.Fprintf(writer, "? %d %d\n", l, r)
   writer.Flush()
   var x int
   fmt.Fscan(reader, &x)
   return x
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, t int
   if _, err := fmt.Fscan(reader, &n, &t); err != nil {
      return
   }
   for ; t > 0; t-- {
      var k int
      fmt.Fscan(reader, &k)
      l, r := 1, n
      for l < r {
         m := (l + r) / 2
         ones := query(writer, reader, 1, m)
         zeros := m - ones
         if zeros >= k {
            r = m
         } else {
            l = m + 1
         }
      }
      fmt.Fprintf(writer, "! %d\n", l)
      writer.Flush()
      var verdict string
      fmt.Fscan(reader, &verdict)
      if verdict != "Correct" && verdict != "" {
         return
      }
   }
}

