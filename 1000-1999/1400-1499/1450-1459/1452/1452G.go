package main

import (
   "bufio"
   "container/heap"
   "fmt"
   "os"
)

// State for priority queue: threshold value and remaining reach
type state struct {
   th  int // original threshold (d_Bob at source)
   rem int // remaining distance to propagate
   u   int // current node
}

// Max-heap by threshold
type pqState []state

func (pq pqState) Len() int { return len(pq) }
func (pq pqState) Less(i, j int) bool {
   return pq[i].th > pq[j].th
}
func (pq pqState) Swap(i, j int) { pq[i], pq[j] = pq[j], pq[i] }
func (pq *pqState) Push(x interface{}) {
   *pq = append(*pq, x.(state))
}
func (pq *pqState) Pop() interface{} {
   old := *pq
   n := len(old)
   x := old[n-1]
   *pq = old[:n-1]
   return x
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   var n int
   fmt.Fscan(in, &n)
   adj := make([][]int, n+1)
   for i := 0; i < n-1; i++ {
       var u, v int
       fmt.Fscan(in, &u, &v)
       adj[u] = append(adj[u], v)
       adj[v] = append(adj[v], u)
   }
   var k int
   fmt.Fscan(in, &k)
   bob := make([]int, k)
   for i := 0; i < k; i++ {
       fmt.Fscan(in, &bob[i])
   }
   // multi-source BFS for d_Bob
   const INF = 1e9
   dBob := make([]int, n+1)
   for i := 1; i <= n; i++ {
       dBob[i] = INF
   }
   q := make([]int, 0, n)
   for _, v := range bob {
       dBob[v] = 0
       q = append(q, v)
   }
   for qi := 0; qi < len(q); qi++ {
       u := q[qi]
       for _, v := range adj[u] {
           if dBob[v] > dBob[u]+1 {
               dBob[v] = dBob[u] + 1
               q = append(q, v)
           }
       }
   }
   // best threshold covering each node
   bestTh := make([]int, n+1)
   for i := 1; i <= n; i++ {
       bestTh[i] = -1
   }
   // priority queue of states
   pqh := &pqState{}
   heap.Init(pqh)
   // initialize states from all nodes as potential capture centers
   for u := 1; u <= n; u++ {
       th := dBob[u]
       // push initial state
       heap.Push(pqh, state{th: th, rem: th, u: u})
   }
   // propagate
   for pqh.Len() > 0 {
       st := heap.Pop(pqh).(state)
       u, th, rem := st.u, st.th, st.rem
       if th <= bestTh[u] {
           continue
       }
       bestTh[u] = th
       if rem == 0 {
           continue
       }
       for _, v := range adj[u] {
           if th > bestTh[v] {
               // can propagate
               heap.Push(pqh, state{th: th, rem: rem - 1, u: v})
           }
       }
   }
   // output
   for i := 1; i <= n; i++ {
       if i > 1 {
           out.WriteByte(' ')
       }
       fmt.Fprint(out, bestTh[i])
   }
   out.WriteByte('\n')
}
