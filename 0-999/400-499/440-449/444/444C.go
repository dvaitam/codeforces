package main

import (
   "bufio"
   "fmt"
   "os"
)

var (
   n, m int
   sumF []int64
   add  []int64
   col  []int
)

func abs(x int) int {
   if x < 0 {
       return -x
   }
   return x
}

func build(idx, l, r int) {
   add[idx] = 0
   sumF[idx] = 0
   if l == r {
       col[idx] = l
   } else {
       mid := (l + r) >> 1
       build(idx<<1, l, mid)
       build(idx<<1|1, mid+1, r)
       if col[idx<<1] == col[idx<<1|1] {
           col[idx] = col[idx<<1]
       } else {
           col[idx] = -1
       }
   }
}

func pushdown(idx, l, r int) {
   if l == r {
       return
   }
   left, right := idx<<1, idx<<1|1
   mid := (l + r) >> 1
   if col[idx] != -1 {
       // propagate color and add to children
       col[left] = col[idx]
       col[right] = col[idx]
       // propagate accumulated f addition
       if add[idx] != 0 {
           add[left] += add[idx]
           sumF[left] += add[idx] * int64(mid-l+1)
           add[right] += add[idx]
           sumF[right] += add[idx] * int64(r-mid)
       }
       // clear tag
       col[idx] = -1
       add[idx] = 0
   } else if add[idx] != 0 {
       // propagate only f addition
       add[left] += add[idx]
       sumF[left] += add[idx] * int64(mid-l+1)
       add[right] += add[idx]
       sumF[right] += add[idx] * int64(r-mid)
       add[idx] = 0
   }
}

func update(idx, l, r, ql, qr, x int) {
   if qr < l || r < ql {
       return
   }
   if ql <= l && r <= qr && col[idx] != -1 {
       // fully covered and uniform color
       d := abs(col[idx] - x)
       add[idx] += int64(d)
       sumF[idx] += int64(d) * int64(r-l+1)
       col[idx] = x
       return
   }
   pushdown(idx, l, r)
   mid := (l + r) >> 1
   if ql <= mid {
       update(idx<<1, l, mid, ql, qr, x)
   }
   if qr > mid {
       update(idx<<1|1, mid+1, r, ql, qr, x)
   }
   sumF[idx] = sumF[idx<<1] + sumF[idx<<1|1]
   // merge color
   if col[idx<<1] != -1 && col[idx<<1] == col[idx<<1|1] {
       col[idx] = col[idx<<1]
   } else {
       col[idx] = -1
   }
}

func query(idx, l, r, ql, qr int) int64 {
   if qr < l || r < ql {
       return 0
   }
   if ql <= l && r <= qr {
       return sumF[idx]
   }
   pushdown(idx, l, r)
   mid := (l + r) >> 1
   var res int64
   if ql <= mid {
       res += query(idx<<1, l, mid, ql, qr)
   }
   if qr > mid {
       res += query(idx<<1|1, mid+1, r, ql, qr)
   }
   return res
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   readInt := func() int {
       var x int
       var sign int = 1
       // skip non-digit
       b, err := reader.ReadByte()
       for err == nil && (b < '0' || b > '9') && b != '-' {
           b, err = reader.ReadByte()
       }
       if err != nil {
           return 0
       }
       if b == '-' {
           sign = -1
           b, _ = reader.ReadByte()
       }
       for err == nil && b >= '0' && b <= '9' {
           x = x*10 + int(b-'0')
           b, err = reader.ReadByte()
       }
       return x * sign
   }

   n = readInt()
   m = readInt()
   sumF = make([]int64, 4*(n+5))
   add = make([]int64, 4*(n+5))
   col = make([]int, 4*(n+5))
   build(1, 1, n)
   for i := 0; i < m; i++ {
       t := readInt()
       if t == 1 {
           l := readInt()
           r := readInt()
           x := readInt()
           update(1, 1, n, l, r, x)
       } else {
           l := readInt()
           r := readInt()
           ans := query(1, 1, n, l, r)
           fmt.Fprintln(writer, ans)
       }
   }
}
