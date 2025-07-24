package main

import (
   "bufio"
   "container/heap"
   "fmt"
   "os"
)

// Event represents a group of servers that will be released at releaseTime
type Event struct {
   releaseTime int
   servers     []int
}

// EventPQ implements heap.Interface based on releaseTime
type EventPQ []*Event

func (pq EventPQ) Len() int { return len(pq) }
func (pq EventPQ) Less(i, j int) bool {
   return pq[i].releaseTime < pq[j].releaseTime
}
func (pq EventPQ) Swap(i, j int) { pq[i], pq[j] = pq[j], pq[i] }
func (pq *EventPQ) Push(x interface{}) {
   *pq = append(*pq, x.(*Event))
}
func (pq *EventPQ) Pop() interface{} {
   old := *pq
   n := len(old)
   x := old[n-1]
   *pq = old[0 : n-1]
   return x
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, q int
   if _, err := fmt.Fscan(reader, &n, &q); err != nil {
       return
   }
   // free[i] == true means server i is available
   free := make([]bool, n+1)
   for i := 1; i <= n; i++ {
       free[i] = true
   }
   // priority queue of release events
   var pq EventPQ
   heap.Init(&pq)

   for t := 0; t < q; t++ {
       var ti, ki, di int
       fmt.Fscan(reader, &ti, &ki, &di)
       // release servers whose time has come
       for pq.Len() > 0 && pq[0].releaseTime <= ti {
           evt := heap.Pop(&pq).(*Event)
           for _, sid := range evt.servers {
               free[sid] = true
           }
       }
       // pick ki smallest free servers
       picked := make([]int, 0, ki)
       for id := 1; id <= n && len(picked) < ki; id++ {
           if free[id] {
               picked = append(picked, id)
           }
       }
       if len(picked) < ki {
           fmt.Fprintln(writer, -1)
           continue
       }
       // allocate servers
       sum := 0
       for _, id := range picked {
           free[id] = false
           sum += id
       }
       // create release event
       evt := &Event{releaseTime: ti + di, servers: picked}
       heap.Push(&pq, evt)
       fmt.Fprintln(writer, sum)
   }
}
