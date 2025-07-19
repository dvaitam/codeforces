package main

import (
   "bufio"
   "fmt"
   "os"
)

type pair struct {
   first, second int
}

var (
   n        int
   adj      [][]int
   ind      []int
   f        [][2]pair
   col      []int
   reader   = bufio.NewReader(os.Stdin)
   writer   = bufio.NewWriter(os.Stdout)
)

func dfs(u, p int) {
   f[u][0] = pair{0, 1}
   f[u][1] = pair{1, ind[u]}
   for _, v := range adj[u] {
       if v == p {
           continue
       }
       dfs(v, u)
       // if u selected (flag=1), children cannot be selected
       f[u][1].first += f[v][0].first
       f[u][1].second += f[v][0].second
       // if u not selected (flag=0), choose best of child
       if f[v][0].first > f[v][1].first {
           f[u][0].first += f[v][0].first
           f[u][0].second += f[v][0].second
       } else if f[v][0].first < f[v][1].first {
           f[u][0].first += f[v][1].first
           f[u][0].second += f[v][1].second
       } else {
           f[u][0].first += f[v][0].first
           if f[v][0].second < f[v][1].second {
               f[u][0].second += f[v][0].second
           } else {
               f[u][0].second += f[v][1].second
           }
       }
   }
}

func work(u, p int, flag bool) {
   if flag {
       col[u] = ind[u]
   } else {
       col[u] = 1
   }
   for _, v := range adj[u] {
       if v == p {
           continue
       }
       if flag {
           work(v, u, false)
       } else {
           if f[v][0].first > f[v][1].first {
               work(v, u, false)
           } else if f[v][0].first < f[v][1].first {
               work(v, u, true)
           } else {
               if f[v][0].second < f[v][1].second {
                   work(v, u, false)
               } else {
                   work(v, u, true)
               }
           }
       }
   }
}

func main() {
   defer writer.Flush()
   fmt.Fscan(reader, &n)
   adj = make([][]int, n+1)
   ind = make([]int, n+1)
   for i := 1; i < n; i++ {
       var u, v int
       fmt.Fscan(reader, &u, &v)
       adj[u] = append(adj[u], v)
       adj[v] = append(adj[v], u)
       ind[u]++
       ind[v]++
   }
   col = make([]int, n+1)
   if n == 2 {
       fmt.Fprintln(writer, "2 2")
       fmt.Fprintln(writer, "1 1")
       return
   }
   f = make([][2]pair, n+1)
   dfs(1, 0)
   var bestFirst, bestSecond int
   var rootFlag bool
   if f[1][0].first > f[1][1].first {
       bestFirst = f[1][0].first
       bestSecond = f[1][0].second
       rootFlag = false
   } else if f[1][0].first < f[1][1].first {
       bestFirst = f[1][1].first
       bestSecond = f[1][1].second
       rootFlag = true
   } else {
       bestFirst = f[1][0].first
       if f[1][0].second < f[1][1].second {
           bestSecond = f[1][0].second
           rootFlag = false
       } else {
           bestSecond = f[1][1].second
           rootFlag = true
       }
   }
   fmt.Fprintln(writer, bestFirst, bestSecond)
   work(1, 0, rootFlag)
   for i := 1; i <= n; i++ {
       if i > 1 {
           writer.WriteByte(' ')
       }
       fmt.Fprint(writer, col[i])
   }
   writer.WriteByte('\n')
}
