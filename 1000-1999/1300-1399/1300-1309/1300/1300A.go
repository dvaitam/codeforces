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
   var t int
   fmt.Fscan(in, &t)
   for ti := 0; ti < t; ti++ {
       var n int
       fmt.Fscan(in, &n)
       sum := 0
       zeros := 0
       for i := 0; i < n; i++ {
           var x int
           fmt.Fscan(in, &x)
           sum += x
           if x == 0 {
               zeros++
           }
       }
       // Each zero requires one step to become 1
       sum += zeros
       ans := zeros
       // If sum becomes zero after fixing zeros, one more step needed
       if sum == 0 {
           ans++
       }
       fmt.Fprintln(out, ans)
   }
}
