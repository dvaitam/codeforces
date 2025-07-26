package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

// Fenwick tree for int64 values
type BIT struct {
   n    int
   tree []int64
}

func NewBIT(n int) *BIT {
   return &BIT{n: n, tree: make([]int64, n+1)}
}

// Update adds v at index i (1-indexed)
func (b *BIT) Update(i int, v int64) {
   for ; i <= b.n; i += i & -i {
       b.tree[i] += v
   }
}

// Query returns sum in [1..i]
func (b *BIT) Query(i int) int64 {
   var s int64
   for ; i > 0; i -= i & -i {
       s += b.tree[i]
   }
   return s
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   xs := make([]int, n)
   vs := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &xs[i])
   }
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &vs[i])
   }
   // sort points by x
   idx := make([]int, n)
   for i := 0; i < n; i++ {
       idx[i] = i
   }
   sort.Slice(idx, func(i, j int) bool {
       return xs[idx[i]] < xs[idx[j]]
   })
   sortedX := make([]int64, n)
   sortedV := make([]int, n)
   for i, id := range idx {
       sortedX[i] = int64(xs[id])
       sortedV[i] = vs[id]
   }
   // compress velocities
   uniqV := make([]int, n)
   copy(uniqV, sortedV)
   sort.Ints(uniqV)
   m := 0
   comp := make(map[int]int, n)
   for _, v := range uniqV {
       if _, ok := comp[v]; !ok {
           m++
           comp[v] = m
       }
   }
   bitCount := NewBIT(m)
   bitSum := NewBIT(m)
   var res int64
   // process in order of increasing x
   for i := 0; i < n; i++ {
       v := sortedV[i]
       xi := sortedX[i]
       ci := comp[v]
       // query all with v_i <= v
       cnt := bitCount.Query(ci)
       sumX := bitSum.Query(ci)
       res += cnt*xi - sumX
       // add current
       bitCount.Update(ci, 1)
       bitSum.Update(ci, xi)
   }
   // output result
   fmt.Println(res)
}
