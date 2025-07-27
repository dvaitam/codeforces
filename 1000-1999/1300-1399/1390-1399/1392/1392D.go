package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var t int
   if _, err := fmt.Fscan(in, &t); err != nil {
       return
   }
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   for t > 0 {
       t--
       var n int
       var s string
       fmt.Fscan(in, &n)
       fmt.Fscan(in, &s)
       // collect runs of consecutive identical attacks
       runs := make([]int, 0, n)
       cur := 1
       for i := 1; i < n; i++ {
           if s[i] == s[i-1] {
               cur++
           } else {
               runs = append(runs, cur)
               cur = 1
           }
       }
       runs = append(runs, cur)
       // merge first and last runs if circular and same character
       if len(runs) > 1 && s[0] == s[n-1] {
           runs[0] += runs[len(runs)-1]
           runs = runs[:len(runs)-1]
       }
       ans := 0
       for _, r := range runs {
           ans += r / 2
       }
       fmt.Fprintln(out, ans)
   }
}
