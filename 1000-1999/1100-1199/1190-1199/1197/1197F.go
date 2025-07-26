package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

const mod = 998244353

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   fmt.Fscan(reader, &n)
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   var m int
   fmt.Fscan(reader, &m)
   fixed := make([]map[int]int, n)
   for i := range fixed {
       fixed[i] = make(map[int]int)
   }
   for i := 0; i < m; i++ {
       var xi, yi, ci int
       fmt.Fscan(reader, &xi, &yi, &ci)
       fixed[xi-1][yi] = ci - 1
   }
   // read f matrix
   fMat := [3][4]bool{}
   for i := 0; i < 3; i++ {
       for j := 1; j <= 3; j++ {
           var v int
           fmt.Fscan(reader, &v)
           fMat[i][j] = (v == 1)
       }
   }
   // state: idx 0..63 represents (g0,g1,g2) = (idx>>4, (idx>>2)&3, idx&3)
   // precompute transitions for x>=4
   var transC [3][64]int
   for c := 0; c < 3; c++ {
       for idx := 0; idx < 64; idx++ {
           prev0 := idx >> 4
           prev1 := (idx >> 2) & 3
           prev2 := idx & 3
           used := [5]bool{}
           if fMat[c][1] {
               used[prev0] = true
           }
           if fMat[c][2] {
               used[prev1] = true
           }
           if fMat[c][3] {
               used[prev2] = true
           }
           mex := 0
           for used[mex] {
               mex++
           }
           newIdx := (mex<<4) | (prev0<<2) | prev1
           transC[c][idx] = newIdx
       }
   }
   // build base matrix M
   const S = 64
   // mpow[k] is M^(1<<k)
   mpow := make([][S][S]int, 31)
   // base M
   for i := 0; i < S; i++ {
       for c := 0; c < 3; c++ {
           j := transC[c][i]
           mpow[0][i][j] = (mpow[0][i][j] + 1) % mod
       }
   }
   // matrix exponentiation
   for k := 1; k < 31; k++ {
       for i := 0; i < S; i++ {
           for j := 0; j < S; j++ {
               sum := 0
               for t := 0; t < S; t++ {
                   sum += mpow[k-1][i][t] * mpow[k-1][t][j]
                   if sum >= mod*mod {
                       sum -= mod * mod
                   }
               }
               mpow[k][i][j] = sum % mod
           }
       }
   }
   // XOR DP across strips
   dpXor := [4]int{1, 0, 0, 0}
   for i := 0; i < n; i++ {
       cnt := processStrip(a[i], fixed[i], fMat, transC, mpow)
       var newDp [4]int
       for x := 0; x < 4; x++ {
           sum := 0
           for g := 0; g < 4; g++ {
               sum = (sum + dpXor[x^g]*cnt[g]) % mod
           }
           newDp[x] = sum
       }
       dpXor = newDp
   }
   fmt.Println(dpXor[0])
}

// process one strip length L with fixed colors
func processStrip(L int, fix map[int]int, fMat [3][4]bool, transC [3][64]int, mpow [][64][64]int) [4]int {
   const S = 64
   // dpVec holds counts for state idx
   dpVec := make([]int, S)
   // initial state at x=0: (0,0,0)
   dpVec[0] = 1
   // initial small x: 1..min(L,3)
   endInit := L
   if endInit > 3 {
       endInit = 3
   }
   for x := 1; x <= endInit; x++ {
       newVec := make([]int, S)
       for idx, val := range dpVec {
           if val == 0 {
               continue
           }
           prev0 := idx >> 4
           prev1 := (idx >> 2) & 3
           prev2 := idx & 3
           if c0, ok := fix[x]; ok {
               // single color
               used := [5]bool{}
               if fMat[c0][1] && x-1 >= 1 {
                   used[prev0] = true
               }
               if fMat[c0][2] && x-2 >= 1 {
                   used[prev1] = true
               }
               if fMat[c0][3] && x-3 >= 1 {
                   used[prev2] = true
               }
               mex := 0
               for used[mex] {
                   mex++
               }
               ni := (mex<<4)|(prev0<<2)|prev1
               newVec[ni] = (newVec[ni] + val) % mod
           } else {
               for c := 0; c < 3; c++ {
                   used := [5]bool{}
                   if fMat[c][1] && x-1 >= 1 {
                       used[prev0] = true
                   }
                   if fMat[c][2] && x-2 >= 1 {
                       used[prev1] = true
                   }
                   if fMat[c][3] && x-3 >= 1 {
                       used[prev2] = true
                   }
                   mex := 0
                   for used[mex] {
                       mex++
                   }
                   ni := (mex<<4)|(prev0<<2)|prev1
                   newVec[ni] = (newVec[ni] + val) % mod
               }
           }
       }
       dpVec = newVec
   }
   if L <= 3 {
       var cnt [4]int
       for idx, val := range dpVec {
           if val == 0 {
               continue
           }
           g := idx >> 4
           cnt[g] = (cnt[g] + val) % mod
       }
       return cnt
   }
   // x >= 4
   cur := 3
   // collect fixed positions >=4
   posList := make([]int, 0, len(fix))
   for pos := range fix {
       if pos >= 4 && pos <= L {
           posList = append(posList, pos)
       }
   }
   sort.Ints(posList)
   // process segments
   for _, pos := range posList {
       seg := pos - cur - 1
       if seg > 0 {
           dpVec = applyMatrix(dpVec, seg, mpow)
       }
       // apply fixed color at pos
       c0 := fix[pos]
       newVec := make([]int, S)
       for idx, val := range dpVec {
           if val == 0 {
               continue
           }
           ni := transC[c0][idx]
           newVec[ni] = (newVec[ni] + val) % mod
       }
       dpVec = newVec
       cur = pos
   }
   // final segment to L
   if L > cur {
       dpVec = applyMatrix(dpVec, L-cur, mpow)
   }
   var cnt [4]int
   for idx, val := range dpVec {
       if val == 0 {
           continue
       }
       g := idx >> 4
       cnt[g] = (cnt[g] + val) % mod
   }
   return cnt
}

// applyMatrix multiplies vector by M^seg (using mpow table)
func applyMatrix(vec []int, seg int, mpow [][64][64]int) []int {
   const S = 64
   res := make([]int, S)
   // copy vec to res as initial
   copy(res, vec)
   for k := 0; seg > 0; k++ {
       if seg&1 == 1 {
           tmp := make([]int, S)
           mat := mpow[k]
           for i := 0; i < S; i++ {
               if res[i] == 0 {
                   continue
               }
               vi := res[i]
               for j := 0; j < S; j++ {
                   tmp[j] = (tmp[j] + vi*mat[i][j]) % mod
               }
           }
           res = tmp
       }
       seg >>= 1
   }
   return res
}
