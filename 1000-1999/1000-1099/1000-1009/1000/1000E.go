package main

import (
   "bufio"
   "fmt"
   "os"
)

var (
   reader = bufio.NewReader(os.Stdin)
   writer = bufio.NewWriter(os.Stdout)
)

func readInt() (int, error) {
   var x int
   var s byte
   var err error
   // skip non-numeric
   for {
       s, err = reader.ReadByte()
       if err != nil {
           return 0, err
       }
       if (s >= '0' && s <= '9') || s == '-' {
           break
       }
   }
   neg := false
   if s == '-' {
       neg = true
       s, err = reader.ReadByte()
       if err != nil {
           return 0, err
       }
   }
   for ; s >= '0' && s <= '9'; s, err = reader.ReadByte() {
       if err != nil {
           break
       }
       x = x*10 + int(s - '0')
   }
   if neg {
       x = -x
   }
   return x, nil
}

func main() {
   defer writer.Flush()
   n, _ := readInt()
   m, _ := readInt()
   // original nodes are 1..n
   // read edges
   u := make([]int, m)
   v := make([]int, m)
   adj := make([][]edge, n+1)
   for i := 0; i < m; i++ {
       ui, _ := readInt()
       vi, _ := readInt()
       u[i], v[i] = ui, vi
       // add both directions, id = 2*i, 2*i+1
       adj[ui] = append(adj[ui], edge{to: vi, id: 2 * i})
       adj[vi] = append(adj[vi], edge{to: ui, id: 2*i + 1})
   }
   tin := make([]int, n+1)
   low := make([]int, n+1)
   isBridge := make([]bool, 2*m)
   timer := 0

   var dfsBridge func(int, int)
   dfsBridge = func(x, peid int) {
       timer++
       tin[x] = timer
       low[x] = timer
       for _, e := range adj[x] {
           if e.id == peid {
               continue
           }
           y := e.to
           if tin[y] == 0 {
               dfsBridge(y, e.id)
               if low[y] < low[x] {
                   low[x] = low[y]
               }
               if low[y] > tin[x] {
                   isBridge[e.id] = true
                   isBridge[e.id^1] = true
               }
           } else {
               if tin[y] < low[x] {
                   low[x] = tin[y]
               }
           }
       }
   }
   for i := 1; i <= n; i++ {
       if tin[i] == 0 {
           dfsBridge(i, -1)
       }
   }
   // assign components
   comp := make([]int, n+1)
   compCnt := 0
   var dfsComp func(int)
   dfsComp = func(x int) {
       comp[x] = compCnt
       for _, e := range adj[x] {
           if comp[e.to] == 0 && !isBridge[e.id] {
               dfsComp(e.to)
           }
       }
   }
   for i := 1; i <= n; i++ {
       if comp[i] == 0 {
           compCnt++
           dfsComp(i)
       }
   }
   // build tree of components
   tree := make([][]int, compCnt+1)
   for i := 0; i < m; i++ {
       a := comp[u[i]]
       b := comp[v[i]]
       if a != b {
           tree[a] = append(tree[a], b)
           tree[b] = append(tree[b], a)
       }
   }
   // tree dp for diameter-like
   dp := make([]int, compCnt+1)
   ans := 0
   var dfsDP func(int, int)
   dfsDP = func(x, p int) {
       for _, y := range tree[x] {
           if y == p {
               continue
           }
           dfsDP(y, x)
           // update ans
           if dp[x]+dp[y]+1 > ans {
               ans = dp[x] + dp[y] + 1
           }
           if dp[y]+1 > dp[x] {
               dp[x] = dp[y] + 1
           }
       }
   }
   // run dp from component of node 1 if exists
   root := 1
   if n >= 1 {
       root = comp[1]
   }
   dfsDP(root, 0)
   fmt.Fprintln(writer, ans)
}

type edge struct {
   to int
   id int
}
