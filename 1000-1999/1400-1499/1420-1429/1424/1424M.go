package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, k int
   if _, err := fmt.Fscan(reader, &n, &k); err != nil {
       return
   }
   type st struct {
       p int
       s []string
   }
   a := make([]st, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i].p)
       a[i].s = make([]string, k)
       for j := 0; j < k; j++ {
           fmt.Fscan(reader, &a[i].s[j])
       }
   }
   sort.Slice(a, func(i, j int) bool { return a[i].p < a[j].p })
   // collect strings and mark used letters
   mark := make([]bool, 26)
   var list []string
   for i := 0; i < n; i++ {
       for j := 0; j < k; j++ {
           str := a[i].s[j]
           list = append(list, str)
           for _, ch := range str {
               mark[ch-'a'] = true
           }
       }
   }
   // build graph of constraints
   adj := make([][]int, 26)
   flag := false
   for i := 1; i < len(list); i++ {
       last := list[i-1]
       curr := list[i]
       minLen := len(last)
       if len(curr) < minLen {
           minLen = len(curr)
       }
       found := false
       for r := 0; r < minLen; r++ {
           if last[r] != curr[r] {
               u := int(last[r] - 'a')
               v := int(curr[r] - 'a')
               adj[u] = append(adj[u], v)
               found = true
               break
           }
       }
       if !found && len(last) > len(curr) {
           flag = true
       }
   }
   if flag {
       fmt.Fprintln(writer, "IMPOSSIBLE")
       return
   }
   // topological sort
   color := make([]int, 26)
   var topsort []int
   var dfs func(u int)
   dfs = func(u int) {
       if flag {
           return
       }
       color[u] = 1
       for _, v := range adj[u] {
           if color[v] == 0 {
               dfs(v)
           } else if color[v] == 1 {
               flag = true
               return
           }
       }
       color[u] = 2
       topsort = append(topsort, u)
   }
   for u := 0; u < 26; u++ {
       if mark[u] && color[u] == 0 {
           dfs(u)
       }
   }
   if flag {
       fmt.Fprintln(writer, "IMPOSSIBLE")
       return
   }
   // output order
   for i, j := 0, len(topsort)-1; i < j; i, j = i+1, j-1 {
       topsort[i], topsort[j] = topsort[j], topsort[i]
   }
   for _, u := range topsort {
       writer.WriteByte(byte(u) + 'a')
   }
   writer.WriteByte('\n')
}
