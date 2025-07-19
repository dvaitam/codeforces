package main

import (
   "bufio"
   "fmt"
   "os"
)

const L = 7
const INF = 1 << 30

// state holds DP backtrack information
type state struct {
   i, j, k int
   r, ax, ay, az bool
}

func idx(b bool) int {
   if b {
       return 1
   }
   return 0
}

func min(a *int, b int) bool {
   if b < *a {
       *a = b
       return true
   }
   return false
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var a0, b0, c0 int
   _, err := fmt.Fscanf(reader, "%d+%d=%d", &a0, &b0, &c0)
   if err != nil {
       return
   }
   // digits 1..n, x[n] is least significant
   var x, y, z [L + 1]int
   sx := fmt.Sprintf("%d", a0)
   sy := fmt.Sprintf("%d", b0)
   sz := fmt.Sprintf("%d", c0)
   lx, ly, lz := len(sx), len(sy), len(sz)
   for i := 1; i <= lx; i++ {
       x[i] = int(sx[i-1] - '0')
   }
   for i := 1; i <= ly; i++ {
       y[i] = int(sy[i-1] - '0')
   }
   for i := 1; i <= lz; i++ {
       z[i] = int(sz[i-1] - '0')
   }
   // DP arrays
   var f [L+1][L+1][L+1][2][2][2][2]int
   var t [L+1][L+1][L+1][2][2][2][2]state
   // init base case: i=j=k=0, r=true
   for ax := 0; ax < 2; ax++ {
       for ay := 0; ay < 2; ay++ {
           for az := 0; az < 2; az++ {
               if ax|ay|az == 0 {
                   f[0][0][0][1][ax][ay][az] = INF
               } else {
                   f[0][0][0][1][ax][ay][az] = 1
               }
               t[0][0][0][1][ax][ay][az] = state{0,0,0,false,false,false,false}
           }
       }
   }
   // DP transitions
   for i := 0; i <= lx; i++ {
       for j := 0; j <= ly; j++ {
           for k := 0; k <= lz; k++ {
               if i|j|k == 0 {
                   continue
               }
               for r := 0; r < 2; r++ {
                   for ax := 0; ax < 2; ax++ {
                       for ay := 0; ay < 2; ay++ {
                           for az := 0; az < 2; az++ {
                               ff := &f[i][j][k][r][ax][ay][az]
                               *ff = INF
                               st := &t[i][j][k][r][ax][ay][az]
                               // match digit
                               if k > 0 && (x[i] + y[j] + r) % 10 == z[k] {
                                   for pax := 0; pax <= ax; pax++ {
                                       for pay := 0; pay <= ay; pay++ {
                                           for paz := 0; paz <= az; paz++ {
                                               pi, pj := i, j
                                               if i > 0 { pi = i-1 }
                                               if j > 0 { pj = j-1 }
                                               nr := 0
                                               if x[i] + y[j] + r > 9 {
                                                   nr = 1
                                               }
                                               cost := f[pi][pj][k-1][nr][pax][pay][paz]
                                               if i == 0 && pax == 1 { cost++ }
                                               if j == 0 && pay == 1 { cost++ }
                                               if min(ff, cost) {
                                                   *st = state{pi, pj, k-1, nr==1, pax==1, pay==1, paz==1}
                                               }
                                           }
                                       }
                                   }
                               }
                               // delete x
                               if ax == 1 {
                                   for pax := 0; pax <= ax; pax++ {
                                       for pay := 0; pay <= ay; pay++ {
                                           for paz := 0; paz <= az; paz++ {
                                               pj := j
                                               if j > 0 { pj = j-1 }
                                               kk := k
                                               if k > 0 { kk = k-1 }
                                               nr := 0
                                               if z[k] - y[j] - r < 0 {
                                                   nr = 1
                                               }
                                               cost := f[i][pj][kk][nr][pax][pay][paz] + 1
                                               if j == 0 && pay == 1 { cost++ }
                                               if k == 0 && paz == 1 { cost++ }
                                               if min(ff, cost) {
                                                   *st = state{i, pj, kk, nr==1, pax==1, pay==1, paz==1}
                                               }
                                           }
                                       }
                                   }
                               }
                               // delete y
                               if ay == 1 {
                                   for pax := 0; pax <= ax; pax++ {
                                       for pay := 0; pay <= ay; pay++ {
                                           for paz := 0; paz <= az; paz++ {
                                               pi := i
                                               if i > 0 { pi = i-1 }
                                               kk := k
                                               if k > 0 { kk = k-1 }
                                               nr := 0
                                               if z[k] - x[i] - r < 0 {
                                                   nr = 1
                                               }
                                               cost := f[pi][j][kk][nr][pax][pay][paz] + 1
                                               if i == 0 && pax == 1 { cost++ }
                                               if k == 0 && paz == 1 { cost++ }
                                               if min(ff, cost) {
                                                   *st = state{pi, j, kk, nr==1, pax==1, pay==1, paz==1}
                                               }
                                           }
                                       }
                                   }
                               }
                               // delete z
                               if az == 1 {
                                   for pax := 0; pax <= ax; pax++ {
                                       for pay := 0; pay <= ay; pay++ {
                                           for paz := 0; paz <= az; paz++ {
                                               pi := i
                                               if i > 0 { pi = i-1 }
                                               pj := j
                                               if j > 0 { pj = j-1 }
                                               nr := 0
                                               if x[i] + y[j] + r > 9 {
                                                   nr = 1
                                               }
                                               cost := f[pi][pj][k][nr][pax][pay][paz] + 1
                                               if i == 0 && pax == 1 { cost++ }
                                               if j == 0 && pay == 1 { cost++ }
                                               if min(ff, cost) {
                                                   *st = state{pi, pj, k, nr==1, pax==1, pay==1, paz==1}
                                               }
                                           }
                                       }
                                   }
                               }
                           }
                       }
                   }
               }
           }
       }
   }
   // backtrack result
   var A, B, C int64
   pw := int64(1)
   i, j, k := lx, ly, lz
   var r bool
   axb, ayb, azb := true, true, true
   for i > 0 || j > 0 || k > 0 || r {
       st := t[i][j][k][idx(r)][idx(axb)][idx(ayb)][idx(azb)]
       // no digits left => leftover carry
       if i == 0 && j == 0 && k == 0 {
           C += pw
       } else {
           // match
           if k > 0 && (x[i]+y[j]+idx(r))%10 == z[k] {
               i1 := i
               if i > 0 { i1 = i - 1 }
               j1 := j
               if j > 0 { j1 = j - 1 }
               k1 := k - 1
               nr := idx(x[i] + y[j] + idx(r) > 9)
               prev := state{i1, j1, k1, nr == 1, false, false, false}
               if prev == st {
                   A += int64(x[i]) * pw
                   B += int64(y[j]) * pw
                   C += int64(z[k]) * pw
                   i, j, k, r, axb, ayb, azb = prev.i, prev.j, prev.k, prev.r, prev.ax, prev.ay, prev.az
                   pw *= 10
                   continue
               }
           }
           // delete x
           if axb {
               j1 := j
               if j > 0 { j1 = j - 1 }
               k1 := k
               if k > 0 { k1 = k - 1 }
               nr := idx(z[k] - y[j] - idx(r) < 0)
               prev := state{i, j1, k1, nr == 1, false, false, false}
               if prev == st {
                   // a digit is computed from result
                   d := (z[k] - y[j] - idx(r)) % 10
                   if d < 0 { d += 10 }
                   A += int64(d) * pw
                   B += int64(y[j]) * pw
                   C += int64(z[k]) * pw
                   i, j, k, r, axb, ayb, azb = prev.i, prev.j, prev.k, prev.r, prev.ax, prev.ay, prev.az
                   pw *= 10
                   continue
               }
           }
           // delete y
           if ayb {
               i1 := i
               if i > 0 { i1 = i - 1 }
               k1 := k
               if k > 0 { k1 = k - 1 }
               nr := idx(z[k] - x[i] - idx(r) < 0)
               prev := state{i1, j, k1, nr == 1, false, false, false}
               if prev == st {
                   d := (z[k] - x[i] - idx(r)) % 10
                   if d < 0 { d += 10 }
                   A += int64(x[i]) * pw
                   B += int64(d) * pw
                   C += int64(z[k]) * pw
                   i, j, k, r, axb, ayb, azb = prev.i, prev.j, prev.k, prev.r, prev.ax, prev.ay, prev.az
                   pw *= 10
                   continue
               }
           }
           // delete z
           // default
           // match deletion of z branch
           {
               i1 := i
               if i > 0 { i1 = i - 1 }
               j1 := j
               if j > 0 { j1 = j - 1 }
               nr := idx(x[i] + y[j] + idx(r) > 9)
               prev := state{i1, j1, k, nr == 1, false, false, false}
               // use this branch
               A += int64(x[i]) * pw
               B += int64(y[j]) * pw
               C += int64((x[i] + y[j] + idx(r)) % 10) * pw
               i, j, k, r, axb, ayb, azb = prev.i, prev.j, prev.k, prev.r, prev.ax, prev.ay, prev.az
               pw *= 10
               continue
           }
       }
       // advance to next
       // for safety, break
       break
   }
   fmt.Printf("%d+%d=%d\n", A, B, C)
}
