package main

import (
   "bufio"
   "fmt"
   "os"
)

var (
   g      [26][]int
   color1 [26]int
   color  [26]int
   ans    []int
)

// detect cycle: return true if cycle found
func dfs1(v int) bool {
   color1[v] = 1
   for _, to := range g[v] {
       if color1[to] == 2 {
           continue
       }
       if color1[to] == 1 {
           return true
       }
       if dfs1(to) {
           return true
       }
   }
   color1[v] = 2
   return false
}

// topo sort DFS
func dfs(v int) {
   color[v] = 1
   for _, to := range g[v] {
       if color[to] == 2 {
           continue
       }
       if color[to] == 1 {
           return
       }
       dfs(to)
   }
   color[v] = 2
   ans = append(ans, v)
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var t int
   if _, err := fmt.Fscan(in, &t); err != nil {
       return
   }
   words := make([]string, t)
   for i := 0; i < t; i++ {
       fmt.Fscan(in, &words[i])
       words[i] += "#"
   }
   // build graph
   for i := 1; i < t; i++ {
       a := words[i-1]
       b := words[i]
       k := 0
       for k < len(a) && a[k] == b[k] {
           k++
       }
       if k < len(a) && a[k] == '#' {
           continue
       }
       if k < len(b) && b[k] == '#' {
           fmt.Fprintln(out, "Impossible")
           return
       }
       u := int(a[k] - 'a')
       v := int(b[k] - 'a')
       g[u] = append(g[u], v)
   }
   // detect cycles
   for i := 0; i < 26; i++ {
       if color1[i] == 0 {
           if dfs1(i) {
               fmt.Fprintln(out, "Impossible")
               return
           }
       }
   }
   // topo sort
   for i := 0; i < 26; i++ {
       if color[i] == 0 {
           dfs(i)
       }
   }
   // print reversed
   for i := len(ans) - 1; i >= 0; i-- {
       out.WriteByte(byte(ans[i]) + 'a')
   }
   out.WriteByte('\n')
}
