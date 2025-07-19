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
   var w int
   if _, err := fmt.Fscan(in, &w); err != nil {
       return
   }
   for w > 0 {
       w--
       solve(in, out)
   }
}

func solve(in *bufio.Reader, out *bufio.Writer) {
   var n int
   fmt.Fscan(in, &n)
   t := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &t[i])
   }
   pref := make([]int, n)
   suf := make([]int, n)
   db := make([]int, n+2)
   mex := 0
   for i := 0; i < n; i++ {
       x := t[i]
       db[x]++
       for db[mex] > 0 {
           mex++
       }
       pref[i] = mex
   }
   for i := 0; i <= n; i++ {
       db[i] = 0
   }
   mex = 0
   for i := n - 1; i >= 0; i-- {
       x := t[i]
       db[x]++
       for db[mex] > 0 {
           mex++
       }
       suf[i] = mex
   }
   pos := -1
   for i := 0; i+1 < n; i++ {
       if pref[i] == suf[i+1] {
           pos = i
       }
   }
   if pos < 0 {
       fmt.Fprintln(out, -1)
   } else {
       fmt.Fprintln(out, 2)
       fmt.Fprintf(out, "1 %d\n", pos+1)
       fmt.Fprintf(out, "%d %d\n", pos+2, n)
   }
}
