package main

import (
   "bufio"
   "container/heap"
   "fmt"
   "os"
   "sort"
)

// Item is an entry in the priority queue
type Item struct {
   id, d int
}

// PriorityQueue implements a min-heap of Items by d, then id
type PriorityQueue []Item

func (pq PriorityQueue) Len() int { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool {
   if pq[i].d != pq[j].d {
       return pq[i].d < pq[j].d
   }
   return pq[i].id < pq[j].id
}
func (pq PriorityQueue) Swap(i, j int) { pq[i], pq[j] = pq[j], pq[i] }
func (pq *PriorityQueue) Push(x interface{}) {
   *pq = append(*pq, x.(Item))
}
func (pq *PriorityQueue) Pop() interface{} {
   old := *pq
   n := len(old)
   it := old[n-1]
   *pq = old[0 : n-1]
   return it
}

// popValid pops until it finds a valid Item matching current d and >0, or returns false
func popValid(pq *PriorityQueue, d []int) (Item, bool) {
   for pq.Len() > 0 {
       it := heap.Pop(pq).(Item)
       if it.d != d[it.id] {
           continue
       }
       if it.d <= 0 {
           continue
       }
       return it, true
   }
   return Item{}, false
}

// EdgeKey for visited edges
type EdgeKey struct{ u, v int }

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n, m int
   fmt.Fscan(in, &n, &m)
   d := make([]int, n+1)
   adj := make([][]int, n+1)
   for i := 0; i < m; i++ {
       var x, y int
       fmt.Fscan(in, &x, &y)
       d[x]++
       d[y]++
       adj[x] = append(adj[x], y)
       adj[y] = append(adj[y], x)
   }
   // initial ceil(deg/2)
   for i := 1; i <= n; i++ {
       d[i] = d[i]/2 + d[i]%2
   }
   // build initial heap
   pq := &PriorityQueue{}
   heap.Init(pq)
   for i := 1; i <= n; i++ {
       heap.Push(pq, Item{id: i, d: d[i]})
   }

   visited := make(map[EdgeKey]struct{})
   type Nd struct{ d, id int }
   var temp []Nd
   var res [][2]int

   for i := 1; i <= n; i++ {
       it, ok := popValid(pq, d)
       if !ok {
           break
       }
       x := it.id
       temp = temp[:0]
       for _, y := range adj[x] {
           // skip if reversed edge visited
           if _, seen := visited[EdgeKey{u: y, v: x}]; seen {
               continue
           }
           temp = append(temp, Nd{d: d[y], id: y})
       }
       sort.Slice(temp, func(i, j int) bool {
           if temp[i].d != temp[j].d {
               return temp[i].d > temp[j].d
           }
           return temp[i].id > temp[j].id
       })
       for _, nd := range temp {
           if d[x] <= 0 {
               break
           }
           y := nd.id
           res = append(res, [2]int{x, y})
           visited[EdgeKey{u: x, v: y}] = struct{}{}
           d[x]--
           d[y]--
           if d[y] >= 0 {
               heap.Push(pq, Item{id: y, d: d[y]})
           }
       }
   }
   // output
   fmt.Fprintln(out, len(res))
   for _, e := range res {
       fmt.Fprintln(out, e[0], e[1])
   }
}
