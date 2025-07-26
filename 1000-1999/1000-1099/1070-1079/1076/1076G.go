package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   var n, m, q int
   fmt.Fscan(reader, &n, &m, &q)
   a := make([]int, n+1)
   for i := 1; i <= n; i++ {
       var x int64
       fmt.Fscan(reader, &x)
       a[i] = int(x & 1) // store parity of a_i
   }
   // build segment tree for XOR with range flip
   size := 1
   for size < n {
       size <<= 1
   }
   seg := make([]int, 2*size)
   lazy := make([]bool, 2*size)
   // initialize leaves
   for i := 0; i < n; i++ {
       seg[size+i] = a[i+1]
   }
   for i := size - 1; i > 0; i-- {
       seg[i] = seg[i<<1] ^ seg[i<<1|1]
   }
   // apply flip on node idx covering length len
   var applyFlip func(idx, length int)
   applyFlip = func(idx, length int) {
       if length&1 == 1 {
           seg[idx] ^= 1
       }
       lazy[idx] = !lazy[idx]
   }
   // push down lazy flag
   var push func(idx, l, r int)
   push = func(idx, l, r int) {
       if !lazy[idx] || idx >= size {
           return
       }
       mid := (l + r) >> 1
       applyFlip(idx<<1, mid-l+1)
       applyFlip(idx<<1|1, r-mid)
       lazy[idx] = false
   }
   // update range [ql,qr]
   var update func(idx, l, r, ql, qr int)
   update = func(idx, l, r, ql, qr int) {
       if ql > r || qr < l {
           return
       }
       if ql <= l && r <= qr {
           applyFlip(idx, r-l+1)
           return
       }
       push(idx, l, r)
       mid := (l + r) >> 1
       update(idx<<1, l, mid, ql, qr)
       update(idx<<1|1, mid+1, r, ql, qr)
       seg[idx] = seg[idx<<1] ^ seg[idx<<1|1]
   }
   // query XOR in [ql,qr]
   var query func(idx, l, r, ql, qr int) int
   query = func(idx, l, r, ql, qr int) int {
       if ql > r || qr < l {
           return 0
       }
       if ql <= l && r <= qr {
           return seg[idx]
       }
       push(idx, l, r)
       mid := (l + r) >> 1
       return query(idx<<1, l, mid, ql, qr) ^ query(idx<<1|1, mid+1, r, ql, qr)
   }
   // process queries
   for i := 0; i < q; i++ {
       var tp int
       fmt.Fscan(reader, &tp)
       if tp == 1 {
           var l, r int
           var d int64
           fmt.Fscan(reader, &l, &r, &d)
           if d&1 == 1 {
               update(1, 1, size, l, r)
           }
       } else {
           var l, r int
           fmt.Fscan(reader, &l, &r)
           // parity of sum = XOR of bits in [l,r]
           res := query(1, 1, size, l, r)
           if res == 0 {
               fmt.Fprintln(writer, 1)
           } else {
               fmt.Fprintln(writer, 2)
           }
       }
   }
