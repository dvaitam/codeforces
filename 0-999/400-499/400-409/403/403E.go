package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

// Segment tree for max query with find-first > threshold
type SegMax struct {
   n    int
   tree []int
}

func NewSegMax(n int) *SegMax {
   size := 1
   for size < n {
       size <<= 1
   }
   return &SegMax{n: size, tree: make([]int, 2*size)}
}
func (st *SegMax) Update(pos, val int) {
   i := pos + st.n - 1
   st.tree[i] = val
   for i >>= 1; i > 0; i >>= 1 {
       if st.tree[2*i] > st.tree[2*i+1] {
           st.tree[i] = st.tree[2*i]
       } else {
           st.tree[i] = st.tree[2*i+1]
       }
   }
}
// find any index in [l,r] (1-based) with value > thresh, or return 0
func (st *SegMax) QueryFirst(l, r, thresh int) int {
   return st.query(1, 1, st.n, l, r, thresh)
}
func (st *SegMax) query(idx, lo, hi, l, r, thr int) int {
   if lo > r || hi < l || st.tree[idx] <= thr {
       return 0
   }
   if lo == hi {
       return lo
   }
   mid := (lo + hi) >> 1
   if res := st.query(2*idx, lo, mid, l, r, thr); res != 0 {
       return res
   }
   return st.query(2*idx+1, mid+1, hi, l, r, thr)
}

