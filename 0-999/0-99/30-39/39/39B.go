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
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   a := make([]int, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(in, &a[i])
   }

   // Map values to their positions
   pos := make(map[int][]int)
   for i := 1; i <= n; i++ {
       v := a[i]
       pos[v] = append(pos[v], i)
   }

   var years []int
   last := 0
   // For each target t = 1,2,... find earliest occurrence after last
   for t := 1; ; t++ {
       found := -1
       for _, idx := range pos[t] {
           if idx > last {
               found = idx
               break
           }
       }
       if found == -1 {
           break
       }
       years = append(years, found)
       last = found
   }

   k := len(years)
   if k == 0 {
       fmt.Fprintln(out, 0)
       return
   }
   fmt.Fprintln(out, k)
   for i, idx := range years {
       year := 2000 + idx
       if i+1 < k {
           fmt.Fprintf(out, "%d ", year)
       } else {
           fmt.Fprintf(out, "%d\n", year)
       }
   }
}
