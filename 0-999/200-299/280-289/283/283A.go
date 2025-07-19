package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n int
   fmt.Fscan(in, &n)
   a := make([]int64, n+5)
   c := make([]int64, n+5)
   num := 1
   var sum int64
   a[1] = 0
   c[1] = 0

   for i := 0; i < n; i++ {
       var com int
       fmt.Fscan(in, &com)
       if com == 1 {
           var x int
           var y int64
           fmt.Fscan(in, &x, &y)
           sum += int64(x) * y
           c[x] += y
       } else if com == 2 {
           var x int64
           fmt.Fscan(in, &x)
           num++
           a[num] = x
           c[num] = 0
           sum += x
       } else if com == 3 {
           sum -= a[num] + c[num]
           c[num-1] += c[num]
           num--
       }
       avg := float64(sum) / float64(num)
       fmt.Fprintf(out, "%.8f\n", avg)
   }
}
