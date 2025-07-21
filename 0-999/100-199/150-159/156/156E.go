package main

import (
   "bufio"
   "fmt"
   "os"
)

var primes = []int{
   2, 3, 5, 7, 11, 13, 17, 19, 23, 29,
   31, 37, 41, 43, 47, 53, 59, 61, 67, 71,
   73, 79, 83, 89, 97,
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n int
   fmt.Fscan(in, &n)
   a := make([]uint64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &a[i])
   }
   // precompute a mod primes
   pc := len(primes)
   aMod := make([][]int, pc)
   for pi := 0; pi < pc; pi++ {
       p := primes[pi]
       arr := make([]int, n)
       for j := 0; j < n; j++ {
           arr[j] = int(a[j] % uint64(p))
       }
       aMod[pi] = arr
   }
   var m int
   fmt.Fscan(in, &m)
   // queries data
   type Q struct { dIdx, LIdx int; s string; c uint64 }
   qs := make([]Q, m)
   // track lengths per base
   lengths := make([]map[int]int, 15) // map L to index
   for i := range lengths {
       lengths[i] = make(map[int]int)
   }
   // temp store raw lens
   rawL := make([]int, m)
   rawD := make([]int, m)
   for i := 0; i < m; i++ {
       var d int
       var s string
       var c uint64
       fmt.Fscan(in, &d, &s, &c)
       rawD[i] = d
       rawL[i] = len(s)
       idx := d - 2
       if _, ok := lengths[idx][len(s)]; !ok {
           lengths[idx][len(s)] = 0
       }
       qs[i] = Q{dIdx: idx, s: s, c: c}
   }
   // build length slices per base
   lengthList := make([][]int, 15)
   for dIdx := 0; dIdx < 15; dIdx++ {
       mp := lengths[dIdx]
       if len(mp) == 0 {
           continue
       }
       list := make([]int, 0, len(mp))
       for L := range mp {
           list = append(list, L)
       }
       // sort small Ls
       for i := 0; i < len(list); i++ {
           for j := i + 1; j < len(list); j++ {
               if list[i] > list[j] {
                   list[i], list[j] = list[j], list[i]
               }
           }
       }
       lengthList[dIdx] = list
       // map back to index
       for i, L := range list {
           lengths[dIdx][L] = i
       }
   }
   // precompute representations
   repDigits := make([][][]uint8, 15)
   repLens := make([][]int, 15)
   for dIdx := 0; dIdx < 15; dIdx++ {
       base := dIdx + 2
       repDigits[dIdx] = make([][]uint8, n)
       repLens[dIdx] = make([]int, n)
       for j := 0; j < n; j++ {
           x := j
           var digs []uint8
           if x == 0 {
               digs = []uint8{0}
           } else {
               for x > 0 {
                   digs = append(digs, uint8(x%base))
                   x /= base
               }
               // reverse
               for l, r := 0, len(digs)-1; l < r; l, r = l+1, r-1 {
                   digs[l], digs[r] = digs[r], digs[l]
               }
           }
           repDigits[dIdx][j] = digs
           repLens[dIdx][j] = len(digs)
       }
   }
   // build indices
   // dims: base dIdx -> LIdx -> pos -> val -> []j
   indices := make([][][][]([]int), 15)
   for dIdx := 0; dIdx < 15; dIdx++ {
       list := lengthList[dIdx]
       if len(list) == 0 {
           continue
       }
       base := dIdx + 2
       Lcnt := len(list)
       indices[dIdx] = make([][][][]int, Lcnt)
       for li, L := range list {
           posList := make([][][]int, L)
           for pos := 0; pos < L; pos++ {
               posList[pos] = make([][]int, base)
           }
           indices[dIdx][li] = posList
           for j := 0; j < n; j++ {
               rl := repLens[dIdx][j]
               offset := L - rl
               for pos := 0; pos < L; pos++ {
                   var dig int
                   if pos < offset {
                       dig = 0
                   } else {
                       dig = int(repDigits[dIdx][j][pos-offset])
                   }
                   indices[dIdx][li][pos][dig] = append(indices[dIdx][li][pos][dig], j)
               }
           }
       }
   }
   // total products per prime and base
   totalProd := make([][]int, pc)
   for pi := 0; pi < pc; pi++ {
       totalProd[pi] = make([]int, 15)
       p := primes[pi]
       for dIdx := 0; dIdx < 15; dIdx++ {
           prod := 1
           for j := 0; j < n; j++ {
               prod = (prod * aMod[pi][j]) % p
           }
           totalProd[pi][dIdx] = prod
       }
   }
   // process queries
   for i := 0; i < m; i++ {
       q := qs[i]
       dIdx := q.dIdx
       L := len(q.s)
       li := lengths[dIdx][L]
       // collect fixed
       type fv struct{ pos, val int }
       var fixed []fv
       for pos := 0; pos < L; pos++ {
           ch := q.s[pos]
           if ch != '?' {
               var v int
               if ch >= '0' && ch <= '9' {
                   v = int(ch - '0')
               } else {
                   v = int(ch - 'A' + 10)
               }
               fixed = append(fixed, fv{pos, v})
           }
       }
       ans := -1
       // iterate primes
       for pi := 0; pi < pc; pi++ {
           p := primes[pi]
           cmod := int(q.c % uint64(p))
           var prod int
           if len(fixed) == 0 {
               prod = totalProd[pi][dIdx]
           } else {
               // pick smallest list
               best := fixed[0]
               bestList := indices[dIdx][li][best.pos][best.val]
               for _, f := range fixed[1:] {
                   lst := indices[dIdx][li][f.pos][f.val]
                   if len(lst) < len(bestList) {
                       best = f
                       bestList = lst
                   }
               }
               prod = 1
               for _, j := range bestList {
                   ok := true
                   rl := repLens[dIdx][j]
                   offset := L - rl
                   for _, f := range fixed {
                       var dig int
                       if f.pos < offset {
                           dig = 0
                       } else {
                           dig = int(repDigits[dIdx][j][f.pos-offset])
                       }
                       if dig != f.val {
                           ok = false
                           break
                       }
                   }
                   if !ok {
                       continue
                   }
                   prod = (prod * aMod[pi][j]) % p
               }
           }
           if (prod + cmod) % p == 0 {
               ans = p
               break
           }
       }
       if ans <= 100 && ans != -1 {
           fmt.Fprintln(out, ans)
       } else {
           fmt.Fprintln(out, -1)
       }
   }
}
