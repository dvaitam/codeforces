package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   var t int
   fmt.Fscan(reader, &t)
   for t > 0 {
       t--
       var n, k, x int
       fmt.Fscan(reader, &n, &k, &x)
       xorSum := 0
       for i := 1; i <= n; i++ {
           xorSum ^= i
       }
       if (k%2 == 1 && xorSum != x) || (k%2 == 0 && xorSum != 0) {
           fmt.Fprintln(writer, "NO")
           continue
       }
       used := make([]bool, n+1)
       groups := make([][]int, 0)
       for i := 1; i <= n; i++ {
           if i == x {
               used[x] = true
               groups = append(groups, []int{x})
               continue
           }
           if used[i] {
               continue
           }
           xi := x ^ i
           if xi <= n {
               groups = append(groups, []int{i, xi})
               used[i] = true
               used[xi] = true
           }
       }
       if len(groups) >= k {
           fmt.Fprintln(writer, "YES")
           // merge extra groups into first
           for h := k; h < len(groups); h++ {
               groups[0] = append(groups[0], groups[h]...)
           }
           // add unused elements
           for i := 1; i <= n; i++ {
               if !used[i] {
                   groups[0] = append(groups[0], i)
               }
           }
           for y := 0; y < k; y++ {
               grp := groups[y]
               fmt.Fprint(writer, len(grp))
               for _, e := range grp {
                   fmt.Fprint(writer, " ", e)
               }
               fmt.Fprint(writer, "\n")
           }
       } else if len(groups) == k-1 {
           fmt.Fprintln(writer, "YES")
           // create last group with remaining elements
           last := make([]int, 0)
           for i := 1; i <= n; i++ {
               if !used[i] {
                   last = append(last, i)
               }
           }
           groups = append(groups, last)
           for y := 0; y < k; y++ {
               grp := groups[y]
               fmt.Fprint(writer, len(grp))
               for _, e := range grp {
                   fmt.Fprint(writer, " ", e)
               }
               fmt.Fprint(writer, "\n")
           }
       } else {
           fmt.Fprintln(writer, "NO")
       }
   }
}
