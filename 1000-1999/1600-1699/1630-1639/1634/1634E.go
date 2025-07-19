package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

type Pai struct {
   x, xh int
}

type Tpi struct {
   x, xh1, xh2 int
}

type Adj struct {
   to, wSelf, wOther, eid int
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var m int
   if _, err := fmt.Fscan(reader, &m); err != nil {
       return
   }
   nArr := make([]int, m)
   a := make([][]Pai, m)
   ansInt := make([][]int, m)
   ta := make([]Tpi, 0)

   for i := 0; i < m; i++ {
       fmt.Fscan(reader, &nArr[i])
       n := nArr[i]
       a[i] = make([]Pai, n)
       ansInt[i] = make([]int, n)
       for j := 0; j < n; j++ {
           var x int
           fmt.Fscan(reader, &x)
           a[i][j] = Pai{x, j}
       }
       sort.Slice(a[i], func(p, q int) bool { return a[i][p].x < a[i][q].x })
       for j := 0; j < n; {
           if j+1 < n && a[i][j].x == a[i][j+1].x {
               ansInt[i][a[i][j].xh] = 0
               ansInt[i][a[i][j+1].xh] = 1
               j += 2
           } else {
               ta = append(ta, Tpi{a[i][j].x, i, a[i][j].xh})
               j++
           }
       }
   }
   tot := len(ta)
   if tot%2 != 0 {
       fmt.Fprint(writer, "NO")
       return
   }
   sort.Slice(ta, func(i, j int) bool { return ta[i].x < ta[j].x })
   adj := make([][]Adj, m)
   pairs := tot / 2
   for i := 0; i < pairs; i++ {
       t1 := ta[2*i]
       t2 := ta[2*i+1]
       if t1.x != t2.x {
           fmt.Fprint(writer, "NO")
           return
       }
       u, v := t1.xh1, t2.xh1
       w1, w2 := t1.xh2, t2.xh2
       eid := i
       adj[u] = append(adj[u], Adj{to: v, wSelf: w1, wOther: w2, eid: eid})
       adj[v] = append(adj[v], Adj{to: u, wSelf: w2, wOther: w1, eid: eid})
   }
   visited := make([]bool, pairs)
   idxEdge := make([]int, m)
   var dfs func(w, st int)
   dfs = func(w, st int) {
       for {
           if idxEdge[w] >= len(adj[w]) {
               return
           }
           e := adj[w][idxEdge[w]]
           idxEdge[w]++
           if visited[e.eid] {
               continue
           }
           visited[e.eid] = true
           ansInt[w][e.wSelf] = 0
           ansInt[e.to][e.wOther] = 1
           if e.to != st {
               dfs(e.to, st)
           }
           return
       }
   }
   for i := 0; i < m; i++ {
       for idxEdge[i] < len(adj[i]) {
           e := adj[i][idxEdge[i]]
           idxEdge[i]++
           if visited[e.eid] {
               continue
           }
           visited[e.eid] = true
           ansInt[i][e.wSelf] = 0
           ansInt[e.to][e.wOther] = 1
           dfs(e.to, i)
       }
   }
   fmt.Fprintln(writer, "YES")
   for i := 0; i < m; i++ {
       for j := 0; j < nArr[i]; j++ {
           if ansInt[i][j] == 0 {
               writer.WriteByte('L')
           } else {
               writer.WriteByte('R')
           }
       }
       writer.WriteByte('\n')
   }
}
