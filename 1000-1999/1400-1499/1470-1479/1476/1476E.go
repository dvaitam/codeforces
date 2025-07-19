package main

import (
   "bufio"
   "container/heap"
   "fmt"
   "os"
)

// IntHeap is a min-heap of ints
type IntHeap []int
func (h IntHeap) Len() int            { return len(h) }
func (h IntHeap) Less(i, j int) bool  { return h[i] < h[j] }
func (h IntHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *IntHeap) Push(x interface{}) { *h = append(*h, x.(int)) }
func (h *IntHeap) Pop() interface{} {
   old := *h
   n := len(old)
   x := old[n-1]
   *h = old[0 : n-1]
   return x
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var N, M, K int
   fmt.Fscan(reader, &N, &M, &K)
   A := make([]string, N)
   for i := 0; i < N; i++ {
       fmt.Fscan(reader, &A[i])
   }
   B := make([]string, M)
   P := make([]int, M)
   // compute sizes
   bMx := 1 << K
   nodes := 1
   for i := 0; i < K; i++ {
       nodes *= 27
   }
   // allocate graph structures
   isValid := make([]bool, nodes)
   id := make([]int, nodes)
   par := make([]int, nodes)
   adj := make([][]int, nodes)

   // helper functions
   getCharId := func(ch byte) int {
       if ch == '_' {
           return 26
       }
       return int(ch - 'a')
   }
   hsh := func(s string) int {
       res := 0
       mult := 1
       for i := 0; i < K; i++ {
           res += mult * getCharId(s[i])
           mult *= 27
       }
       return res
   }
   match := func(s, t string) bool {
       for i := 0; i < K; i++ {
           if s[i] != '_' && t[i] != '_' && t[i] != s[i] {
               return false
           }
       }
       return true
   }
   // mark valid nodes
   for i := 0; i < N; i++ {
       u := hsh(A[i])
       isValid[u] = true
       id[u] = i + 1
   }
   // add edges
   for i := 0; i < M; i++ {
       fmt.Fscan(reader, &B[i], &P[i])
       p := P[i] - 1
       // check match
       if !match(B[i], A[p]) {
           fmt.Fprintln(writer, "NO")
           return
       }
       sBytes := []byte(B[i])
       t := A[p]
       u := hsh(t)
       // generate masks
       for mask := 0; mask < bMx; mask++ {
           tmp := make([]byte, K)
           copy(tmp, sBytes)
           for j := 0; j < K; j++ {
               if mask&(1<<j) != 0 {
                   tmp[j] = '_'
               }
           }
           v := hsh(string(tmp))
           if u == v || !isValid[u] || !isValid[v] {
               continue
           }
           adj[u] = append(adj[u], v)
           par[v]++
       }
   }
   // topo sort
   h := &IntHeap{}
   heap.Init(h)
   for i := 0; i < N; i++ {
       u := hsh(A[i])
       if par[u] == 0 {
           heap.Push(h, u)
       }
   }
   result := make([]int, 0, N)
   for len(result) < N {
       if h.Len() == 0 {
           fmt.Fprintln(writer, "NO")
           os.Exit(0)
       }
       u := heap.Pop(h).(int)
       result = append(result, id[u])
       for _, v := range adj[u] {
           par[v]--
           if par[v] == 0 {
               heap.Push(h, v)
           }
       }
   }
   fmt.Fprintln(writer, "YES")
   for i, v := range result {
       if i > 0 {
           writer.WriteByte(' ')
       }
       fmt.Fprint(writer, v)
   }
   fmt.Fprintln(writer)
}
