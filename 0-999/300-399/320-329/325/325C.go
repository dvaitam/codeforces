package main

import (
   "bufio"
   "fmt"
   "os"
)

const capDiam = 314000000

// Item is an entry in priority queue: rule index and its current value
type Item struct { r int; val int64 }

// PQueue is a min-heap of Items
type PQueue struct { data []Item }

// Push inserts an item into the heap
func (h *PQueue) Push(it Item) {
   h.data = append(h.data, it)
   i := len(h.data) - 1
   for i > 0 {
       p := (i - 1) / 2
       if h.data[p].val <= h.data[i].val {
           break
       }
       h.data[p], h.data[i] = h.data[i], h.data[p]
       i = p
   }
}

// Pop removes and returns the smallest item; second return false if empty
func (h *PQueue) Pop() (Item, bool) {
   if len(h.data) == 0 {
       return Item{}, false
   }
   res := h.data[0]
   last := h.data[len(h.data)-1]
   h.data = h.data[:len(h.data)-1]
   if len(h.data) > 0 {
       h.data[0] = last
       i := 0
       for {
           l := 2*i + 1
           r := l + 1
           if l >= len(h.data) {
               break
           }
           mn := l
           if r < len(h.data) && h.data[r].val < h.data[l].val {
               mn = r
           }
           if h.data[mn].val >= h.data[i].val {
               break
           }
           h.data[mn], h.data[i] = h.data[i], h.data[mn]
           i = mn
       }
   }
   return res, true
}

