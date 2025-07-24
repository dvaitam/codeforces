package main

import (
   "bufio"
   "container/heap"
   "fmt"
   "os"
)

// entry represents a subject in the priority queue
type entry struct {
   deadline int
   id       int
}

// minHeap implements heap.Interface for entries sorted by deadline
type minHeap []entry

func (h minHeap) Len() int { return len(h) }
func (h minHeap) Less(i, j int) bool { return h[i].deadline < h[j].deadline }
func (h minHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }
func (h *minHeap) Push(x interface{}) { *h = append(*h, x.(entry)) }
func (h *minHeap) Pop() interface{} {
   old := *h
   n := len(old)
   x := old[n-1]
   *h = old[:n-1]
   return x
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m int
   fmt.Fscan(reader, &n, &m)
   d := make([]int, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &d[i])
   }
   a := make([]int, m+1)
   for i := 1; i <= m; i++ {
       fmt.Fscan(reader, &a[i])
   }
   // binary search minimal days
   lo, hi, ans := 1, n, -1
   for lo <= hi {
       mid := (lo + hi) / 2
       if can(mid, n, m, d, a) {
           ans = mid
           hi = mid - 1
       } else {
           lo = mid + 1
       }
   }
   fmt.Println(ans)
}

// can determines if it's possible to pass all exams within first mid days
func can(mid, n, m int, d, a []int) bool {
   // find last exam day for each subject
   last := make([]int, m+1)
   for i := 1; i <= mid; i++ {
       subj := d[i]
       if subj > 0 {
           last[subj] = i
       }
   }
   // scheduledExam[t] = subject to take exam on day t, or 0
   sched := make([]int, mid+1)
   for subj := 1; subj <= m; subj++ {
       if last[subj] == 0 {
           return false
       }
       sched[last[subj]] = subj
   }
   // remaining preparation days
   rem := make([]int, m+1)
   for i := 1; i <= m; i++ {
       rem[i] = a[i]
   }
   // build priority queue of subjects by deadline
   h := &minHeap{}
   heap.Init(h)
   for subj := 1; subj <= m; subj++ {
       heap.Push(h, entry{deadline: last[subj], id: subj})
   }
   // simulate days
   for day := 1; day <= mid; day++ {
       if sched[day] > 0 {
           subj := sched[day]
           if rem[subj] > 0 {
               return false
           }
           // exam taken, mark rem zero
           rem[subj] = 0
       } else {
           // free day, prepare for earliest deadline subject
           for h.Len() > 0 {
               top := (*h)[0]
               if rem[top.id] == 0 {
                   heap.Pop(h)
                   continue
               }
               break
           }
           if h.Len() == 0 {
               continue
           }
           top := (*h)[0]
           if top.deadline <= day {
               return false
           }
           // use this day to prepare
           rem[top.id]--
       }
   }
   return true
}
