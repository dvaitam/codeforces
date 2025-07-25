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
   if n%2 == 1 {
       fmt.Fprintln(out, 0)
   } else {
       // For even n, number of ways is 2^(n/2)
       ans := 1 << uint(n/2)
       fmt.Fprintln(out, ans)
   }
}