type Rule struct {
   dest int
   src  []int
   d    int
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   var m, n int
   fmt.Fscan(reader, &m, &n)
   rules := make([]Rule, m)
   srcAdj := make([][]int, n+1)
   for i := 0; i < m; i++ {
       var mi, li int
       fmt.Fscan(reader, &mi, &li)
       rules[i].dest = mi
       rules[i].src = make([]int, 0, li)
       cntD := 0
       for j := 0; j < li; j++ {
           var x int
           fmt.Fscan(reader, &x)
           if x < 0 {
               cntD++
           } else {
               rules[i].src = append(rules[i].src, x)
               srcAdj[x] = append(srcAdj[x], i)
           }
       }
       rules[i].d = cntD
   }
   // compute nullable
   remNull := make([]int, m)
   for i := range rules {
       remNull[i] = len(rules[i].src)
   }
   isNull := make([]bool, n+1)
   queueR := make([]int, 0, m)
   for i := range rules {
       if remNull[i] == 0 {
           queueR = append(queueR, i)
       }
   }
   for qi := 0; qi < len(queueR); qi++ {
       r := queueR[qi]
       dest := rules[r].dest
       if isNull[dest] {
           continue
       }
       isNull[dest] = true
       for _, r2 := range srcAdj[dest] {
           remNull[r2]--
           if remNull[r2] == 0 {
               queueR = append(queueR, r2)
           }
       }
   }
   // prepare for min DP
   good := make([]bool, m)
   remMin := make([]int, m)
   sumMin := make([]int64, m)
   for i := range rules {
       ok := true
       for _, j := range rules[i].src {
           if !isNull[j] {
               ok = false
               break
           }
       }
       if ok {
           good[i] = true
           remMin[i] = len(rules[i].src)
           sumMin[i] = int64(rules[i].d)
       }
   }
   // min hyper-Dijkstra
   m0 := make([]int64, n+1)
   const INF = int64(1e18)
   for i := 1; i <= n; i++ {
       m0[i] = INF
   }
   // priority queue for rules (min-heap)
   pq := &PQueue{data: make([]Item, 0, m)}
   for i := range rules {
       if good[i] && remMin[i] == 0 {
           pq.Push(Item{i, sumMin[i]})
       }
   }
   for {
       it, ok := pq.Pop()
       if !ok {
           break
       }
       r := it.r; v := it.val
       dest := rules[r].dest
       if v >= m0[dest] {
           continue
       }
       m0[dest] = v
       for _, r2 := range srcAdj[dest] {
           if !good[r2] {
               continue
           }
           remMin[r2]--
           sumMin[r2] += m0[dest]
           if remMin[r2] == 0 {
               pq.Push(Item{r2, sumMin[r2]})
           }
       }
   }
   // clamp min
   for i := 1; i <= n; i++ {
       if m0[i] > capDiam {
           m0[i] = capDiam
       }
   }
   // build graph for SCC
   adj := make([][]int, n+1)
   rev := make([][]int, n+1)
   selfLoop := make([]bool, n+1)
   for i := range rules {
       if !good[i] {
           continue
       }
       u := rules[i].dest
       for _, v := range rules[i].src {
           adj[u] = append(adj[u], v)
           rev[v] = append(rev[v], u)
           if u == v {
               selfLoop[u] = true
           }
       }
   }
   // Tarjan SCC
   var (
       index   = make([]int, n+1)
       lowlink = make([]int, n+1)
       onstack = make([]bool, n+1)
       stack   []int
       idx     = 1
       sccID   = make([]int, n+1)
       sccCnt  = 0
       sccSz   = make([]int, n+1)
   )
   var dfs func(u int)
   dfs = func(u int) {
       index[u] = idx
       lowlink[u] = idx
       idx++
       stack = append(stack, u)
       onstack[u] = true
       for _, v := range adj[u] {
           if !isNull[v] {
               continue
           }
           if index[v] == 0 {
               dfs(v)
               if lowlink[v] < lowlink[u] {
                   lowlink[u] = lowlink[v]
               }
           } else if onstack[v] && index[v] < lowlink[u] {
               lowlink[u] = index[v]
           }
       }
       if lowlink[u] == index[u] {
           sccCnt++
           for {
               v := stack[len(stack)-1]
               stack = stack[:len(stack)-1]
               onstack[v] = false
               sccID[v] = sccCnt
               sccSz[sccCnt]++
               if v == u {
                   break
               }
           }
       }
   }
   for i := 1; i <= n; i++ {
       if isNull[i] && index[i] == 0 {
           dfs(i)
       }
   }
   infinite := make([]bool, n+1)
   for i := 1; i <= n; i++ {
       if !isNull[i] {
           continue
       }
       id := sccID[i]
       if sccSz[id] > 1 || selfLoop[i] {
           infinite[i] = true
       }
   }
   // propagate infinite backwards
   q := make([]int, 0, n)
   visInf := make([]bool, n+1)
   for i := 1; i <= n; i++ {
       if infinite[i] {
           visInf[i] = true
           q = append(q, i)
       }
   }
   for qi := 0; qi < len(q); qi++ {
       u := q[qi]
       for _, v := range rev[u] {
           if !isNull[v] || visInf[v] {
               continue
           }
           visInf[v] = true
           q = append(q, v)
       }
   }
   // max DP
   remMax := make([]int, m)
   sumMax := make([]int64, m)
   for i := range rules {
       if good[i] {
           remMax[i] = len(rules[i].src)
           sumMax[i] = int64(rules[i].d)
       }
   }
   m1 := make([]int64, n+1)
   rqueue := make([]int, 0, m)
   for i := range rules {
       if good[i] && remMax[i] == 0 {
           rqueue = append(rqueue, i)
       }
   }
   for qi := 0; qi < len(rqueue); qi++ {
       r := rqueue[qi]
       dest := rules[r].dest
       if visInf[dest] {
           continue
       }
       cand := sumMax[r]
       if cand > m1[dest] {
           m1[dest] = cand
       }
       for _, r2 := range srcAdj[dest] {
           if !good[r2] {
               continue
           }
           remMax[r2]--
           sumMax[r2] += m1[dest]
           if remMax[r2] == 0 {
               rqueue = append(rqueue, r2)
           }
       }
   }
   // clamp max
   for i := 1; i <= n; i++ {
       if m1[i] > capDiam {
           m1[i] = capDiam
       }
   }
   // output
   for i := 1; i <= n; i++ {
       if !isNull[i] {
           fmt.Fprintln(writer, "-1 -1")
       } else if visInf[i] {
           fmt.Fprintf(writer, "%d -2\n", m0[i])
       } else {
           fmt.Fprintf(writer, "%d %d\n", m0[i], m1[i])
       }
   }
}
