package main

import (
   "bufio"
   "os"
   "strconv"
)

const INF = int64(1e18)

type Query struct { l, r, id int }

var (
   n, q    int
   p       []int
   w, sum  []int64
   son     [][]int
   dfnr    []int
   ia      []int64
   lz, mn  []int64
   qr      [][]Query
   ans     []int64
   reader  = bufio.NewReader(os.Stdin)
   writer  = bufio.NewWriter(os.Stdout)
)

func readInt() int {
   var x int
   var neg bool
   c, err := reader.ReadByte()
   for err == nil && (c < '0' || c > '9') && c != '-' {
       c, err = reader.ReadByte()
   }
   if c == '-' {
       neg = true
       c, _ = reader.ReadByte()
   }
   for err == nil && c >= '0' && c <= '9' {
       x = x*10 + int(c-'0')
       c, err = reader.ReadByte()
   }
   if neg {
       return -x
   }
   return x
}

func readInt64() int64 {
   var x int64
   var neg bool
   c, err := reader.ReadByte()
   for err == nil && (c < '0' || c > '9') && c != '-' {
       c, err = reader.ReadByte()
   }
   if c == '-' {
       neg = true
       c, _ = reader.ReadByte()
   }
   for err == nil && c >= '0' && c <= '9' {
       x = x*10 + int64(c-'0')
       c, err = reader.ReadByte()
   }
   if neg {
       return -x
   }
   return x
}

func buildSeg(l, r, pidx int) {
   lz[pidx] = 0
   if l == r {
       mn[pidx] = ia[l]
       return
   }
   mid := (l + r) >> 1
   buildSeg(l, mid, pidx<<1)
   buildSeg(mid+1, r, pidx<<1|1)
   if mn[pidx<<1] < mn[pidx<<1|1] {
       mn[pidx] = mn[pidx<<1]
   } else {
       mn[pidx] = mn[pidx<<1|1]
   }
}

func addSeg(L, R int, v int64, l, r, pidx int) {
   if L <= l && r <= R {
       lz[pidx] += v
       mn[pidx] += v
       return
   }
   if r < L || R < l {
       return
   }
   // pushdown
   if lz[pidx] != 0 {
       lc, rc := pidx<<1, pidx<<1|1
       lz[lc] += lz[pidx]
       mn[lc] += lz[pidx]
       lz[rc] += lz[pidx]
       mn[rc] += lz[pidx]
       lz[pidx] = 0
   }
   mid := (l + r) >> 1
   if L <= mid {
       addSeg(L, R, v, l, mid, pidx<<1)
   }
   if R > mid {
       addSeg(L, R, v, mid+1, r, pidx<<1|1)
   }
   lc, rc := pidx<<1, pidx<<1|1
   if mn[lc] < mn[rc] {
       mn[pidx] = mn[lc]
   } else {
       mn[pidx] = mn[rc]
   }
}

func querySeg(L, R, l, r, pidx int) int64 {
   if L <= l && r <= R {
       return mn[pidx]
   }
   if r < L || R < l {
       return INF
   }
   // pushdown
   if lz[pidx] != 0 {
       lc, rc := pidx<<1, pidx<<1|1
       lz[lc] += lz[pidx]
       mn[lc] += lz[pidx]
       lz[rc] += lz[pidx]
       mn[rc] += lz[pidx]
       lz[pidx] = 0
   }
   mid := (l + r) >> 1
   res := INF
   if L <= mid {
       v := querySeg(L, R, l, mid, pidx<<1)
       if v < res {
           res = v
       }
   }
   if R > mid {
       v := querySeg(L, R, mid+1, r, pidx<<1|1)
       if v < res {
           res = v
       }
   }
   return res
}

func pdfs(x int) {
   if len(son[x]) == 0 {
       ia[x] = sum[x]
   } else {
       ia[x] = INF
   }
   for _, u := range son[x] {
       sum[u] = sum[x] + w[u]
       pdfs(u)
   }
}

func dfs(x int) {
   for _, qu := range qr[x] {
       res := querySeg(qu.l, qu.r, 1, n, 1)
       ans[qu.id] = sum[x] + res
   }
   for _, u := range son[x] {
       addSeg(u, dfnr[u], -2*w[u], 1, n, 1)
       dfs(u)
       addSeg(u, dfnr[u], 2*w[u], 1, n, 1)
   }
}

func main() {
   defer writer.Flush()
   n = readInt()
   q = readInt()
   p = make([]int, n+1)
   w = make([]int64, n+1)
   son = make([][]int, n+1)
   for i := 2; i <= n; i++ {
       p[i] = readInt()
       w[i] = readInt64()
       son[p[i]] = append(son[p[i]], i)
   }
   sum = make([]int64, n+1)
   ia = make([]int64, n+2)
   pdfs(1)
   // build segment tree
   lz = make([]int64, (n+2)*4)
   mn = make([]int64, (n+2)*4)
   buildSeg(1, n, 1)
   // dfnr
   dfnr = make([]int, n+1)
   for i := n; i >= 2; i-- {
       if dfnr[i] < i {
           dfnr[i] = i
       }
       if dfnr[p[i]] < dfnr[i] {
           dfnr[p[i]] = dfnr[i]
       }
   }
   // queries
   qr = make([][]Query, n+1)
   ans = make([]int64, q+1)
   for i := 1; i <= q; i++ {
       x := readInt()
       l := readInt()
       r := readInt()
       qr[x] = append(qr[x], Query{l, r, i})
   }
   dfs(1)
   // output
   for i := 1; i <= q; i++ {
       writer.WriteString(strconv.FormatInt(ans[i], 10))
       writer.WriteByte('\n')
   }
}
