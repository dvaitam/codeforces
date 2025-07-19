package main

import (
   "bufio"
   "fmt"
   "os"
)

// Edge represents a directed edge with capacity and reverse edge index
type Edge struct {
   to, cap, rev int
}

func min(a, b int) int {
   if a < b {
       return a
   }
   return b
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n1, n2, m, q int
   fmt.Fscan(reader, &n1, &n2, &m, &q)
   // node indices: 1..n1 (left), n1+1..n1+n2 (right), st=n1+n2+1, ed=st+1
   st := n1 + n2 + 1
   ed := st + 1
   tot := ed
   // build graph
   graph := make([][]Edge, tot+1)
   level := make([]int, tot+1)
   iter := make([]int, tot+1)

   // addEdge adds forward and reverse edges
   addEdge := func(from, to, cap int) {
       graph[from] = append(graph[from], Edge{to, cap, len(graph[to])})
       graph[to] = append(graph[to], Edge{from, 0, len(graph[from]) - 1})
   }

   ex := make([]int, m+1)
   ey := make([]int, m+1)
   idNode := make([]int, m+1)
   idIdx := make([]int, m+1)
   // read edges and add to flow network
   for i := 1; i <= m; i++ {
       fmt.Fscan(reader, &ex[i], &ey[i])
       u := ex[i]
       v := ey[i] + n1
       // before adding, record length of graph[v]
       // add forward edge u->v cap 1 and reverse
       graph[u] = append(graph[u], Edge{v, 1, len(graph[v])})
       graph[v] = append(graph[v], Edge{u, 0, len(graph[u]) - 1})
       // reverse edge is last in graph[v]
       idNode[i] = v
       idIdx[i] = len(graph[v]) - 1
   }
   // connect source and sink
   for i := 1; i <= n1; i++ {
       addEdge(st, i, 1)
   }
   for i := 1; i <= n2; i++ {
       addEdge(i+n1, ed, 1)
   }

   // Dinic: bfs to build levels
   bfs := func() bool {
       for i := 1; i <= tot; i++ {
           level[i] = -1
       }
       queue := make([]int, 0, tot)
       level[st] = 0
       queue = append(queue, st)
       for qi := 0; qi < len(queue); qi++ {
           v := queue[qi]
           for _, e := range graph[v] {
               if e.cap > 0 && level[e.to] < 0 {
                   level[e.to] = level[v] + 1
                   queue = append(queue, e.to)
                   if e.to == ed {
                       return true
                   }
               }
           }
       }
       return level[ed] >= 0
   }
   // dfs to send flow
   var dfs func(int, int) int
   dfs = func(v, f int) int {
       if v == ed {
           return f
       }
       res := 0
       for i := iter[v]; i < len(graph[v]); i++ {
           e := &graph[v][i]
           if e.cap > 0 && level[v] < level[e.to] {
               d := dfs(e.to, min(f-res, e.cap))
               if d > 0 {
                   e.cap -= d
                   graph[e.to][e.rev].cap += d
                   res += d
                   if res == f {
                       iter[v] = i
                       return res
                   }
               }
           }
       }
       iter[v] = len(graph[v])
       return res
   }

   // max flow
   flow := 0
   const INF = 1<<60
   for bfs() {
       for i := 1; i <= tot; i++ {
           iter[i] = 0
       }
       for {
           f := dfs(st, int(INF))
           if f == 0 {
               break
           }
           flow += f
       }
   }
   // build reachable after max flow
   bfs()
   // find minimum vertex cover nodes p
   p := make([]int, 0, flow)
   vmap := make([]int, tot+1)
   for i := 1; i <= n1; i++ {
       if level[i] < 0 {
           p = append(p, i)
       }
   }
   for i := 1; i <= n2; i++ {
       if level[i+n1] >= 0 {
           p = append(p, i+n1)
       }
   }
   ucnt := len(p)
   // map node to position
   for idx, node := range p {
       vmap[node] = idx + 1
   }
   // assign edge ids in order of cover
   eid := make([]int, ucnt+1)
   for i := 1; i <= m; i++ {
       // reverse edge cap > 0 means matched
       if graph[idNode[i]][idIdx[i]].cap > 0 {
           if vmap[ex[i]] > 0 {
               eid[vmap[ex[i]]] = i
           } else {
               eid[vmap[ey[i]+n1]] = i
           }
       }
   }
   // prefix sums of eid
   sum := make([]int, ucnt+1)
   for i := 1; i <= ucnt; i++ {
       sum[i] = sum[i-1] + eid[i]
   }
   // handle queries
   cur := ucnt
   for qi := 0; qi < q; qi++ {
       var op int
       fmt.Fscan(reader, &op)
       if op == 1 {
           writer.WriteString("1\n")
           node := p[cur-1]
           if node <= n1 {
               writer.WriteString(fmt.Sprintf("%d\n", node))
           } else {
               writer.WriteString(fmt.Sprintf("%d\n", -(node-n1)))
           }
           cur--
           writer.WriteString(fmt.Sprintf("%d\n", sum[cur]))
       } else {
           writer.WriteString(fmt.Sprintf("%d\n", cur))
           for i := 1; i <= cur; i++ {
               writer.WriteString(fmt.Sprintf("%d ", eid[i]))
           }
           writer.WriteByte('\n')
       }
   }
}
