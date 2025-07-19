package main

import (
   "bufio"
   "container/heap"
   "fmt"
   "os"
   "sort"
)

type AddBasic struct {
   typ, id int
   val int64
}

type AddOp struct {
   typ, id int
   a, b int64
}

// MulOp represents an operation in the heap
type MulOp struct {
   typ, id, idx, addIdx int
   a, b int64
}

// MaxHeap implements a max-heap for MulOp by ratio a/b
type MaxHeap []MulOp

func (h MaxHeap) Len() int { return len(h) }
func (h MaxHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }
func (h MaxHeap) Less(i, j int) bool {
   // we want h[i] to be popped first if a/b is larger
   return h[i].a*h[j].b > h[j].a*h[i].b
}
func (h *MaxHeap) Push(x interface{}) {
   *h = append(*h, x.(MulOp))
}
func (h *MaxHeap) Pop() interface{} {
   old := *h
   n := len(old)
   x := old[n-1]
   *h = old[0 : n-1]
   return x
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   var k, n, m int
   fmt.Fscan(in, &k, &n, &m)
   num := make([]int64, k+1)
   for i := 1; i <= k; i++ {
       fmt.Fscan(in, &num[i])
   }
   assignMx := make([]int64, k+1)
   assignId := make([]int, k+1)
   for i := 1; i <= k; i++ {
       assignMx[i] = num[i]
   }
   addsBasic := make([][]AddBasic, k+1)
   // operations type 3 will be pushed to heap directly
   var h MaxHeap
   for j := 1; j <= n; j++ {
       var typ, idx int
       var b int64
       fmt.Fscan(in, &typ, &idx, &b)
       if typ == 1 {
           if b > assignMx[idx] {
               assignMx[idx] = b
               assignId[idx] = j
           }
       } else if typ == 2 {
           addsBasic[idx] = append(addsBasic[idx], AddBasic{typ, j, b})
       } else if typ == 3 {
           // multiply by b => ratio b/1 => a=b-1, b=1
           h = append(h, MulOp{typ: 3, id: j, a: b - 1, b: 1})
       }
   }
   // prepare additions per skill and initial heap entries
   addsOps := make([][]AddOp, k+1)
   for i := 1; i <= k; i++ {
       // include best assignment as an addition
       if assignMx[i] > num[i] {
           addsBasic[i] = append(addsBasic[i], AddBasic{typ: 1, id: assignId[i], val: assignMx[i] - num[i]})
       }
       if len(addsBasic[i]) == 0 {
           continue
       }
       // sort descending by val
       bs := addsBasic[i]
       sort.Slice(bs, func(a, b int) bool {
           return bs[a].val > bs[b].val
       })
       // build AddOp list with precomputed a, b
       sum := num[i]
       for j, it := range bs {
           addsOps[i] = append(addsOps[i], AddOp{typ: it.typ, id: it.id, a: it.val, b: sum})
           sum += it.val
       }
       // include first addition for this skill in heap
       first := addsOps[i][0]
       h = append(h, MulOp{typ: first.typ, id: first.id, idx: i, addIdx: 0, a: first.a, b: first.b})
   }
   // initialize heap with all type3 and first-addition ops
   heap.Init(&h)
   // select up to m operations
   K := m
   var res []struct{ typ, id int }
   for K > 0 && h.Len() > 0 {
       cur := heap.Pop(&h).(MulOp)
       if cur.a <= 0 {
           break
       }
       res = append(res, struct{ typ, id int }{cur.typ, cur.id})
       if cur.typ < 3 {
           // schedule next addition for this skill
           i := cur.idx
           j := cur.addIdx + 1
           if j < len(addsOps[i]) {
               op := addsOps[i][j]
               heap.Push(&h, MulOp{typ: op.typ, id: op.id, idx: i, addIdx: j, a: op.a, b: op.b})
           }
       }
       K--
   }
   // sort selected operations by type (1,2,3)
   sort.Slice(res, func(i, j int) bool { return res[i].typ < res[j].typ })
   // output
   fmt.Fprintln(out, len(res))
   for i, v := range res {
       if i > 0 {
           out.WriteByte(' ')
       }
       fmt.Fprint(out, v.id)
   }
   if len(res) > 0 {
       out.WriteByte('\n')
   }
}
