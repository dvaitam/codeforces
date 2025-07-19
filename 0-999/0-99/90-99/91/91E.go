package main

import (
   "bufio"
   "os"
   "sort"
   "strconv"
)

const Mdeep = 20

var (
   a, b []int64
   tz, tp []int
   Tr [][]int
   idArr []int
)

type Query struct {
   id, L, R, T int
}

// read ints quickly
var reader = bufio.NewReader(os.Stdin)
var writer = bufio.NewWriter(os.Stdout)

func readInt() int {
   var c byte
   var err error
   // skip non-digits
   for {
       c, err = reader.ReadByte()
       if err != nil {
           return 0
       }
       if (c >= '0' && c <= '9') || c == '-' {
           break
       }
   }
   sign := 1
   if c == '-' {
       sign = -1
       c, _ = reader.ReadByte()
   }
   x := 0
   for c >= '0' && c <= '9' {
       x = x*10 + int(c-'0')
       c, err = reader.ReadByte()
       if err != nil {
           break
       }
   }
   return x * sign
}

func cmp(i, j int) bool {
   if a[i] == a[j] {
       return b[i] > b[j]
   }
   return a[i] > a[j]
}

func cmp3(x, y, t int) bool {
   return a[x] + b[x]*int64(t) <= a[y] + b[y]*int64(t)
}

func cmps(i, j, k int) bool {
   // compare cross: (a[j]-a[k])*(b[j]-b[i]) <= (b[k]-b[j])*(a[i]-a[j])
   t1 := (a[j]-a[k]) * (b[j]-b[i])
   t2 := (b[k]-b[j]) * (a[i]-a[j])
   return t1 <= t2
}

// build segment tree convex hulls
func build(p, d, L, R int) {
   // initialize indices
   for i := L; i <= R; i++ {
       idArr[i] = i
   }
   // sort by decreasing a, then decreasing b
   tmp := idArr[L : R+1]
   sort.Slice(tmp, func(i, j int) bool {
       return cmp(tmp[i], tmp[j])
   })
   tz[p] = L
   tp[p] = L - 1
   for _, idx := range tmp {
       if tz[p] <= tp[p] && b[idx] <= b[Tr[d][tp[p]]] {
           continue
       }
       for tz[p] < tp[p] && cmps(Tr[d][tp[p]-1], Tr[d][tp[p]], idx) {
           tp[p]--
       }
       tp[p]++
       Tr[d][tp[p]] = idx
   }
   if L >= R {
       return
   }
   mid := (L + R) >> 1
   build(p<<1, d+1, L, mid)
   build(p<<1|1, d+1, mid+1, R)
}

// query returns index of best line
func query(p, d, L, R, l, r, t int) int {
   if l > R || r < L {
       return -1
   }
   if L >= l && R <= r {
       for tz[p] < tp[p] && cmp3(Tr[d][tz[p]], Tr[d][tz[p]+1], t) {
           tz[p]++
       }
       return Tr[d][tz[p]]
   }
   mid := (L + R) >> 1
   x := query(p<<1, d+1, L, mid, l, r, t)
   y := query(p<<1|1, d+1, mid+1, R, l, r, t)
   if x < 0 || y < 0 {
       if x < 0 {
           return y
       }
       return x
   }
   if !cmp3(x, y, t) {
       return x
   }
   return y
}

func main() {
   defer writer.Flush()
   n := readInt()
   m := readInt()
   a = make([]int64, n)
   b = make([]int64, n)
   for i := 0; i < n; i++ {
       ai := readInt()
       bi := readInt()
       a[i] = int64(ai)
       b[i] = int64(bi)
   }
   qs := make([]Query, m)
   for i := 0; i < m; i++ {
       L := readInt()
       R := readInt()
       T := readInt()
       qs[i] = Query{i, L - 1, R - 1, T}
   }
   // sort queries by T
   sort.Slice(qs, func(i, j int) bool {
       return qs[i].T < qs[j].T
   })
   // allocate structures
   tz = make([]int, 4*n)
   tp = make([]int, 4*n)
   Tr = make([][]int, Mdeep)
   for i := 0; i < Mdeep; i++ {
       Tr[i] = make([]int, n)
   }
   idArr = make([]int, n)
   // build hulls
   build(1, 0, 0, n-1)
   ans := make([]int, m)
   for _, q := range qs {
       idx := query(1, 0, 0, n-1, q.L, q.R, q.T)
       ans[q.id] = idx + 1
   }
   // output
   for i := 0; i < m; i++ {
       writer.WriteString(strconv.Itoa(ans[i]))
       writer.WriteByte('\n')
   }
}
