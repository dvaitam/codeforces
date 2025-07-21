package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

// BIT implements Fenwick tree for int counts
type BIT struct {
   n int
   t []int
}

func NewBIT(n int) *BIT {
   return &BIT{n, make([]int, n+1)}
}

// Add v at position i
func (b *BIT) Add(i, v int) {
   for x := i; x <= b.n; x += x & -x {
       b.t[x] += v
   }
}

// Sum1 returns sum 1..i
func (b *BIT) Sum1(i int) int {
   s := 0
   for x := i; x > 0; x -= x & -x {
       s += b.t[x]
   }
   return s
}

// Sum returns sum l..r
func (b *BIT) Sum(l, r int) int {
   if r < l {
       return 0
   }
   return b.Sum1(r) - b.Sum1(l-1)
}

// Select returns smallest i such that sum1(i) >= k; assume 1 <= k <= total sum
func (b *BIT) Select(k int) int {
   idx := 0
   // compute highest power of two <= n
   bitMask := 1
   for bitMask<<1 <= b.n {
       bitMask <<= 1
   }
   for d := bitMask; d > 0; d >>= 1 {
       nxt := idx + d
       if nxt <= b.n && b.t[nxt] < k {
           idx = nxt
           k -= b.t[nxt]
       }
   }
   return idx + 1
}

type Person struct {
   t  int64
   s, f int
   idx int
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n, m int
   fmt.Fscan(in, &n, &m)
   persons := make([]Person, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &persons[i].t, &persons[i].s, &persons[i].f)
       persons[i].idx = i
   }
   events := make([]Person, n)
   copy(events, persons)
   sort.Slice(events, func(i, j int) bool { return events[i].t < events[j].t })

   waitingBIT := NewBIT(m)
   insideBIT := NewBIT(m)
   waitingQ := make([][]int, m+1)
   insideQ := make([][]int, m+1)
   dest := make([]int, n)
   for i := range persons {
       dest[persons[i].idx] = persons[i].f
   }
   ans := make([]int64, n)

   t := int64(0)
   x := 1
   ptr := 0
   delivered := 0

   // get current pup and pdown, returns dir: +1 up, -1 down, 0 idle
   dir := func() int {
       up := waitingBIT.Sum(x+1, m) + insideBIT.Sum(x+1, m)
       down := waitingBIT.Sum(1, x-1) + insideBIT.Sum(1, x-1)
       if up+down == 0 {
           return 0
       }
       if up >= down {
           return 1
       }
       return -1
   }

   for delivered < n {
       d := dir()
       // idle
       if d == 0 {
           if ptr < n {
               // jump to next arrival
               t = max64(t, events[ptr].t)
               // process all arrivals at t
               for ptr < n && events[ptr].t == t {
                   e := &events[ptr]
                   waitingQ[e.s] = append(waitingQ[e.s], e.idx)
                   waitingBIT.Add(e.s, 1)
                   ptr++
               }
               continue
           } else {
               break
           }
       }
       // next arrival time
       tNextArr := int64(1<<62)
       if ptr < n {
           tNextArr = events[ptr].t
       }
       // find next exit
       var fExit, dExit int64
       if d == 1 && insideBIT.Sum(x+1, m) > 0 {
           // first dest > x
           s0 := insideBIT.Sum1(x)
           fExit = int64(insideBIT.Select(s0 + 1))
           dExit = fExit - int64(x)
       } else if d == -1 && insideBIT.Sum(1, x-1) > 0 {
           s1 := insideBIT.Sum1(x - 1)
           fExit = int64(insideBIT.Select(s1))
           dExit = int64(x) - fExit
       } else {
           dExit = int64(1<<62)
       }
       // find next entry
       var fEntry, dEntry int64
       if d == 1 && waitingBIT.Sum(x+1, m) > 0 {
           s0 := waitingBIT.Sum1(x)
           fEntry = int64(waitingBIT.Select(s0 + 1))
           dEntry = fEntry - int64(x)
       } else if d == -1 && waitingBIT.Sum(1, x-1) > 0 {
           s1 := waitingBIT.Sum1(x - 1)
           fEntry = int64(waitingBIT.Select(s1))
           dEntry = int64(x) - fEntry
       } else {
           dEntry = int64(1<<62)
       }
       // next stop distance
       dEvent := dExit
       fEvent := fExit
       if dEntry < dEvent {
           dEvent = dEntry
           fEvent = fEntry
       }
       // time to reach fEvent
       tEvent := t + dEvent
       if tNextArr <= tEvent {
           // arrival occurs first or same time
           // advance to tNextArr
           dt := tNextArr - t
           x += d * int(dt)
           t = tNextArr
           // process all arrivals at t
           for ptr < n && events[ptr].t == t {
               e := &events[ptr]
               waitingQ[e.s] = append(waitingQ[e.s], e.idx)
               waitingBIT.Add(e.s, 1)
               ptr++
           }
           // process release and entry at current floor x
           // release
           if len(insideQ[x]) > 0 {
               for _, id := range insideQ[x] {
                   ans[id] = t
                   delivered++
               }
               insideBIT.Add(x, -len(insideQ[x]))
               insideQ[x] = nil
           }
           // entry
           if len(waitingQ[x]) > 0 {
               for _, id := range waitingQ[x] {
                   f := dest[id]
                   insideQ[f] = append(insideQ[f], id)
                   insideBIT.Add(f, 1)
               }
               waitingBIT.Add(x, -len(waitingQ[x]))
               waitingQ[x] = nil
           }
       } else {
           // stop at fEvent first
           // advance
           t = tEvent
           x = int(fEvent)
           // release
           if len(insideQ[x]) > 0 {
               for _, id := range insideQ[x] {
                   ans[id] = t
                   delivered++
               }
               insideBIT.Add(x, -len(insideQ[x]))
               insideQ[x] = nil
           }
           // entry
           if len(waitingQ[x]) > 0 {
               for _, id := range waitingQ[x] {
                   f := dest[id]
                   insideQ[f] = append(insideQ[f], id)
                   insideBIT.Add(f, 1)
               }
               waitingBIT.Add(x, -len(waitingQ[x]))
               waitingQ[x] = nil
           }
       }
   }
   // output answers
   for i := 0; i < n; i++ {
       fmt.Fprintln(out, ans[i])
   }
}

func max64(a, b int64) int64 {
   if a > b {
       return a
   }
   return b
}
