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

   var x int64
   fmt.Fscan(in, &x)
   type op struct{ u, v int64; t int }
   var ans []op

   for x > 1 {
       // find largest i such that (1<<(i+1)) <= x
       i := 0
       for (int64(1) << uint(i+1)) <= x {
           ans = append(ans, op{x << uint(i), x << uint(i), 0})
           i++
       }
       a := x << uint(i)
       b := a ^ x
       ans = append(ans, op{a, x, 1})
       for j := 0; j < i; j++ {
           ans = append(ans, op{x << uint(j), a, 0})
           ans = append(ans, op{x << uint(j), b, 0})
           a += x << uint(j)
           b += x << uint(j)
       }
       ans = append(ans, op{a, b, 1})
       x = a ^ b
   }

   fmt.Fprintln(out, len(ans))
   for _, e := range ans {
       var sym byte
       if e.t == 0 {
           sym = '+'
       } else {
           sym = '^'
       }
       fmt.Fprintf(out, "%d %c %d\n", e.u, sym, e.v)
   }
}
