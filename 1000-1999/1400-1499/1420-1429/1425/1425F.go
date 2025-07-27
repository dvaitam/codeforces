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

   var N int
   // read number of cages
   fmt.Fscan(in, &N)

   ans := make([]int, N+1)
   x := make([]int, N+1)

   // ask for sum of [1,2] and [2,3]
   x[1] = query(1, 2, in, out)
   x[2] = query(2, 3, in, out)
   // ask for sum of [1,3]
   sum13 := query(1, 3, in, out)
   // compute first three
   ans[1] = sum13 - x[2]
   ans[2] = x[1] - ans[1]
   ans[3] = x[2] - ans[2]

   // ask remaining adjacent sums and compute
   for i := 3; i < N; i++ {
       x[i] = query(i, i+1, in, out)
       ans[i+1] = x[i] - ans[i]
   }

   // output result
   fmt.Fprint(out, "!")
   for i := 1; i <= N; i++ {
       fmt.Fprintf(out, " %d", ans[i])
   }
   fmt.Fprintln(out)
   out.Flush()
}

// query asks sum of flamingoes in range [l, r]
func query(l, r int, in *bufio.Reader, out *bufio.Writer) int {
   fmt.Fprintf(out, "? %d %d\n", l, r)
   out.Flush()
   var v int
   fmt.Fscan(in, &v)
   return v
}
