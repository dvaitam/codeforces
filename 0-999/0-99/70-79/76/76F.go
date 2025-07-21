package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

const NEG_INF = -1000000

// Event holds transformed coordinates for DP
type Event struct {
   A    int64 // x + V*t
   Bidx int   // compressed index for B = x - V*t
   t    int   // time
   x    int64 // original x
}

// BIT supports prefix maximum query and point updates, with clearing
type BIT struct {
   n        int
   tree     []int
   modified []int
}

// NewBIT creates a BIT of size n
func NewBIT(n int) *BIT {
   tree := make([]int, n+1)
   for i := 1; i <= n; i++ {
       tree[i] = NEG_INF
   }
   return &BIT{n: n, tree: tree, modified: make([]int, 0)}
}

// update sets position pos to max(current, val)
func (b *BIT) update(pos, val int) {
   for i := pos; i <= b.n; i += i & -i {
       if b.tree[i] < val {
           b.tree[i] = val
           b.modified = append(b.modified, i)
       }
   }
}

// query returns max over [1..pos]
func (b *BIT) query(pos int) int {
   res := NEG_INF
   for i := pos; i > 0; i -= i & -i {
       if b.tree[i] > res {
           res = b.tree[i]
       }
   }
   return res
}

// clear resets all modified entries to NEG_INF
func (b *BIT) clear() {
   for _, i := range b.modified {
       b.tree[i] = NEG_INF
   }
   b.modified = b.modified[:0]
}

// cdqDP performs CDQ divide-and-conquer DP on events[l:r]
func cdqDP(events []Event, dp []int, bit *BIT, l, r int) {
   if r-l <= 1 {
       return
   }
   // split ensuring no equal times across halves
   mid0 := (l + r) >> 1
   tmid := events[mid0].t
   mid := mid0 + 1
   for mid < r && events[mid].t == tmid {
       mid++
   }
   if mid == r {
       mid = mid0
       for mid > l && events[mid-1].t == tmid {
           mid--
       }
   }
   if mid <= l || mid >= r {
       return
   }
   // process left half
   cdqDP(events, dp, bit, l, mid)
   // prepare indices sorted by A
   left := make([]int, mid-l)
   for i := l; i < mid; i++ {
       left[i-l] = i
   }
   right := make([]int, r-mid)
   for i := mid; i < r; i++ {
       right[i-mid] = i
   }
   sort.Slice(left, func(i, j int) bool {
       return events[left[i]].A < events[left[j]].A
   })
   sort.Slice(right, func(i, j int) bool {
       return events[right[i]].A < events[right[j]].A
   })
   // merge updates
   pi := 0
   for _, i := range right {
       for pi < len(left) && events[left[pi]].A <= events[i].A {
           idxj := left[pi]
           bit.update(events[idxj].Bidx, dp[idxj])
           pi++
       }
       best := bit.query(events[i].Bidx)
       if best+1 > dp[i] {
           dp[i] = best + 1
       }
   }
   // clear BIT for next use
   bit.clear()
   // process right half
   cdqDP(events, dp, bit, mid, r)
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var N int
   fmt.Fscan(in, &N)
   rawX := make([]int64, N)
   rawT := make([]int, N)
   for i := 0; i < N; i++ {
       fmt.Fscan(in, &rawX[i], &rawT[i])
   }
   var V int64
   fmt.Fscan(in, &V)
   // compute A and B
   Bs := make([]int64, N)
   events := make([]Event, N)
   for i := 0; i < N; i++ {
       x := rawX[i]
       t := rawT[i]
       A := x + V*int64(t)
       B := x - V*int64(t)
       Bs[i] = B
       events[i] = Event{A: A, Bidx: 0, t: t, x: x}
   }
   // compress B values (ascending Bs, then Bidx = M - pos)
   sort.Slice(Bs, func(i, j int) bool { return Bs[i] < Bs[j] })
   uniq := Bs[:0]
   for i, v := range Bs {
       if i == 0 || v != Bs[i-1] {
           uniq = append(uniq, v)
       }
   }
   M := len(uniq)
   // map B to index
   bmap := make(map[int64]int, M)
   for i, v := range uniq {
       // largest B has highest pos i=M-1 => idx=1
       bmap[v] = M - i
   }
   for i := 0; i < N; i++ {
       B := rawX[i] - V*int64(rawT[i])
       events[i].Bidx = bmap[B]
   }
   // sort events by time
   sort.Slice(events, func(i, j int) bool {
       return events[i].t < events[j].t
   })
   // DP1: starting at origin
   dp1 := make([]int, N)
   for i := 0; i < N; i++ {
       // reachable from (0,0)?
       if abs64(events[i].x) <= V*int64(events[i].t) {
           dp1[i] = 1
       } else {
           dp1[i] = NEG_INF
       }
   }
   bit1 := NewBIT(M)
   cdqDP(events, dp1, bit1, 0, N)
   ans1 := 0
   for i := 0; i < N; i++ {
       if dp1[i] > ans1 {
           ans1 = dp1[i]
       }
   }
   // DP2: free start
   dp2 := make([]int, N)
   for i := 0; i < N; i++ {
       dp2[i] = 1
   }
   bit2 := NewBIT(M)
   cdqDP(events, dp2, bit2, 0, N)
   ans2 := 0
   for i := 0; i < N; i++ {
       if dp2[i] > ans2 {
           ans2 = dp2[i]
       }
   }
   // output results
   out := bufio.NewWriter(os.Stdout)
   fmt.Fprintf(out, "%d %d", ans1, ans2)
   out.Flush()
}

func abs64(x int64) int64 {
   if x < 0 {
       return -x
   }
   return x
}
