package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

var (
   A, B int64
   p    []int64
   w    *bufio.Writer
)

func work(l, beg, end int, st int64) int64 {
   if beg > end {
       return A
   }
   if l == 0 {
       return int64(end-beg+1) * B
   }
   half := int64(1) << (l - 1)
   m := st + half
   // find first index in p[beg:end+1] where p[i] >= m
   cnt := end - beg + 1
   idx := sort.Search(cnt, func(i int) bool { return p[beg+i] >= m })
   mid := beg + idx
   total := int64(cnt)
   costAll := total * (int64(1)<<l) * B
   costSplit := work(l-1, beg, mid-1, st) + work(l-1, mid, end, m)
   if costAll < costSplit {
       return costAll
   }
   return costSplit
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   w = bufio.NewWriter(os.Stdout)
   defer w.Flush()
   var n int
   var k int
   if _, err := fmt.Fscan(reader, &n, &k, &A, &B); err != nil {
       return
   }
   p = make([]int64, k)
   for i := 0; i < k; i++ {
       fmt.Fscan(reader, &p[i])
   }
   sort.Slice(p, func(i, j int) bool { return p[i] < p[j] })
   res := work(n, 0, k-1, 1)
   fmt.Fprintln(w, res)
}
