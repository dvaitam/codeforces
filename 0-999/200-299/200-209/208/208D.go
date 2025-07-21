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
   p := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &p[i])
   }
   var a, b, c, d, e int64
   fmt.Fscan(reader, &a, &b, &c, &d, &e)

   counts := make([]int64, 5)
   var points int64
   for _, v := range p {
       points += v
       for {
           switch {
           case points >= e:
               points -= e
               counts[4]++
           case points >= d:
               points -= d
               counts[3]++
           case points >= c:
               points -= c
               counts[2]++
           case points >= b:
               points -= b
               counts[1]++
           case points >= a:
               points -= a
               counts[0]++
           default:
               break
           }
           // break out of loop when no more prizes
           if points < a {
               break
           }
       }
   }

   // output counts
   fmt.Fprintf(writer, "%d %d %d %d %d\n", counts[0], counts[1], counts[2], counts[3], counts[4])
   // output remaining points
   fmt.Fprintln(writer, points)
}
