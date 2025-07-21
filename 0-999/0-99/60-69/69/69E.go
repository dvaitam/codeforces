package main

import (
   "bufio"
   "fmt"
   "os"
)

// MaxHeap implements a max-heap of ints.
type MaxHeap struct { data []int }

func (h *MaxHeap) push(x int) {
   h.data = append(h.data, x)
   h.up(len(h.data) - 1)
}

func (h *MaxHeap) pop() int {
   n := len(h.data)
   top := h.data[0]
   h.data[0] = h.data[n-1]
   h.data = h.data[:n-1]
   if len(h.data) > 0 {
       h.down(0)
   }
   return top
}

func (h *MaxHeap) up(i int) {
   for i > 0 {
       p := (i - 1) / 2
       if h.data[i] > h.data[p] {
           h.data[i], h.data[p] = h.data[p], h.data[i]
           i = p
           continue
       }
       break
   }
}

func (h *MaxHeap) down(i int) {
   n := len(h.data)
   for {
       l := 2*i + 1
       if l >= n {
           break
       }
       largest := l
       r := l + 1
       if r < n && h.data[r] > h.data[l] {
           largest = r
       }
       if h.data[largest] > h.data[i] {
           h.data[i], h.data[largest] = h.data[largest], h.data[i]
           i = largest
       } else {
           break
       }
   }
}

func (h *MaxHeap) top() (int, bool) {
   if len(h.data) == 0 {
       return 0, false
   }
   return h.data[0], true
}

func (h *MaxHeap) len() int {
   return len(h.data)
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, k int
   if _, err := fmt.Fscan(reader, &n, &k); err != nil {
       return
   }
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }

   cnt := make(map[int]int, n)
   h := &MaxHeap{data: make([]int, 0, k)}
   // Initialize first window
   for i := 0; i < k; i++ {
       v := a[i]
       cnt[v]++
       if cnt[v] == 1 {
           h.push(v)
       }
   }

   // Process windows
   for i := 0; i <= n-k; i++ {
       // Clean invalid entries
       // Clean invalid entries
       for {
           top, ok := h.top()
           if !ok || cnt[top] != 1 {
               if !ok {
                   break
               }
               h.pop()
               continue
           }
           break
       }
       if h.len() == 0 {
           writer.WriteString("Nothing\n")
       } else {
           top, _ := h.top()
           fmt.Fprintln(writer, top)
       }
       // Slide window
       if i == n-k {
           break
       }
       // Remove a[i]
       v := a[i]
       prev := cnt[v]
       cnt[v] = prev - 1
       if cnt[v] == 1 {
           h.push(v)
       }
       // Add a[i+k]
       u := a[i+k]
       prev2 := cnt[u]
       cnt[u] = prev2 + 1
       if cnt[u] == 1 {
           h.push(u)
       }
   }
}
