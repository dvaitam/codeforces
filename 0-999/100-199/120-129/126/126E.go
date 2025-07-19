package main

import (
   "bufio"
   "fmt"
   "os"
)

const (
   le = 20
   ri = 21
)

var (
   U   = make(map[int64]struct{})
   a   [7][8]int
   P   [4][4]int
   L   [10]int
   R   [10]int
   A   [10]int
   Bv  [10]int
   T   = -1
   p   [22][22]int
   v   [7][8]bool
   w   [10][10]bool
   r2  [7][8]bool
   d2  [7][8]bool
   used [22]bool
   s   [13][15]byte
   S   [13][15]byte
   z   = []byte{'B','R','W','Y'}
   tAdd = 0
)

func add(x, y int) {
   P[x][y] = tAdd
   P[y][x] = tAdd
   L[tAdd] = x
   R[tAdd] = y
   tAdd++
}

func f2(x int) bool {
   if x == ri {
       return true
   }
   used[x] = true
   for i := ri; i >= 0; i-- {
       if p[x][i] > 0 && !used[i] && f2(i) {
           p[x][i]--
           p[i][x]++
           return true
       }
   }
   return false
}

func gg() {
   var c [10]int
   var D [10][10]int
   t := 0
   // hash Bv
   var h int64
   for i := 0; i < 10; i++ {
       h = h*1000000007 + int64(Bv[i])
   }
   if _, ok := U[h]; ok {
       return
   }
   U[h] = struct{}{}
   copy(c[:], A[:])
   // D init
   for i := 0; i < 10; i++ {
       D[i][i] = min(A[i], Bv[i])
       t += D[i][i] * 2
   }
   if 28 + t/2 <= T {
       return
   }
   // flow init
   for i := range p {
       for j := range p[i] {
           p[i][j] = 0
       }
   }
   for i := 0; i < 10; i++ {
       if A[i] > Bv[i] {
           p[le][i] = A[i] - Bv[i]
       } else if A[i] < Bv[i] {
           p[i+10][ri] = Bv[i] - A[i]
       }
   }
   for i := 0; i < 10; i++ {
       for j := 0; j < 10; j++ {
           if w[i][j] {
               u := min(p[le][i], p[j+10][ri])
               if u > 0 {
                   p[le][i] -= u
                   p[j+10][ri] -= u
                   p[i][j+10] = 28
                   t += u
                   D[i][j] += u
               }
           }
       }
   }
   // max flow augment
   for {
       for i := range used {
           used[i] = false
       }
       if f2(le) {
           t++
       } else {
           break
       }
   }
   if t <= T {
       return
   }
   T = t
   for i := 0; i < 10; i++ {
       for j := 0; j < 10; j++ {
           D[i][j] += p[j+10][i]
       }
   }
   // reset s
   for i := 0; i < 7; i++ {
       for j := 0; j < 8; j++ {
           s[i*2][j*2] = '.'
       }
   }
   // fill horizontal
   for i := 0; i < 7; i++ {
       for j := 0; j < 8; j++ {
           if r2[i][j] {
               o := P[a[i][j]][a[i][j+1]]
               for k := 0; k < 10; k++ {
                   if D[k][o] > 0 {
                       D[k][o]--
                       c[k]--
                       if a[i][j] == R[k] || a[i][j+1] == L[k] {
                           L[k], R[k] = R[k], L[k]
                       }
                       s[i*2][j*2] = z[L[k]]
                       s[i*2][j*2+2] = z[R[k]]
                       break
                   }
               }
           }
       }
   }
   // fill vertical
   for i := 0; i < 7; i++ {
       for j := 0; j < 8; j++ {
           if d2[i][j] {
               o := P[a[i][j]][a[i+1][j]]
               for k := 0; k < 10; k++ {
                   if D[k][o] > 0 {
                       D[k][o]--
                       c[k]--
                       if a[i][j] == R[k] || a[i+1][j] == L[k] {
                           L[k], R[k] = R[k], L[k]
                       }
                       s[i*2][j*2] = z[L[k]]
                       s[i*2+2][j*2] = z[R[k]]
                       break
                   }
               }
           }
       }
   }
   // fill remaining horizontal
   for i := 0; i < 7; i++ {
       for j := 0; j < 8; j++ {
           if s[i*2][j*2] == '.' && r2[i][j] {
               for k := 0; k < 10; k++ {
                   if c[k] > 0 {
                       c[k]--
                       s[i*2][j*2] = z[L[k]]
                       s[i*2][j*2+2] = z[R[k]]
                       break
                   }
               }
           }
       }
   }
   // fill remaining vertical
   for i := 0; i < 7; i++ {
       for j := 0; j < 8; j++ {
           if s[i*2][j*2] == '.' && d2[i][j] {
               for k := 0; k < 10; k++ {
                   if c[k] > 0 {
                       c[k]--
                       s[i*2][j*2] = z[L[k]]
                       s[i*2+2][j*2] = z[R[k]]
                       break
                   }
               }
           }
       }
   }
   // copy to S
   for i := 0; i < 13; i++ {
       for j := 0; j < 15; j++ {
           S[i][j] = s[i][j]
       }
   }
}

func ff(x, y int) {
   if x == 7 {
       gg()
       return
   }
   if y == 8 {
       ff(x+1, 0)
       return
   }
   if v[x][y] {
       ff(x, y+1)
       return
   }
   // horizontal
   if y < 7 && !v[x][y+1] {
       v[x][y], v[x][y+1] = true, true
       idx := P[a[x][y]][a[x][y+1]]
       Bv[idx]++
       s[x*2][y*2+1] = '-'
       r2[x][y] = true
       ff(x, y+1)
       v[x][y], v[x][y+1] = false, false
       Bv[idx]--
       s[x*2][y*2+1] = '.'
       r2[x][y] = false
   }
   // vertical
   if x < 6 && !v[x+1][y] {
       v[x][y], v[x+1][y] = true, true
       idx := P[a[x][y]][a[x+1][y]]
       Bv[idx]++
       s[x*2+1][y*2] = '|'
       d2[x][y] = true
       ff(x, y+1)
       v[x][y], v[x+1][y] = false, false
       Bv[idx]--
       s[x*2+1][y*2] = '.'
       d2[x][y] = false
   }
}

func min(a, b int) int {
   if a < b {
       return a
   }
   return b
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var row string
   for i := 0; i < 7; i++ {
       fmt.Fscan(in, &row)
       for j := 0; j < 8; j++ {
           switch row[j] {
           case 'B': a[i][j] = 0
           case 'R': a[i][j] = 1
           case 'W': a[i][j] = 2
           case 'Y': a[i][j] = 3
           }
       }
   }
   // define arcs
   add(0,3); add(0,2); add(0,1); add(0,0)
   add(1,3); add(1,2); add(1,1)
   add(2,3); add(2,2)
   add(3,3)
   // conflict matrix
   for i := 0; i < 10; i++ {
       for j := 0; j < 10; j++ {
           if L[i] == L[j] || L[i] == R[j] || R[i] == L[j] || R[i] == R[j] {
               w[i][j] = true
           }
       }
   }
   for i := 0; i < 10; i++ {
       fmt.Fscan(in, &A[i])
   }
   // init s grid
   for i := 0; i < 13; i++ {
       for j := 0; j < 15; j++ {
           s[i][j] = '.'
       }
   }
   ff(0, 0)
   // output
   wout := bufio.NewWriter(os.Stdout)
   defer wout.Flush()
   fmt.Fprintln(wout, T)
   for i := 0; i < 13; i++ {
       wout.Write(s[i][:15])
       wout.WriteByte('\n')
   }
}
