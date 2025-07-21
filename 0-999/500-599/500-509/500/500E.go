package main

import (
   "bufio"
   "fmt"
   "os"
)

type block struct {
   start int   // start index of block (0-based)
   reach int64 // maximum reach coordinate of the block
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n int
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   p := make([]int64, n)
   r := make([]int64, n)
   for i := 0; i < n; i++ {
       var li int64
       fmt.Fscan(in, &p[i], &li)
       r[i] = p[i] + li
   }
   // build blocks of overlapping intervals [p[i], r[i]]
   blocks := make([]block, 0, n)
   for i := 0; i < n; i++ {
       start := i
       reach := r[i]
       // merge with previous blocks if overlapping
       for len(blocks) > 0 && blocks[len(blocks)-1].reach >= p[i] {
           prev := blocks[len(blocks)-1]
           blocks = blocks[:len(blocks)-1]
           if prev.start < start {
               start = prev.start
           }
           if prev.reach > reach {
               reach = prev.reach
           }
       }
       blocks = append(blocks, block{start, reach})
   }
   // assign block id for each domino
   bcnt := len(blocks)
   blk := make([]int, n)
   for bi := 0; bi < bcnt; bi++ {
       s := blocks[bi].start
       var e int
       if bi+1 < bcnt {
           e = blocks[bi+1].start - 1
       } else {
           e = n - 1
       }
       for i := s; i <= e; i++ {
           blk[i] = bi
       }
   }
   // compute gaps between blocks
   gaps := make([]int64, bcnt-1)
   for bi := 0; bi+1 < bcnt; bi++ {
       nextStart := blocks[bi+1].start
       gaps[bi] = p[nextStart] - blocks[bi].reach
   }
   // prefix sums of gaps
   ps := make([]int64, bcnt)
   for i := 1; i < bcnt; i++ {
       ps[i] = ps[i-1] + gaps[i-1]
   }
   // process queries
   var q int
   fmt.Fscan(in, &q)
   for qi := 0; qi < q; qi++ {
       var x, y int
       fmt.Fscan(in, &x, &y)
       x--
       y--
       bx := blk[x]
       by := blk[y]
       if bx >= by {
           fmt.Fprintln(out, 0)
       } else {
           fmt.Fprintln(out, ps[by]-ps[bx])
       }
   }
}
