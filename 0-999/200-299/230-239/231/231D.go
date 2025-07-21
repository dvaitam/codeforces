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

   var x, y, z int
   if _, err := fmt.Fscan(reader, &x, &y, &z); err != nil {
       return
   }
   var x1, y1, z1 int
   fmt.Fscan(reader, &x1, &y1, &z1)
   a := make([]int, 7)
   for i := 1; i <= 6; i++ {
       fmt.Fscan(reader, &a[i])
   }

   sum := 0
   if y < 0 {
       sum += a[1]
   }
   if y > y1 {
       sum += a[2]
   }
   if z < 0 {
       sum += a[3]
   }
   if z > z1 {
       sum += a[4]
   }
   if x < 0 {
       sum += a[5]
   }
   if x > x1 {
       sum += a[6]
   }

   fmt.Fprint(writer, sum)
}
