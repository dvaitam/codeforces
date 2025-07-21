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
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   xs := make([]int, n)
   ys := make([]int, n)
   maxColor := 0
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &xs[i], &ys[i])
       if xs[i] > maxColor {
           maxColor = xs[i]
       }
       if ys[i] > maxColor {
           maxColor = ys[i]
       }
   }
   count := make([]int, maxColor+1)
   for i := 0; i < n; i++ {
       count[xs[i]]++
   }
   for i := 0; i < n; i++ {
       conflicts := count[ys[i]]
       home := (n - 1) + conflicts
       away := (n - 1) - conflicts
       fmt.Fprintln(writer, home, away)
   }
}
