package main

import (
   "bufio"
   "container/heap"
   "fmt"
   "os"
)

// Segment represents an empty interval between occupied positions l and r
type Segment struct {
   l, r int
}

// PriorityQueue implements heap.Interface for segments with custom priority
type PriorityQueue struct {
   data      []Segment
   active    map[int64]bool
   n         int
}

func (pq PriorityQueue) Len() int { return len(pq.data) }
// We want max-heap by distance, then min position
func (pq PriorityQueue) Less(i, j int) bool {
   si, sj := pq.data[i], pq.data[j]
   // compute position and distance for si
   di, pi := pq.eval(si)
   dj, pj := pq.eval(sj)
   if di != dj {
       return di > dj
   }
   return pi < pj
}
func (pq PriorityQueue) Swap(i, j int) { pq.data[i], pq.data[j] = pq.data[j], pq.data[i] }
func (pq *PriorityQueue) Push(x interface{}) { pq.data = append(pq.data, x.(Segment)) }
func (pq *PriorityQueue) Pop() interface{} {
   old := pq.data
   n := len(old)
   x := old[n-1]
   pq.data = old[0 : n-1]
   return x
}

// eval returns distance and chosen position for segment s
func (pq *PriorityQueue) eval(s Segment) (dist, pos int) {
   l, r := s.l, s.r
   if l == 0 {
       pos = 1
       dist = r - 1
   } else if r == pq.n+1 {
       pos = pq.n
       dist = pq.n - l
   } else {
       pos = (l + r) / 2
       dist = pos - l
   }
   return
}

// key generates a unique key for segment
func key(l, r int) int64 {
   return (int64(l) << 32) | int64(r)
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   var n, m int
   fmt.Fscan(in, &n, &m)
   // neighbors in occupied linked list, include 0 and n+1
   lneigh := make([]int, n+2)
   rneigh := make([]int, n+2)
   for i := 0; i <= n+1; i++ {
       lneigh[i] = i - 1
       rneigh[i] = i + 1
   }
   rneigh[n+1] = -1
   lneigh[0] = -1
   // priority queue of segments
   pq := &PriorityQueue{data: make([]Segment, 0, m+1), active: make(map[int64]bool), n: n}
   heap.Init(pq)
   // initial segment [0, n+1]
   pq.active[key(0, n+1)] = true
   heap.Push(pq, Segment{0, n + 1})
   // map car id to position
   carPos := make(map[int]int, m)
   for i := 0; i < m; i++ {
       var t, id int
       fmt.Fscan(in, &t, &id)
       if t == 1 {
           // arrival: find best segment
           var s Segment
           for {
               if pq.Len() == 0 {
                   break
               }
               s = heap.Pop(pq).(Segment)
               k := key(s.l, s.r)
               if pq.active[k] {
                   delete(pq.active, k)
                   break
               }
           }
           // choose position
           _, pos := pq.eval(s)
           fmt.Fprintln(out, pos)
           carPos[id] = pos
           l, r := s.l, s.r
           // update neighbors linked list: insert pos between l and r
           rneigh[l] = pos
           lneigh[r] = pos
           lneigh[pos] = l
           rneigh[pos] = r
           // create new segments [l,pos] and [pos,r]
           if pos-l > 1 {
               a, b := l, pos
               pq.active[key(a, b)] = true
               heap.Push(pq, Segment{a, b})
           }
           if r-pos > 1 {
               a, b := pos, r
               pq.active[key(a, b)] = true
               heap.Push(pq, Segment{a, b})
           }
       } else {
           // departure
           pos := carPos[id]
           delete(carPos, id)
           l, r := lneigh[pos], rneigh[pos]
           // remove segments [l,pos] and [pos,r]
           if pos-l > 1 {
               delete(pq.active, key(l, pos))
           }
           if r-pos > 1 {
               delete(pq.active, key(pos, r))
           }
           // link l and r
           rneigh[l] = r
           lneigh[r] = l
           // add merged segment
           if r-l > 1 {
               pq.active[key(l, r)] = true
               heap.Push(pq, Segment{l, r})
           }
       }
   }
}
