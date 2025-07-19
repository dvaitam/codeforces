package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   var T int
   fmt.Fscan(reader, &T)
   for T > 0 {
       T--
       solve(reader, writer)
   }
}

func solve(reader *bufio.Reader, writer *bufio.Writer) {
   var n int
   fmt.Fscan(reader, &n)
   // operations pairs
   type op struct{ x, y int }
   ops := make([]op, 0, 0)
   // add function
   add := func(x, y int) {
       ops = append(ops, op{x, y})
   }
   for n > 4 {
       q := int(math.Sqrt(float64(n))) + 1
       for i := q + 1; i < n; i++ {
           add(i, n)
       }
       add(n, q)
       add(n, (n+q-1)/q)
       n = q
   }
   if n == 3 {
       add(3, 2)
       add(3, 2)
   } else if n == 4 {
       add(3, 4)
       add(4, 2)
       add(4, 2)
   }
   // output
   fmt.Fprintln(writer, len(ops))
   for _, p := range ops {
       fmt.Fprintln(writer, p.x, p.y)
   }
}
