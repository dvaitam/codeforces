package main

import (
   "bufio"
   "fmt"
   "os"
)

const D = 5

type Node struct {
   num, same int
   a, b      [D]int
}

var n, m, p0, p int
var arr []int
var st []Node

func combine(n1, n2 Node) Node {
   t := Node{}
   t.same = 0
   t.num = n1.num
   for i := 0; i < t.num; i++ {
       t.a[i] = n1.a[i]
       t.b[i] = n1.b[i]
   }
   for i := 0; i < n2.num; i++ {
       ai, bi := n2.a[i], n2.b[i]
       j := 0
       for ; j < t.num; j++ {
           if t.a[j] == ai {
               t.b[j] += bi
               break
           }
       }
       if j < t.num {
           continue
       }
       if t.num < p {
           t.a[t.num] = ai
           t.b[t.num] = bi
           t.num++
           continue
       }
       // find minimal b
       k := 0
       for j = 1; j < t.num; j++ {
           if t.b[j] < t.b[k] {
               k = j
           }
       }
       if bi < t.b[k] {
           for j = 0; j < t.num; j++ {
               t.b[j] -= bi
           }
       } else {
           tmp := t.b[k]
           t.a[k] = ai
           t.b[k] = bi
           for j = 0; j < t.num; j++ {
               t.b[j] -= tmp
           }
       }
   }
   return t
}

func build(l, r, idx int) {
   if l == r {
       st[idx].num = 1
       st[idx].same = 0
       st[idx].a[0] = arr[l]
       st[idx].b[0] = 1
       return
   }
   mid := (l + r) >> 1
   build(l, mid, idx<<1)
   build(mid+1, r, idx<<1|1)
   st[idx] = combine(st[idx<<1], st[idx<<1|1])
}

func pushDown(idx, l, r int) {
   if st[idx].same != 0 {
       mid := (l + r) >> 1
       lc, rc := idx<<1, idx<<1|1
       v := st[idx].same
       st[lc].same = v
       st[rc].same = v
       st[lc].num = 1
       st[rc].num = 1
       st[lc].a[0] = v
       st[rc].a[0] = v
       st[lc].b[0] = mid - l + 1
       st[rc].b[0] = r - mid
       st[idx].same = 0
   }
}

func update(l, r, ql, qr, v, idx int) {
   if ql <= l && r <= qr {
       st[idx].same = v
       st[idx].num = 1
       st[idx].a[0] = v
       st[idx].b[0] = r - l + 1
       return
   }
   pushDown(idx, l, r)
   mid := (l + r) >> 1
   if qr <= mid {
       update(l, mid, ql, qr, v, idx<<1)
   } else if ql > mid {
       update(mid+1, r, ql, qr, v, idx<<1|1)
   } else {
       update(l, mid, ql, mid, v, idx<<1)
       update(mid+1, r, mid+1, qr, v, idx<<1|1)
   }
   st[idx] = combine(st[idx<<1], st[idx<<1|1])
}

func query(l, r, ql, qr, idx int) Node {
   if ql <= l && r <= qr {
       return st[idx]
   }
   pushDown(idx, l, r)
   mid := (l + r) >> 1
   if qr <= mid {
       return query(l, mid, ql, qr, idx<<1)
   }
   if ql > mid {
       return query(mid+1, r, ql, qr, idx<<1|1)
   }
   left := query(l, mid, ql, mid, idx<<1)
   right := query(mid+1, r, mid+1, qr, idx<<1|1)
   return combine(left, right)
}

func main() {
   br := bufio.NewReader(os.Stdin)
   bw := bufio.NewWriter(os.Stdout)
   defer bw.Flush()
   readInt := func() int {
       x := 0; sign := 1
       ch, _ := br.ReadByte()
       for (ch < '0' || ch > '9') && ch != '-' {
           ch, _ = br.ReadByte()
       }
       if ch == '-' {
           sign = -1
           ch, _ = br.ReadByte()
       }
       for ch >= '0' && ch <= '9' {
           x = x*10 + int(ch-'0')
           ch, _ = br.ReadByte()
       }
       return x * sign
   }
   n = readInt()
   m = readInt()
   p0 = readInt()
   p = 100 / p0
   arr = make([]int, n+1)
   for i := 1; i <= n; i++ {
       arr[i] = readInt()
   }
   st = make([]Node, 4*n+4)
   build(1, n, 1)
   for i := 0; i < m; i++ {
       o := readInt()
       l := readInt()
       r := readInt()
       if o == 1 {
           x := readInt()
           update(1, n, l, r, x, 1)
       } else {
           ans := query(1, n, l, r, 1)
           fmt.Fprint(bw, ans.num)
           for j := 0; j < ans.num; j++ {
               fmt.Fprint(bw, " ", ans.a[j])
           }
           fmt.Fprint(bw, '\n')
       }
   }
}
