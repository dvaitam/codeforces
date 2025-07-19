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
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &a[i])
       a[i]--
   }
   b := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &b[i])
       b[i]--
   }
   // label cycles in a
   vis := make([]bool, n)
   aLab := make([]int, n)
   n1 := 0
   for i := 0; i < n; i++ {
       if vis[i] {
           continue
       }
       for j := i; !vis[j]; j = a[j] {
           vis[j] = true
           aLab[j] = n1
       }
       n1++
   }
   // label cycles in b
   for i := range vis {
       vis[i] = false
   }
   bLab := make([]int, n)
   n2 := 0
   for i := 0; i < n; i++ {
       if vis[i] {
           continue
       }
       for j := i; !vis[j]; j = b[j] {
           vis[j] = true
           bLab[j] = n2
       }
       n2++
   }
   // build bipartite graph
   e := make([][]int, n1)
   for i := 0; i < n; i++ {
       u := aLab[i]
       v := bLab[i]
       e[u] = append(e[u], v)
   }
   // bipartite matching
   x := make([]int, n1)
   for i := range x {
       x[i] = -1
   }
   y := make([]int, n2)
   for i := range y {
       y[i] = -1
   }
   visL := make([]bool, n1)
   var dfs func(int) bool
   dfs = func(u int) bool {
       if visL[u] {
           return false
       }
       visL[u] = true
       for _, v := range e[u] {
           if y[v] == -1 || dfs(y[v]) {
               x[u] = v
               y[v] = u
               return true
           }
       }
       return false
   }
   match := 0
   for i := 0; i < n1; i++ {
       for j := range visL {
           visL[j] = false
       }
       if dfs(i) {
           match++
       }
   }
   // operations are elements where their edge is unmatched
   k := n - match
   fmt.Fprintln(out, k)
   first := true
   for i := 0; i < n; i++ {
       if x[aLab[i]] != bLab[i] {
           if !first {
               out.WriteByte(' ')
           }
           first = false
           fmt.Fprint(out, i+1)
       }
   }
   fmt.Fprintln(out)
}
