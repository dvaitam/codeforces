package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

const N = 80

var n int
var L [N+9]int
var R [N+9]int
var v []int

var p [N+9]int
var pre, pn, cnt int
var inArr [N+9]float64
var smallArr [N+9]float64
var f [N+9][N+9][N+9]float64
var g [N+9][N+9]float64
var b [8][N+9][N+9]float64
var ans [N+9][N+9]float64

func backup(d int) {
   for i := 0; i < pn; i++ {
       for j := 0; i+j < pn; j++ {
           b[d][i][j] = g[i][j]
       }
   }
}

func recoverG(d int) {
   for i := 0; i < pn; i++ {
       for j := 0; i+j < pn; j++ {
           g[i][j] = b[d][i][j]
       }
   }
}

func ins(l, r int) {
   for i := l; i <= r; i++ {
       x := inArr[i]
       y := smallArr[i]
       for a := cnt; a >= 0; a-- {
           for bb := cnt - a; bb >= 0; bb-- {
               g[a+1][bb] += g[a][bb] * y
               g[a][bb+1] += g[a][bb] * x
               g[a][bb] *= 1 - x - y
           }
       }
       cnt++
   }
}

func dac(l, r, d int) {
   if l > r {
       return
   }
   if l == r {
       for i := 0; i < pn; i++ {
           for j := 0; i+j < pn; j++ {
               f[p[l]][i+pre][j+1] += g[i][j] * inArr[l]
           }
       }
       return
   }
   mid := (l + r) >> 1
   backup(d)
   ins(l, mid)
   dac(mid+1, r, d+1)
   recoverG(d)
   cnt -= mid - l + 1

   backup(d)
   ins(mid+1, r)
   dac(l, mid, d+1)
   recoverG(d)
   cnt -= r - mid
}

func work(lCoord, rCoord int) {
   pn = 0
   pre = 0
   for i := 1; i <= n; i++ {
       if L[i] <= lCoord && rCoord <= R[i] {
           pn++
           p[pn] = i
           length := float64(R[i] - L[i])
           inArr[pn] = float64(rCoord-lCoord) / length
           smallArr[pn] = float64(lCoord-L[i]) / length
       } else if R[i] <= lCoord {
           pre++
       }
   }
   for i := 0; i <= pn; i++ {
       for j := 0; j <= pn; j++ {
           g[i][j] = 0
       }
   }
   cnt = 0
   g[0][0] = 1.0
   dac(1, pn, 0)
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   fmt.Fscan(reader, &n)
   v = make([]int, 0, n*2)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &L[i], &R[i])
       v = append(v, L[i], R[i])
   }
   sort.Ints(v)
   uniq := v[:0]
   for i, x := range v {
       if i == 0 || x != v[i-1] {
           uniq = append(uniq, x)
       }
   }
   v = uniq
   for idx := 0; idx+1 < len(v); idx++ {
       work(v[idx], v[idx+1])
   }
   for i := 1; i <= n; i++ {
       for j := 0; j < n; j++ {
           for k := 1; k+j <= n; k++ {
               t := f[i][j][k] / float64(k)
               for r := 1; r <= k; r++ {
                   ans[i][j+r] += t
               }
           }
       }
   }
   for i := 1; i <= n; i++ {
       for j := 1; j <= n; j++ {
           fmt.Fprintf(writer, "%.10f", ans[i][j])
           if j < n {
               writer.WriteByte(' ')
           }
       }
       writer.WriteByte('\n')
   }
}
