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
   var t, n int
   if _, err := fmt.Fscan(in, &t); err != nil {
       return
   }
   for ; t > 0; t-- {
       fmt.Fscan(in, &n)
       switch n {
       case 1:
           fmt.Fprintln(out, 1)
           continue
       case 2:
           fmt.Fprintln(out, "1 2")
           continue
       }
       ans := make([]int, 0, n)
       ans = append(ans, 2)
       half1 := (n - 3) / 2
       for i := 1; i <= half1; i++ {
           ans = append(ans, i+3)
       }
       ans = append(ans, 1)
       half2 := (n - 2) / 2
       for i := 1; i <= half2; i++ {
           ans = append(ans, i+3+half1)
       }
       ans = append(ans, 3)
       // output the sequence
       fmt.Fprint(out, ans[0])
       for i := 1; i < len(ans); i++ {
           fmt.Fprint(out, " ", ans[i])
       }
       fmt.Fprint(out, "\n")
   }
}
