package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
   "sort"
   "time"
)

const V = 200003
const segN = 1 << 16

var (
   seed      uint64
   vec       = make([][]int, V)
   bucketVer = make([]int, V)
   version   = 1
   x, y, xx  []int
   n, L, k   int
   B         int
   segSum    [2 * segN]int
   segPre    [2 * segN]int
)

func hash(cx, cy int) int {
   h := ((uint64(cx) << 32) | uint64(cy)) ^ seed
   return int(h % V)
}

func ensure(h int) {
   if bucketVer[h] != version {
       vec[h] = vec[h][:0]
       bucketVer[h] = version
   }
}

func build(l, r int, R float64) {
   if B != 0 {
       for i := l; i <= r; i++ {
           h := xx[i]
           ensure(h)
           vec[h] = vec[h][:0]
       }
   }
   B = int(R*2) + 1
   for i := l; i <= r; i++ {
       cx := x[i] / B
       cy := y[i] / B
       h := hash(cx, cy)
       ensure(h)
       xx[i] = h
       vec[h] = append(vec[h], i)
   }
}

func remove(p int) {
   h := xx[p]
   ensure(h)
   vec[h] = vec[h][1:]
}

func add(p int) {
   cx := x[p] / B
   cy := y[p] / B
   h := hash(cx, cy)
   ensure(h)
   xx[p] = h
   vec[h] = append(vec[h], p)
}

func pup(p int) {
   left, right := p<<1, p<<1|1
   segSum[p] = segSum[left] + segSum[right]
   if segPre[left] > segSum[left]+segPre[right] {
       segPre[p] = segPre[left]
   } else {
       segPre[p] = segSum[left] + segPre[right]
   }
}

func upd(q, v int) {
   p := q + segN
   segSum[p] += v
   if segSum[p] > 0 {
       segPre[p] = segSum[p]
   } else {
       segPre[p] = 0
   }
   for p >>= 1; p > 0; p >>= 1 {
       pup(p)
   }
}

func clr(q int) {
   p := q + segN
   for p > 0 {
       segSum[p] = 0
       segPre[p] = 0
       p >>= 1
   }
}

type event struct{
   ang float64
   id  int
}

func check(p int, R float64) bool {
   // clear previous tree segments are cleaned by clr
   var evs []event
   cx := x[p] / B
   cy := y[p] / B
   // mdf: update range [q, q+L)
   mdf := func(q, v int) {
       upd(q, v)
       if q+L <= n {
           upd(q+L, -v)
       }
   }
   for dx := -1; dx <= 1; dx++ {
       for dy := -1; dy <= 1; dy++ {
           h0 := hash(cx+dx, cy+dy)
           ensure(h0)
           for _, j := range vec[h0] {
               if j != p && abs(j-p) < L {
                   dx0 := float64(x[p]-x[j])
                   dy0 := float64(y[p]-y[j])
                   dist := math.Hypot(dx0, dy0) / (2 * R)
                   if dist >= 1 {
                       continue
                   }
                   phi := math.Atan2(dy0, dx0)
                   dphi := math.Acos(dist)
                   argl := phi - dphi
                   argr := phi + dphi
                   if argl < -math.Pi || argr > math.Pi {
                       mdf(j, 1)
                   }
                   if argl < -math.Pi {
                       argl += 2 * math.Pi
                   }
                   if argr > math.Pi {
                       argr -= 2 * math.Pi
                   }
                   evs = append(evs, event{argl, j})
                   evs = append(evs, event{argr, -j})
               }
           }
       }
   }
   sort.Slice(evs, func(i, j int) bool { return evs[i].ang < evs[j].ang })
   flg := segPre[1] >= k-1
   for _, e := range evs {
       if flg {
           break
       }
       id := e.id
       if id > 0 {
           mdf(id, 1)
       } else {
           mdf(-id, -1)
       }
       if segPre[1] >= k-1 {
           flg = true
       }
   }
   for _, e := range evs {
       j := abs(e.id)
       clr(j)
       if j+L <= n {
           clr(j + L)
       }
   }
   return flg
}

func solve() {
   version++
   B = 0
   fmt.Fscan(in, &n, &L, &k)
   R := 225675834.0 * math.Sqrt(float64(k-1)/float64(L))
   x = make([]int, n+2)
   y = make([]int, n+2)
   xx = make([]int, n+2)
   for i := 1; i <= n; i++ {
       fmt.Fscan(in, &x[i], &y[i])
       x[i] += 100000000
       y[i] += 100000000
   }
   build(1, min(L, n), R)
   for i := 1; i <= n; i++ {
       if i > L {
           remove(i - L)
       }
       if check(i, R) {
           l, r := 0.0, R
           for r-l >= l*1e-10 {
               md := (l + r) / 2
               if check(i, md) {
                   r = md
               } else {
                   l = md
               }
           }
           R = l
           l0 := i - L + 1
           if l0 < 1 {
               l0 = 1
           }
           r0 := i + L - 1
           if r0 > n {
               r0 = n
           }
           build(l0, r0, R)
       }
       if i+L <= n {
           add(i + L)
       }
   }
   // clear remaining
   for i := max(1, n-L+1); i <= n; i++ {
       h := xx[i]
       ensure(h)
       vec[h] = vec[h][:0]
   }
   fmt.Printf("%.10f\n", R)
}

var in = bufio.NewReader(os.Stdin)

func main() {
   seed = uint64(time.Now().UnixNano())
   var T int
   fmt.Fscan(in, &T)
   for T > 0 {
       solve()
       T--
   }
}

func abs(a int) int {
   if a < 0 {
       return -a
   }
   return a
}

func min(a, b int) int {
   if a < b {
       return a
   }
   return b
}

func max(a, b int) int {
   if a > b {
       return a
   }
   return b
}