// Segment tree for min query with find-first < threshold
type SegMin struct {
   n    int
   tree []int
}
func NewSegMin(n int) *SegMin {
   size := 1
   for size < n {
       size <<= 1
   }
   inf := int(1e9)
   tree := make([]int, 2*size)
   for i := range tree {
       tree[i] = inf
   }
   return &SegMin{n: size, tree: tree}
}
func (st *SegMin) Update(pos, val int) {
   i := pos + st.n - 1
   st.tree[i] = val
   for i >>= 1; i > 0; i >>= 1 {
       if st.tree[2*i] < st.tree[2*i+1] {
           st.tree[i] = st.tree[2*i]
       } else {
           st.tree[i] = st.tree[2*i+1]
       }
   }
}
// find any index in [l,r] (1-based) with value < thr, or return 0
func (st *SegMin) QueryFirst(l, r, thr int) int {
   return st.query(1, 1, st.n, l, r, thr)
}
func (st *SegMin) query(idx, lo, hi, l, r, thr int) int {
   if lo > r || hi < l || st.tree[idx] >= thr {
       return 0
   }
   if lo == hi {
       return lo
   }
   mid := (lo + hi) >> 1
   if res := st.query(2*idx, lo, mid, l, r, thr); res != 0 {
       return res
   }
   return st.query(2*idx+1, mid+1, hi, l, r, thr)
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   var n int
   fmt.Fscan(in, &n)
   a := make([]int, n+1)
   for i := 2; i <= n; i++ {
       fmt.Fscan(in, &a[i])
   }
   b := make([]int, n+1)
   for i := 2; i <= n; i++ {
       fmt.Fscan(in, &b[i])
   }
   var idx0 int
   fmt.Fscan(in, &idx0)
   // build blue tree
   blueAdj := make([][]int, n+1)
   for i := 2; i <= n; i++ {
       blueAdj[a[i]] = append(blueAdj[a[i]], i)
   }
   blueTin := make([]int, n+1)
   blueTout := make([]int, n+1)
   timer := 0
   var dfsBlue func(int)
   dfsBlue = func(u int) {
       timer++
       blueTin[u] = timer
       for _, v := range blueAdj[u] {
           dfsBlue(v)
       }
       blueTout[u] = timer
   }
   dfsBlue(1)
   // build red tree
   redAdj := make([][]int, n+1)
   for i := 2; i <= n; i++ {
       redAdj[b[i]] = append(redAdj[b[i]], i)
   }
   redTin := make([]int, n+1)
   redTout := make([]int, n+1)
   timer = 0
   var dfsRed func(int)
   dfsRed = func(u int) {
       timer++
       redTin[u] = timer
       for _, v := range redAdj[u] {
           dfsRed(v)
       }
       redTout[u] = timer
   }
   dfsRed(1)
   // prepare red edges data on blue tin
   aaR := make([]int, n) // 1-based edges
   bbR := make([]int, n)
   listARA := make([][]int, n+2)
   listARB := make([][]int, n+2)
   for ei := 1; ei < n; ei++ {
       u := b[ei+1]
       v := ei + 1
       // endpoints u and v
       t1 := blueTin[u]
       t2 := blueTin[v]
       if t1 < t2 {
           aaR[ei] = t1; bbR[ei] = t2
       } else {
           aaR[ei] = t2; bbR[ei] = t1
       }
       listARA[aaR[ei]] = append(listARA[aaR[ei]], ei)
       listARB[bbR[ei]] = append(listARB[bbR[ei]], ei)
   }
   // sort per list
   for i := 1; i <= n; i++ {
       // listARA by bb desc
       if len(listARA[i]) > 1 {
           sort.Slice(listARA[i], func(x, y int) bool {
               return bbR[listARA[i][x]] > bbR[listARA[i][y]]
           })
       }
       // listARB by aa asc
       if len(listARB[i]) > 1 {
           sort.Slice(listARB[i], func(x, y int) bool {
               return aaR[listARB[i][x]] < aaR[listARB[i][y]]
           })
       }
   }
   // initialize segment trees for red edges
   segRA := NewSegMax(n)
   segRB := NewSegMin(n)
   // pointers to current front in lists
   ptrRA := make([]int, n+2)
   ptrRB := make([]int, n+2)
   const INF = int(1e9)
   for i := 1; i <= n; i++ {
       if ptrRA[i] < len(listARA[i]) {
           eid := listARA[i][ptrRA[i]]
           segRA.Update(i, bbR[eid])
       }
       if ptrARB[i] < len(listARB[i]) {
           eid := listARB[i][ptrRB[i]]
           segRB.Update(i, aaR[eid])
       }
   }
   // prepare blue edges data on red tin similarly
   aaB := make([]int, n)
   bbB := make([]int, n)
   listBAA := make([][]int, n+2)
   listBAB := make([][]int, n+2)
   for ei := 1; ei < n; ei++ {
       u := a[ei+1]
       v := ei + 1
       t1 := redTin[u]
       t2 := redTin[v]
       if t1 < t2 {
           aaB[ei] = t1; bbB[ei] = t2
       } else {
           aaB[ei] = t2; bbB[ei] = t1
       }
       listBAA[aaB[ei]] = append(listBAA[aaB[ei]], ei)
       listBAB[bbB[ei]] = append(listBAB[bbB[ei]], ei)
   }
   for i := 1; i <= n; i++ {
       if len(listBAA[i]) > 1 {
           sort.Slice(listBAA[i], func(x, y int) bool {
               return bbB[listBAA[i][x]] > bbB[listBAA[i][y]]
           })
       }
       if len(listBAB[i]) > 1 {
           sort.Slice(listBAB[i], func(x, y int) bool {
               return aaB[listBAB[i][x]] < aaB[listBAB[i][y]]
           })
       }
   }
   segBA := NewSegMax(n)
   segBB := NewSegMin(n)
   ptrBA := make([]int, n+2)
   ptrBB := make([]int, n+2)
   for i := 1; i <= n; i++ {
       if ptrBA[i] < len(listBAA[i]) {
           eid := listBAA[i][ptrBA[i]]
           segBA.Update(i, bbB[eid])
       }
       if ptrBB[i] < len(listBAB[i]) {
           eid := listBAB[i][ptrBB[i]]
           segBB.Update(i, aaB[eid])
       }
   }
   // BFS
   deletedR := make([]bool, n)
   deletedB := make([]bool, n)
   type Stage struct { color string; edges []int }
   var stages []Stage
   blueQ := []int{idx0}
   for len(blueQ) > 0 {
       // Blue stage
       sort.Ints(blueQ)
       stages = append(stages, Stage{"Blue", append([]int(nil), blueQ...)})
       // process blueQ to find redQ
       redQ := make([]int, 0)
       for _, be := range blueQ {
           v := be + 1
           l, r := blueTin[v], blueTout[v]
           // case1: b in [l,r] with a<l using segRB
           for {
               bb := segRB.QueryFirst(l, r, l)
               if bb == 0 { break }
               // get edge
               // pop until find undeleted and a<l
               for ptrRB[bb] < len(listARB[bb]) {
                   e := listARB[bb][ptrRB[bb]]
                   if deletedR[e] || aaR[e] >= l {
                       ptrRB[bb]++
                       continue
                   }
                   // delete e
                   deletedR[e] = true
                   redQ = append(redQ, e)
                   ptrRB[bb]++
                   // update segRB at bb
                   if ptrRB[bb] < len(listARB[bb]) {
                       segRB.Update(bb, aaR[listARB[bb][ptrRB[bb]]])
                   } else {
                       segRB.Update(bb, INF)
                   }
                   // also update segRA at aaR[e]
                   aa0 := aaR[e]
                   for ptrRA[aa0] < len(listARA[aa0]) && deletedR[listARA[aa0][ptrRA[aa0]]] {
                       ptrRA[aa0]++
                   }
                   if ptrRA[aa0] < len(listARA[aa0]) {
                       segRA.Update(aa0, bbR[listARA[aa0][ptrRA[aa0]]])
                   } else {
                       segRA.Update(aa0, 0)
                   }
                   break
               }
           }
           // case2: a in [l,r] with b>r using segRA
           for {
               aa := segRA.QueryFirst(l, r, r)
               if aa == 0 { break }
               for ptrRA[aa] < len(listARA[aa]) {
                   e := listARA[aa][ptrRA[aa]]
                   if deletedR[e] || bbR[e] <= r {
                       ptrRA[aa]++
                       continue
                   }
                   deletedR[e] = true
                   redQ = append(redQ, e)
                   ptrRA[aa]++
                   if ptrRA[aa] < len(listARA[aa]) {
                       segRA.Update(aa, bbR[listARA[aa][ptrRA[aa]]])
                   } else {
                       segRA.Update(aa, 0)
                   }
                   // update segRB at bbR[e]
                   bb0 := bbR[e]
                   for ptrRB[bb0] < len(listARB[bb0]) && deletedR[listARB[bb0][ptrRB[bb0]]] {
                       ptrRB[bb0]++
                   }
                   if ptrRB[bb0] < len(listARB[bb0]) {
                       segRB.Update(bb0, aaR[listARB[bb0][ptrRB[bb0]]])
                   } else {
                       segRB.Update(bb0, INF)
                   }
                   break
               }
           }
       }
       if len(redQ) == 0 { break }
       // Red stage
       sort.Ints(redQ)
       stages = append(stages, Stage{"Red", append([]int(nil), redQ...)})
       // process redQ to find next blueQ
       blueQ = make([]int, 0)
       for _, re := range redQ {
           v := re + 1
           l, r := redTin[v], redTout[v]
           // case1 on blue: b in [l,r] with a<l
           for {
               bb := segBB.QueryFirst(l, r, l)
               if bb == 0 { break }
               for ptrBB[bb] < len(listBAB[bb]) {
                   e := listBAB[bb][ptrBB[bb]]
                   if deletedB[e] || aaB[e] >= l {
                       ptrBB[bb]++
                       continue
                   }
                   deletedB[e] = true
                   blueQ = append(blueQ, e)
                   ptrBB[bb]++
                   if ptrBB[bb] < len(listBAB[bb]) {
                       segBB.Update(bb, aaB[listBAB[bb][ptrBB[bb]]])
                   } else {
                       segBB.Update(bb, INF)
                   }
                   // update segBA
                   aa0 := aaB[e]
                   for ptrBA[aa0] < len(listBAA[aa0]) && deletedB[listBAA[aa0][ptrBA[aa0]]] {
                       ptrBA[aa0]++
                   }
                   if ptrBA[aa0] < len(listBAA[aa0]) {
                       segBA.Update(aa0, bbB[listBAA[aa0][ptrBA[aa0]]])
                   } else {
                       segBA.Update(aa0, 0)
                   }
                   break
               }
           }
           // case2 on blue: a in [l,r] with b>r
           for {
               aa := segBA.QueryFirst(l, r, r)
               if aa == 0 { break }
               for ptrBA[aa] < len(listBAA[aa]) {
                   e := listBAA[aa][ptrBA[aa]]
                   if deletedB[e] || bbB[e] <= r {
                       ptrBA[aa]++
                       continue
                   }
                   deletedB[e] = true
                   blueQ = append(blueQ, e)
                   ptrBA[aa]++
                   if ptrBA[aa] < len(listBAA[aa]) {
                       segBA.Update(aa, bbB[listBAA[aa][ptrBA[aa]]])
                   } else {
                       segBA.Update(aa, 0)
                   }
                   // update segBB
                   bb0 := bbB[e]
                   for ptrBB[bb0] < len(listBAB[bb0]) && deletedB[listBAB[bb0][ptrBB[bb0]]] {
                       ptrBB[bb0]++
                   }
                   if ptrBB[bb0] < len(listBAB[bb0]) {
                       segBB.Update(bb0, aaB[listBAB[bb0][ptrBB[bb0]]])
                   } else {
                       segBB.Update(bb0, INF)
                   }
                   break
               }
           }
       }
   }
   // output
   for _, stg := range stages {
       fmt.Fprintln(out, stg.color)
       for i, e := range stg.edges {
           if i > 0 {
               out.WriteByte(' ')
           }
           fmt.Fprint(out, e)
       }
       out.WriteByte('\n')
   }
}
