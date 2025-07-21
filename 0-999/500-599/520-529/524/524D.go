package main

import (
   "bufio"
   "container/heap"
   "fmt"
   "os"
)

// Node represents an active user with its end time and identifier
type Node struct {
   end int
   id  int
}

// MinHeap implements a min-heap based on end time
type MinHeap []Node
func (h MinHeap) Len() int           { return len(h) }
func (h MinHeap) Less(i, j int) bool { return h[i].end < h[j].end }
func (h MinHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *MinHeap) Push(x interface{}) {
   *h = append(*h, x.(Node))
}
func (h *MinHeap) Pop() interface{} {
   old := *h
   n := len(old)
   x := old[n-1]
   *h = old[:n-1]
   return x
}

// MaxHeap implements a max-heap based on end time
type MaxHeap []Node
func (h MaxHeap) Len() int           { return len(h) }
func (h MaxHeap) Less(i, j int) bool { return h[i].end > h[j].end }
func (h MaxHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *MaxHeap) Push(x interface{}) {
   *h = append(*h, x.(Node))
}
func (h *MaxHeap) Pop() interface{} {
   old := *h
   n := len(old)
   x := old[n-1]
   *h = old[:n-1]
   return x
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, M, T int
   if _, err := fmt.Fscanf(in, "%d %d %d\n", &n, &M, &T); err != nil {
       return
   }
   times := make([]int, n)
   for i := 0; i < n; i++ {
       var s string
       fmt.Fscanf(in, "%s\n", &s)
       h := int(s[0]-'0')*10 + int(s[1]-'0')
       m := int(s[3]-'0')*10 + int(s[4]-'0')
       sec := int(s[6]-'0')*10 + int(s[7]-'0')
       times[i] = h*3600 + m*60 + sec
   }
   // currentEnd stores the current end time for each user id
   currentEnd := make([]int, n+2)
   ans := make([]int, n)
   activeMin := &MinHeap{}
   activeMax := &MaxHeap{}
   heap.Init(activeMin)
   heap.Init(activeMax)
   totalUsers := 0
   activeCount := 0
   maxActive := 0
   for i, s := range times {
       // expire users whose interval ended before s
       for activeMin.Len() > 0 {
           top := (*activeMin)[0]
           if top.end >= s {
               break
           }
           heap.Pop(activeMin)
           // skip stale entries
           if currentEnd[top.id] != top.end {
               continue
           }
           // expire this user
           activeCount--
           currentEnd[top.id] = 0
       }
       if activeCount < M {
           // assign new user
           totalUsers++
           uid := totalUsers
           ans[i] = uid
           end := s + T - 1
           currentEnd[uid] = end
           heap.Push(activeMin, Node{end, uid})
           heap.Push(activeMax, Node{end, uid})
           activeCount++
           if activeCount > maxActive {
               maxActive = activeCount
           }
       } else {
           // reuse an active user with the latest end time
           var id int
           var end0 int
           for {
               top := heap.Pop(activeMax).(Node)
               if currentEnd[top.id] == top.end {
                   id = top.id
                   end0 = top.end
                   break
               }
           }
           ans[i] = id
           newEnd := end0
           if s+T-1 > end0 {
               newEnd = s + T - 1
               currentEnd[id] = newEnd
               heap.Push(activeMin, Node{newEnd, id})
               heap.Push(activeMax, Node{newEnd, id})
           }
       }
   }
   if maxActive < M {
       fmt.Println("No solution")
       return
   }
   // output result
   fmt.Println(totalUsers)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   for _, id := range ans {
       fmt.Fprintln(writer, id)
   }
}
