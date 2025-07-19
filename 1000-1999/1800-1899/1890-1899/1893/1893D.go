package main

import (
   "bufio"
   "fmt"
   "os"
)

// Node represents an element with count and value
type Node struct {
   cnt int
   val int
}

// MaxHeap implements a max-heap of Nodes
type MaxHeap struct { data []Node }

// comparator: return true if a > b
func (h *MaxHeap) greater(i, j int) bool {
   a, b := h.data[i], h.data[j]
   if a.cnt == b.cnt {
       return a.val > b.val
   }
   return a.cnt > b.cnt
}

// Push adds a node to the heap
func (h *MaxHeap) Push(node Node) {
   h.data = append(h.data, node)
   i := len(h.data) - 1
   for i > 0 {
       p := (i - 1) / 2
       if h.greater(i, p) {
           h.data[i], h.data[p] = h.data[p], h.data[i]
           i = p
       } else {
           break
       }
   }
}

// Pop removes and returns the max node
func (h *MaxHeap) Pop() Node {
   n := len(h.data)
   top := h.data[0]
   if n == 1 {
       h.data = h.data[:0]
       return top
   }
   h.data[0] = h.data[n-1]
   h.data = h.data[:n-1]
   i := 0
   for {
       left := 2*i + 1
       right := 2*i + 2
       largest := i
       if left < len(h.data) && h.greater(left, largest) {
           largest = left
       }
       if right < len(h.data) && h.greater(right, largest) {
           largest = right
       }
       if largest == i {
           break
       }
       h.data[i], h.data[largest] = h.data[largest], h.data[i]
       i = largest
   }
   return top
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var t int
   if _, err := fmt.Fscan(reader, &t); err != nil {
       return
   }
   for t > 0 {
       t--
       var n, m int
       fmt.Fscan(reader, &n, &m)
       a := make([]int, n)
       for i := 0; i < n; i++ {
           fmt.Fscan(reader, &a[i])
       }
       s := make([]int, m)
       d := make([]int, m)
       for i := 0; i < m; i++ {
           fmt.Fscan(reader, &s[i])
       }
       for i := 0; i < m; i++ {
           fmt.Fscan(reader, &d[i])
       }
       cnt := make([]int, n+1)
       for _, x := range a {
           if x >= 1 && x <= n {
               cnt[x]++
           }
       }
       // initialize heap
       h := &MaxHeap{}
       for x := 1; x <= n; x++ {
           if cnt[x] > 0 {
               h.Push(Node{cnt: cnt[x], val: x})
           }
       }
       ans := make([][]int, m)
       ok := true
       for i := 0; i < m; i++ {
           ans[i] = make([]int, s[i])
           // fill s[i] elements
           for j := 0; j < s[i]; j++ {
               if j >= d[i] {
                   prev := ans[i][j-d[i]]
                   if cnt[prev] > 0 {
                       h.Push(Node{cnt: cnt[prev], val: prev})
                   }
               }
               if len(h.data) == 0 {
                   fmt.Fprintln(writer, -1)
                   ok = false
                   break
               }
               node := h.Pop()
               ans[i][j] = node.val
               cnt[node.val]--
           }
           if !ok {
               break
           }
           // reinsert remaining delayed elements
           for j := s[i]; j < s[i]+d[i]; j++ {
               prev := ans[i][j-d[i]]
               if cnt[prev] > 0 {
                   h.Push(Node{cnt: cnt[prev], val: prev})
               }
           }
       }
       if !ok {
           continue
       }
       // output answer
       for i := 0; i < m; i++ {
           for j := 0; j < len(ans[i]); j++ {
               if j > 0 {
                   writer.WriteByte(' ')
               }
               fmt.Fprint(writer, ans[i][j])
           }
           writer.WriteByte('\n')
       }
   }
}
