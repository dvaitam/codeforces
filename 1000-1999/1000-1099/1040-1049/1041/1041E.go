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
   cnt := make([]int, n+1)
   for i := 1; i < n; i++ {
       var a, b int
       fmt.Fscan(in, &a, &b)
       cnt[a]++
       cnt[b]++
   }
   // check degree of n
   if cnt[n] != n-1 {
       fmt.Print("NO")
       return
   }
   // check cnt[i] <= i
   for i := 1; i <= n; i++ {
       if cnt[i] > i {
           fmt.Print("NO")
           return
       }
   }
   // collect zeros
   all := make([]int, 0, n)
   for i := 1; i <= n; i++ {
       if cnt[i] == 0 {
           all = append(all, i)
       }
   }
   // sort descending
   sort.Sort(sort.Reverse(sort.IntSlice(all)))
   sol := make([][2]int, 0, n)
   prv := n
   idx := 0
   // build solution
   for i := n - 1; i >= 1; i-- {
       if cnt[i] == 0 {
           continue
       }
       // connect extra edges
       for cnt[i] > 1 {
           if idx >= len(all) {
               fmt.Print("NO")
               return
           }
           u := prv
           v := all[idx]
           if v > i {
               fmt.Print("NO")
               return
           }
           sol = append(sol, [2]int{u, v})
           prv = v
           idx++
           cnt[i]--
       }
       // connect the last
       sol = append(sol, [2]int{i, prv})
       prv = i
   }
   // output
   w := bufio.NewWriter(os.Stdout)
   defer w.Flush()
   fmt.Fprintln(w, "YES")
   for _, p := range sol {
       fmt.Fprintf(w, "%d %d\n", p[0], p[1])
   }
}
