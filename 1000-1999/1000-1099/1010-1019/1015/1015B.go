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
   var s, t string
   fmt.Fscan(in, &s, &t)
   // Check if s and t have same multiset of characters
   if len(s) != n || len(t) != n {
       fmt.Println(-1)
       return
   }
   a := []rune(s)
   b := []rune(t)
   sa := make([]rune, n)
   sb := make([]rune, n)
   copy(sa, a)
   copy(sb, b)
   sort.Slice(sa, func(i, j int) bool { return sa[i] < sa[j] })
   sort.Slice(sb, func(i, j int) bool { return sb[i] < sb[j] })
   for i := 0; i < n; i++ {
       if sa[i] != sb[i] {
           fmt.Println(-1)
           return
       }
   }
   // Build sequence of adjacent swaps
   ops := make([]int, 0)
   for i := 0; i < n; i++ {
       if a[i] != b[i] {
           // find the matching character
           j := i + 1
           for ; j < n; j++ {
               if a[j] == b[i] {
                   break
               }
           }
           // bubble it to position i
           for k := j; k > i; k-- {
               a[k], a[k-1] = a[k-1], a[k]
               ops = append(ops, k)
           }
       }
   }
   // Output
   w := bufio.NewWriter(os.Stdout)
   defer w.Flush()
   fmt.Fprintln(w, len(ops))
   for _, op := range ops {
       fmt.Fprint(w, op, " ")
   }
   fmt.Fprintln(w)
}
