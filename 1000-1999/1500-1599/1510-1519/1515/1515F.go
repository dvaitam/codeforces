package main

import (
   "bufio"
   "fmt"
   "os"
)

type Edge struct {
   to, id int
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   // fast int reader
   var readInt func() (int, error)
   readInt = func() (int, error) {
       var x int
       var sign = 1
       ch, err := reader.ReadByte()
       if err != nil {
           return 0, err
       }
       for (ch < '0' || ch > '9') && ch != '-' {
           ch, err = reader.ReadByte()
           if err != nil {
               return 0, err
           }
       }
       if ch == '-' {
           sign = -1
           ch, err = reader.ReadByte()
           if err != nil {
               return 0, err
           }
       }
       for ch >= '0' && ch <= '9' {
           x = x*10 + int(ch-'0')
           ch, err = reader.ReadByte()
           if err != nil {
               break
           }
       }
       return x * sign, nil
   }
   // read N, M, X
   N, _ := readInt()
   M, _ := readInt()
   X, _ := readInt()
   A := make([]int64, N)
   var total int64
   for i := 0; i < N; i++ {
       ai, _ := readInt()
       A[i] = int64(ai)
       total += A[i]
   }
   U := make([]int, M)
   V := make([]int, M)
   for i := 0; i < M; i++ {
       ui, _ := readInt()
       vi, _ := readInt()
       U[i] = ui - 1
       V[i] = vi - 1
   }
   if total < int64(X)*int64(N-1) {
       fmt.Fprintln(writer, "NO")
       return
   }
   fmt.Fprintln(writer, "YES")
   // DSU
   parent := make([]int, N)
   for i := range parent {
       parent[i] = -1
   }
   var find func(int) int
   find = func(x int) int {
       if parent[x] < 0 {
           return x
       }
       parent[x] = find(parent[x])
       return parent[x]
   }
   union := func(a, b int) bool {
       a = find(a)
       b = find(b)
       if a == b {
           return false
       }
       if parent[a] < parent[b] {
           parent[a] += parent[b]
           parent[b] = a
       } else {
           parent[b] += parent[a]
           parent[a] = b
       }
       return true
   }
   // select edges for spanning tree
   type Pair struct{u, v, id int}
   edges := make([]Pair, 0, N-1)
   for i := 0; i < M; i++ {
       if union(U[i], V[i]) {
           edges = append(edges, Pair{U[i], V[i], i + 1})
       }
   }
   // build adjacency
   adj := make([][]Edge, N)
   for _, e := range edges {
       adj[e.u] = append(adj[e.u], Edge{e.v, e.id})
       adj[e.v] = append(adj[e.v], Edge{e.u, e.id})
   }
   // result order
   res := make([]int, N-1)
   r1, r2 := 0, 0
   // iterative post-order DFS
   type St struct{node, parent, peid int; post bool}
   stack := make([]St, 0, 2*N)
   stack = append(stack, St{0, -1, -1, false})
   for len(stack) > 0 {
       st := stack[len(stack)-1]
       stack = stack[:len(stack)-1]
       if st.post {
           if st.peid >= 0 {
               if A[st.node] >= int64(X) {
                   A[st.parent] += A[st.node] - int64(X)
                   res[r1] = st.peid
                   r1++
               } else {
                   res[len(res)-1-r2] = st.peid
                   r2++
               }
           }
       } else {
           // pre
           stack = append(stack, St{st.node, st.parent, st.peid, true})
           // children
           for i := len(adj[st.node]) - 1; i >= 0; i-- {
               e := adj[st.node][i]
               if e.to == st.parent {
                   continue
               }
               stack = append(stack, St{e.to, st.node, e.id, false})
           }
       }
   }
   // output
   for i := 0; i < len(res); i++ {
       fmt.Fprintln(writer, res[i])
   }
}
