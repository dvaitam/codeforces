package main

import (
   "bufio"
   "fmt"
   "math/rand"
   "os"
   "sort"
   "time"
)

var (
   rd  = bufio.NewReader(os.Stdin)
   wr  = bufio.NewWriter(os.Stdout)
)

func readInt() int {
   var x int
   var c byte
   buf := make([]byte, 0)
   for {
       b, err := rd.ReadByte()
       if err != nil {
           break
       }
       c = b
       if c >= '0' && c <= '9' {
           buf = append(buf, c)
           break
       }
   }
   for {
       b, err := rd.ReadByte()
       if err != nil || b < '0' || b > '9' {
           break
       }
       buf = append(buf, b)
   }
   for _, d := range buf {
       x = x*10 + int(d-'0')
   }
   return x
}

type Edge struct{ c, v int }

func main() {
   defer wr.Flush()
   n := readInt()
   m := readInt()
   nm := n * m
   // random hashes
   p := make([]uint64, nm+1)
   rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
   for i := 1; i <= nm; i++ {
       p[i] = rnd.Uint64() >> 1
   }
   size := n + m
   A := make([][]Edge, size+1)
   B := make([][]Edge, size+1)
   hA := make([]uint64, size+1)
   hB := make([]uint64, size+1)
   // read matrix A
   for i := 1; i <= n; i++ {
       for j := 1; j <= m; j++ {
           c := readInt()
           A[i] = append(A[i], Edge{c, j + n})
           A[j+n] = append(A[j+n], Edge{c, i})
           hA[i] ^= p[c]
           hA[j+n] ^= p[c]
       }
   }
   // sort adjacency
   for i := 1; i <= size; i++ {
       ai := A[i]
       sort.Slice(ai, func(a, b int) bool {
           if ai[a].c != ai[b].c {
               return ai[a].c < ai[b].c
           }
           return ai[a].v < ai[b].v
       })
       A[i] = ai
   }
   // read matrix B
   for i := 1; i <= n; i++ {
       for j := 1; j <= m; j++ {
           c := readInt()
           B[i] = append(B[i], Edge{c, j + n})
           B[j+n] = append(B[j+n], Edge{c, i})
           hB[i] ^= p[c]
           hB[j+n] ^= p[c]
       }
   }
   for i := 1; i <= size; i++ {
       bi := B[i]
       sort.Slice(bi, func(a, b int) bool {
           if bi[a].c != bi[b].c {
               return bi[a].c < bi[b].c
           }
           return bi[a].v < bi[b].v
       })
       B[i] = bi
   }
   // matchings
   wA := make([]int, size+1)
   wB := make([]int, size+1)
   st := make([]int, size+1)
   var tstk int

   var dfs func(x, y int) bool
   dfs = func(x, y int) bool {
       if hA[x] != hB[y] {
           return false
       }
       if wA[x] == y && wB[y] == x {
           return true
       }
       if wA[x] != 0 || wB[y] != 0 {
           return false
       }
       wA[x], wB[y] = y, x
       st[tstk] = x
       tstk++
       la := len(A[x])
       for i := 0; i < la; i++ {
           if !dfs(A[x][i].v, B[y][i].v) {
               return false
           }
       }
       return true
   }
   // initial matching by hash
   l, r := 1, 0
   if n < m {
       l, r = 1, n
   } else {
       l, r = n+1, n+m
   }
   for i := l; i <= r; i++ {
       if wA[i] != 0 {
           continue
       }
       flag := false
       for j := l; j <= r; j++ {
           if wB[j] == 0 && hA[i] == hB[j] {
               start := tstk
               if dfs(i, j) {
                   flag = true
                   break
               }
               // rollback
               for tstk > start {
                   tstk--
                   x := st[tstk]
                   y := wA[x]
                   wA[x], wB[y] = 0, 0
               }
           }
       }
       if !flag {
           fmt.Fprintln(wr, -1)
           return
       }
   }
   // match remaining arbitrarily
   jj := 1
   for i := 1; i <= size; i++ {
       if wA[i] == 0 {
           for wB[jj] != 0 {
               jj++
           }
           wA[i], wB[jj] = jj, i
       }
   }
   // generate operations
   type Op struct{ p, x, y int }
   ops := make([]Op, 0, size)
   // rows
   for i := 1; i <= n; i++ {
       for wA[i] != 0 && wA[i] < i {
           ops = append(ops, Op{1, i, wA[i]})
           // swap in wA mapping
           wi := wA[i]
           wA[i], wA[wi] = wA[wi], wA[i]
       }
   }
   // columns
   for i := n + 1; i <= size; i++ {
       for wA[i] != 0 && wA[i] < i {
           ops = append(ops, Op{2, i - n, wA[i] - n})
           wi := wA[i]
           wA[i], wA[wi] = wA[wi], wA[i]
       }
   }
   // output
   fmt.Fprintln(wr, len(ops))
   for _, op := range ops {
       fmt.Fprintln(wr, op.p, op.x, op.y)
   }
}
