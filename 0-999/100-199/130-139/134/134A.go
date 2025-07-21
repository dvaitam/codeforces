package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   a := make([]int, n)
   var sum int64
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
       sum += int64(a[i])
   }
   // Check if average (including this) is integer
   if sum%int64(n) != 0 {
       fmt.Println(0)
       return
   }
   target := sum / int64(n)
   // Collect indices where a[i] == target
   var res []int
   for i, v := range a {
       if int64(v) == target {
           res = append(res, i+1)
       }
   }
   // Output
   w := bufio.NewWriter(os.Stdout)
   defer w.Flush()
   fmt.Fprintln(w, len(res))
   if len(res) > 0 {
       for i, idx := range res {
           if i > 0 {
               w.WriteByte(' ')
           }
           fmt.Fprint(w, idx)
       }
       fmt.Fprintln(w)
   }
}
