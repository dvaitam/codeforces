package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   a := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &a[i])
   }
   // positions are numbered from end: a[0] is pos1
   c := make([]int64, n)
   for i := 0; i < n; i++ {
       // i from 0, position = i+1
       c[i] = a[i] + int64(i+1)
   }
   sort.Slice(c, func(i, j int) bool {
       return c[i] < c[j]
   })
   // check feasibility and compute result
   res := make([]int64, n)
   for i := 0; i < n; i++ {
       // position (i+1) must be <= c[i]
       if c[i] < int64(i+1) {
           fmt.Println(":(")
           return
       }
       res[i] = c[i] - int64(i+1)
   }
   // output res as final sequence at positions 1..n
   w := bufio.NewWriter(os.Stdout)
   defer w.Flush()
   for i, v := range res {
       if i > 0 {
           w.WriteByte(' ')
       }
       fmt.Fprint(w, v)
   }
   w.WriteByte('\n')
}
