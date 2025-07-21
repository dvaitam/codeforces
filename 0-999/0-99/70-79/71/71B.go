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

   var n, k, t int
   if _, err := fmt.Fscan(in, &n, &k, &t); err != nil {
       return
   }
   // total saturation sum
   S := t * n * k / 100
   full := S / k
   rem := S % k
   // build progress bar
   a := make([]int, n)
   for i := 0; i < n; i++ {
       switch {
       case i < full:
           a[i] = k
       case i == full:
           a[i] = rem
       default:
           a[i] = 0
       }
   }
   // output result
   for i, v := range a {
       if i > 0 {
           fmt.Fprint(out, " ")
       }
       fmt.Fprint(out, v)
   }
   fmt.Fprintln(out)
}
