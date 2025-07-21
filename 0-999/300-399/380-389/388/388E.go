package main

import (
   "bufio"
   "fmt"
   "math/big"
   "os"
)

// gcd for int64
func gcd(a, b int64) int64 {
   if b == 0 {
       return a
   }
   return gcd(b, a%b)
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   fmt.Fscan(in, &n)
   // store line params
   numVx := make([]int64, n)
   numVy := make([]int64, n)
   numOx := make([]int64, n)
   numOy := make([]int64, n)
   den := make([]int64, n)
   for i := 0; i < n; i++ {
       var t1, x1, y1, t2, x2, y2 int64
       fmt.Fscan(in, &t1, &x1, &y1, &t2, &x2, &y2)
       di := t2 - t1
       den[i] = di
       numVx[i] = x2 - x1
       numVy[i] = y2 - y1
       // offset oi = (x1*t2 - x2*t1) / di, same for y
       numOx[i] = x1*t2 - x2*t1
       numOy[i] = y1*t2 - y2*t1
   }
   // adjacency matrix
   adj := make([][]bool, n)
   for i := range adj {
       adj[i] = make([]bool, n)
   }
   concurrencyMax := 1
   // big.Int for cross-multiplication
   var bi1, bi2 big.Int
   // temp map for t grouping
   type tKey struct{ num, den int64 }
   for i := 0; i < n; i++ {
       count := make(map[tKey]int)
       for j := i + 1; j < n; j++ {
           // dvx_num = numVx[i]*den[j] - numVx[j]*den[i]
           dvx := numVx[i]*den[j] - numVx[j]*den[i]
           dvy := numVy[i]*den[j] - numVy[j]*den[i]
           dox := numOx[j]*den[i] - numOx[i]*den[j]
           doy := numOy[j]*den[i] - numOy[i]*den[j]
           if dvx == 0 && dvy == 0 {
               continue
           }
           // check dvx*doy == dvy*dox
           bi1.SetInt64(dvx)
           bi1.Mul(&bi1, bi1.SetInt64(doy))
           bi2.SetInt64(dvy)
           bi2.Mul(&bi2, bi2.SetInt64(dox))
           if bi1.Cmp(&bi2) != 0 {
               continue
           }
           // they meet
           adj[i][j] = true
           adj[j][i] = true
           // compute meeting time t = dox/dvx or doy/dvy
           var tn, td int64
           if dvx != 0 {
               tn = dox
               td = dvx
           } else {
               tn = doy
               td = dvy
           }
           if td < 0 {
               tn = -tn
               td = -td
           }
           g := gcd(tn<0?-tn:tn, td)
           tn /= g
           td /= g
           key := tKey{tn, td}
           count[key]++
       }
       // update concurrency
       for _, c := range count {
           if c+1 > concurrencyMax {
               concurrencyMax = c + 1
           }
       }
   }
   // build bitsets for triangle detection
   if concurrencyMax < 3 {
       w := (n + 63) >> 6
       bits := make([][]uint64, n)
       for i := 0; i < n; i++ {
           bits[i] = make([]uint64, w)
           for j := 0; j < n; j++ {
               if adj[i][j] {
                   bits[i][j>>6] |= 1 << (uint(j) & 63)
               }
           }
       }
       // detect triangle
       foundTri := false
       for i := 0; i < n && !foundTri; i++ {
           for j := i + 1; j < n; j++ {
               if !adj[i][j] {
                   continue
               }
               // common neighbors
               for k := 0; k < w; k++ {
                   if bits[i][k]&bits[j][k] != 0 {
                       foundTri = true
                       break
                   }
               }
               if foundTri {
                   break
               }
           }
       }
       if foundTri && concurrencyMax < 3 {
           concurrencyMax = 3
       }
   }
   // if no meetings
   if concurrencyMax < 1 {
       concurrencyMax = 1
   }
   fmt.Println(concurrencyMax)
}
