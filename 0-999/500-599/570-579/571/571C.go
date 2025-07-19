package main

import (
   "bufio"
   "fmt"
   "os"
)

// Edge represents an edge in the graph
type Edge struct {
   to, nxt, from int
}

var (
   head  []int
   edges []Edge
   le    int
   du, rdu []int
   cate  []int
   p1, p2 []int
   have  []bool
   anse  []int
   used  []bool
   ans   []byte
   nc, n int
)

func dfs(u, fe int) bool {
   flag := false
   have[u] = true
   if rdu[u] != du[u] {
       flag = true
   }
   for i := head[u]; i != -1; i = edges[i].nxt {
       if i == fe {
           continue
       }
       id := i >> 1
       if used[id] {
           continue
       }
       v := edges[i].to
       if have[v] {
           anse[u] = i
           flag = true
           used[id] = true
       } else {
           if dfs(v, i^1) {
               anse[u] = i
               flag = true
           } else {
               anse[v] = i ^ 1
           }
           used[id] = true
       }
   }
   return flag
}

func addEdge(u, v, from int) {
   edges[le] = Edge{to: v, nxt: head[u], from: from}
   head[u] = le
   le++
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   in := bufio.NewReader(reader)
   var err error
   // read nc, n
   if _, err = fmt.Fscan(in, &nc, &n); err != nil {
       return
   }
   // init
   head = make([]int, nc)
   for i := range head {
       head[i] = -1
   }
   du = make([]int, nc)
   rdu = make([]int, nc)
   cate = make([]int, n+1)
   p1 = make([]int, n+1)
   p2 = make([]int, n+1)
   have = make([]bool, nc)
   anse = make([]int, nc)
   for i := range anse {
       anse[i] = -1
   }
   ans = make([]byte, n+1)
   // prepare edges slice (max possible 2*n edges)
   edges = make([]Edge, 2*n+5)
   used = make([]bool, n+1)
   // read categories
   for i := 0; i < nc; i++ {
       var k int
       fmt.Fscan(in, &k)
       rdu[i] = k
       for j := 0; j < k; j++ {
           var a int
           fmt.Fscan(in, &a)
           if a < 0 {
               idx := -a
               p2[idx] = i
               cate[idx] |= 2
           } else {
               p1[a] = i
               cate[a] |= 1
           }
       }
   }
   // build graph and default answers
   for i := 1; i <= n; i++ {
       if cate[i] != 3 || p1[i] == p2[i] {
           if cate[i]&2 != 0 {
               ans[i] = '0'
           } else {
               ans[i] = '1'
           }
       } else {
           // two-category item
           addEdge(p1[i], p2[i], i)
           addEdge(p2[i], p1[i], -i)
           du[p1[i]]++
           du[p2[i]]++
       }
   }
   // DFS to orient edges
   for i := 0; i < nc; i++ {
       if !have[i] {
           if !dfs(i, -1) {
               fmt.Println("NO")
               return
           }
       }
   }
   // assign answers from oriented edges
   for u := 0; u < nc; u++ {
       ei := anse[u]
       if ei != -1 {
           if ei&1 == 1 {
               idx := -edges[ei].from
               ans[idx] = '0'
           } else {
               idx := edges[ei].from
               ans[idx] = '1'
           }
       }
   }
   // fill defaults
   for i := 1; i <= n; i++ {
       if ans[i] == 0 {
           ans[i] = '1'
       }
   }
   // output
   fmt.Println("YES")
   fmt.Println(string(ans[1:]))
}
