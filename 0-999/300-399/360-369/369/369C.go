package main

import (
   "bufio"
   "fmt"
   "os"
)

type edge struct { to, flag int }

var (
   r = bufio.NewReader(os.Stdin)
   w = bufio.NewWriter(os.Stdout)
)

// readInt reads next integer from stdin.
func readInt() (int, error) {
   var x int
   // skip non-digits
   for {
       b, err := r.ReadByte()
       if err != nil {
           return 0, err
       }
       if b >= '0' && b <= '9' {
           x = int(b - '0')
           break
       }
   }
   // read remaining digits
   for {
       b, err := r.ReadByte()
       if err != nil || b < '0' || b > '9' {
           break
       }
       x = x*10 + int(b-'0')
   }
   return x, nil
}

func main() {
   defer w.Flush()
   for {
       N, err := readInt()
       if err != nil {
           break
       }
       // build undirected graph
       graph := make([][]edge, N+1)
       for i := 1; i < N; i++ {
           a, _ := readInt()
           b, _ := readInt()
           c, _ := readInt()
           graph[a] = append(graph[a], edge{b, c})
           graph[b] = append(graph[b], edge{a, c})
       }
       // build rooted tree at 1
       children := make([][]edge, N+1)
       indeg := make([]int, N+1)
       parent := make([]int, N+1)
       // BFS/DFS to orient edges
       stack := make([]int, 0, N)
       stack = append(stack, 1)
       parent[1] = 0
       for i := 0; i < len(stack); i++ {
           u := stack[i]
           for _, e1 := range graph[u] {
               v := e1.to
               if v == parent[u] {
                   continue
               }
               parent[v] = u
               children[u] = append(children[u], edge{v, e1.flag})
               indeg[v]++
               stack = append(stack, v)
           }
       }
       // initialize queue with leaves
       q := make([]int, 0, N)
       for i := 1; i <= N; i++ {
           if indeg[i] == 0 {
               q = append(q, i)
           }
       }
       vis := make([]bool, N+1)
       ans := make([]int, 0)
       // process in topological order
       for head := 0; head < len(q); head++ {
           u := q[head]
           for _, e1 := range children[u] {
               v, f := e1.to, e1.flag
               indeg[v]--
               if indeg[v] == 0 {
                   q = append(q, v)
               }
               if f == 1 || vis[u] {
                   continue
               }
               ans = append(ans, u)
               // mark u and all descendants as visited
               st2 := []int{u}
               for len(st2) > 0 {
                   x := st2[len(st2)-1]
                   st2 = st2[:len(st2)-1]
                   if vis[x] {
                       continue
                   }
                   vis[x] = true
                   for _, e2 := range children[x] {
                       st2 = append(st2, e2.to)
                   }
               }
           }
       }
       // output result
       fmt.Fprintln(w, len(ans))
       for i, v := range ans {
           if i > 0 {
               w.WriteByte(' ')
           }
           fmt.Fprint(w, v)
       }
       w.WriteByte('\n')
   }
}
