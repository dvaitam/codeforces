package main

import (
   "bufio"
   "fmt"
   "os"
)

// up2 returns ceil(n/2)
func up2(n int) int {
   return n/2 + n%2
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var t, n int
   if _, err := fmt.Fscan(in, &t); err != nil {
       return
   }
   for ; t > 0; t-- {
       fmt.Fscan(in, &n)
       // Always 2 operations as per solution
       fmt.Fprintln(out, 2)
       cur := n
       for i := n - 1; i >= 1; i-- {
           fmt.Fprint(out, i, " ", cur, '\n')
           cur = up2(i + cur)
       }
   }
}
