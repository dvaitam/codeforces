package main

import (
   "bufio"
   "container/heap"
   "fmt"
   "os"
   "sort"
)

var (
   reader = bufio.NewReader(os.Stdin)
   writer = bufio.NewWriter(os.Stdout)
)

func readInt() int {
   var x int
   fmt.Fscan(reader, &x)
   return x
}

type Edge struct{ idx, t int }
type Pair struct{ val, idx int }
type MaxHeap []Pair

func (h MaxHeap) Len() int            { return len(h) }
func (h MaxHeap) Less(i, j int) bool  { return h[i].val > h[j].val }
func (h MaxHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *MaxHeap) Push(x interface{}) { *h = append(*h, x.(Pair)) }
func (h *MaxHeap) Pop() interface{} {
   old := *h
   n := len(old)
   item := old[n-1]
   *h = old[:n-1]
   return item
}

func solve() {
   n := readInt()
   m := readInt()
   a := make([]int, n)
   for i := 0; i < n; i++ {
       a[i] = readInt()
   }
   adj := make([][]Edge, n)
   ver := make([][3]int, m)
   f := make([][3]int, m)
   for i := 0; i < m; i++ {
       u := readInt() - 1
       v := readInt() - 1
       w := readInt() - 1
       ver[i] = [3]int{u, v, w}
       adj[u] = append(adj[u], Edge{i, 0})
       adj[v] = append(adj[v], Edge{i, 1})
       adj[w] = append(adj[w], Edge{i, 2})
       b := []int{0, 1, 2}
       sort.Slice(b, func(x, y int) bool {
           return ver[i][b[x]] < ver[i][b[y]]
       })
       for j := 0; j < 3; j++ {
           pos := b[j]
           a[ver[i][pos]] += j + 3
           f[i][pos] += j + 3
       }
   }
   h := &MaxHeap{}
   heap.Init(h)
   for i := 0; i < n; i++ {
       heap.Push(h, Pair{a[i], i})
   }
   for h.Len() > 0 {
       p := heap.Pop(h).(Pair)
       x, v := p.idx, p.val
       if a[x] != v {
           continue
       }
       for _, e := range adj[x] {
           j, t := e.idx, e.t
           y := ver[j][(t+1)%3]
           z := ver[j][(t+2)%3]
           if a[x] == a[y] && a[x] == a[z] {
               a[y]++
               f[j][(t+1)%3]++
               heap.Push(h, Pair{a[y], y})
               a[z]++
               f[j][(t+2)%3]++
               heap.Push(h, Pair{a[z], z})
           } else if a[x] == a[y] {
               a[y] += 2
               f[j][(t+1)%3] += 2
               heap.Push(h, Pair{a[y], y})
           } else if a[x] == a[z] {
               a[z] += 2
               f[j][(t+2)%3] += 2
               heap.Push(h, Pair{a[z], z})
           }
       }
   }
   for i := 0; i < m; i++ {
       u, v, w := ver[i][0], ver[i][1], ver[i][2]
       aa := (f[i][0] + f[i][1] - f[i][2]) / 2
       bb := (f[i][1] + f[i][2] - f[i][0]) / 2
       cc := (f[i][2] + f[i][0] - f[i][1]) / 2
       fmt.Fprintln(writer, aa, bb, cc)
   }
}

func main() {
   defer writer.Flush()
   t := readInt()
   for ; t > 0; t-- {
       solve()
   }
}
