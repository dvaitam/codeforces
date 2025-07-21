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

   var n int
   fmt.Fscan(reader, &n)
   mid := n / 2
   for i := 0; i < n; i++ {
       d := i - mid
       if d < 0 {
           d = -d
       }
       stars := d
       ds := n - 2*d
       for j := 0; j < stars; j++ {
           writer.WriteByte('*')
       }
       for j := 0; j < ds; j++ {
           writer.WriteByte('D')
       }
       for j := 0; j < stars; j++ {
           writer.WriteByte('*')
       }
       writer.WriteByte('\n')
   }
}
