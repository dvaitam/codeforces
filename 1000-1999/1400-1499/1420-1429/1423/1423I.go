package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   var K, Q int
   fmt.Fscan(in, &K, &Q)
   M := 1 << K
   L := make([]uint32, Q)
   R := make([]uint32, Q)
   V := make([]uint16, Q)
   for i := 0; i < Q; i++ {
       var l, r uint32
       var v uint16
       fmt.Fscan(in, &l, &r, &v)
       L[i], R[i], V[i] = l, r, v
   }
   A := make([]uint16, M)
   B := make([]uint16, M)
   for b := 0; b < 16; b++ {
       mask := uint16(1 << b)
       rowDiff := make([]int, M+1)
       colForce := make([]bool, M)
       colAll := false
       // ones constraints
       for i := 0; i < Q; i++ {
           if V[i]&mask != 0 {
               hL := int(L[i] >> K)
               hR := int(R[i] >> K)
               rowDiff[hL]++
               rowDiff[hR+1]--
               if hL == hR {
                   lmod := int(L[i] & uint32(M-1))
                   rmod := int(R[i] & uint32(M-1))
                   for j := lmod; j <= rmod; j++ {
                       colForce[j] = true
                   }
               } else {
                   colAll = true
               }
           }
       }
       // build rowForce and prefix sum
       rowForce := make([]bool, M)
       rowPS := make([]int, M)
       cnt := 0
       for i := 0; i < M; i++ {
           cnt += rowDiff[i]
           if cnt > 0 {
               rowForce[i] = true
               rowPS[i] = 1
           }
           if i > 0 {
               rowPS[i] += rowPS[i-1]
           }
       }
       // finalize columns
       if colAll {
           for j := 0; j < M; j++ {
               colForce[j] = true
           }
       }
       // prefix sum for cols
       colPS := make([]int, M)
       if colForce[0] {
           colPS[0] = 1
       }
       for i := 1; i < M; i++ {
           colPS[i] = colPS[i-1]
           if colForce[i] {
               colPS[i]++
           }
       }
       totalCols := colPS[M-1]
       // zeros constraints
       for i := 0; i < Q; i++ {
           if V[i]&mask == 0 {
               hL := int(L[i] >> K)
               hR := int(R[i] >> K)
               lmod := int(L[i] & uint32(M-1))
               rmod := int(R[i] & uint32(M-1))
               // interior rows full columns
               if hR >= hL+2 {
                   if rowPS[hR-1]-rowPS[hL] > 0 && totalCols > 0 {
                       fmt.Fprintln(out, "impossible")
                       return
                   }
               }
               // head row
               endj := rmod
               if hL < hR {
                   endj = M - 1
               }
               if rowForce[hL] {
                   before := 0
                   if lmod > 0 {
                       before = colPS[lmod-1]
                   }
                   if colPS[endj]-before > 0 {
                       fmt.Fprintln(out, "impossible")
                       return
                   }
               }
               // tail row
               if hR != hL && rowForce[hR] {
                   if colPS[rmod] > 0 {
                       fmt.Fprintln(out, "impossible")
                       return
                   }
               }
           }
       }
       // assign bits
       for j := 0; j < M; j++ {
           if colForce[j] {
               A[j] |= mask
           }
       }
       for h := 0; h < M; h++ {
           if rowForce[h] {
               B[h] |= mask
           }
       }
   }
   // output
   fmt.Fprintln(out, "possible")
   for j := 0; j < M; j++ {
       fmt.Fprintln(out, A[j])
   }
   for h := 0; h < M; h++ {
       fmt.Fprintln(out, B[h])
   }
}
