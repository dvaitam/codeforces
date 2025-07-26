package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func ceilDiv(a, b int64) int64 {
   if b <= 0 {
       return 0
   }
   if a <= 0 {
       return 0
   }
   return (a + b - 1) / b
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   var n, m int
   var t int64
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   l := make([]int, n)
   r := make([]int, n)
   var totalL, totalR int64
   for i := 0; i < n; i++ {
       var li, ri int
       fmt.Fscan(reader, &li, &ri)
       l[i] = li
       r[i] = ri
       totalL += int64(li)
       totalR += int64(ri)
   }
   var q int
   fmt.Fscan(reader, &q)
   eqS := make([]int, m+2)
   for i := range eqS {
       eqS[i] = -1
   }
   // read known ranks
   knownMax := int64(0)
   var ps []int
   var ss []int64
   for i := 0; i < q; i++ {
       var p, s int64
       fmt.Fscan(reader, &p, &s)
       eqS[p] = int(s)
       if s > knownMax {
           knownMax = s
       }
       ps = append(ps, int(p))
       ss = append(ss, s)
   }
   fmt.Fscan(reader, &t)
   // total sum check
   if t < totalL || t > totalR {
       fmt.Fprintln(writer, "-1 -1")
       return
   }
   // prefix for equal constraints
   prefixCommon := make([]int, m+2)
   prefixOk := make([]bool, m+2)
   prefixCommon[0] = -1
   prefixOk[0] = true
   for i := 1; i <= m; i++ {
       prefixOk[i] = prefixOk[i-1]
       prefixCommon[i] = prefixCommon[i-1]
       if eqS[i] >= 0 {
           if prefixCommon[i-1] == -1 {
               prefixCommon[i] = eqS[i]
           } else if prefixCommon[i-1] != eqS[i] {
               prefixOk[i] = false
           }
       }
   }
   // suffix max for below group
   suffixMax := make([]int64, m+2)
   suffixMax[m] = -1
   for i := m - 1; i >= 0; i-- {
       cur := suffixMax[i+1]
       if eqS[i+1] >= 0 {
           if int64(eqS[i+1]) > cur {
               cur = int64(eqS[i+1])
           }
       }
       suffixMax[i] = cur
   }
   // sort l and r
   sort.Ints(r)
   sort.Ints(l)
   prefixR := make([]int64, n+1)
   prefixL := make([]int64, n+1)
   for i := 0; i < n; i++ {
       prefixR[i+1] = prefixR[i] + int64(r[i])
       prefixL[i+1] = prefixL[i] + int64(l[i])
   }
   // helper to compute R_top(k): sum min(r_i, k)
   Rtop := make([]int64, m+1)
   Rbot := make([]int64, m+1)
   Ltop := make([]int64, m+1)
   Lbot := make([]int64, m+1)
   for k := 0; k <= m; k++ {
       // Rtop[k]
       ki := k
       pos := sort.Search(n, func(i int) bool { return r[i] >= ki })
       Rtop[k] = prefixR[pos] + int64(n-pos)*int64(ki)
       // Rbot[k] as R_top(m-k)
       bi := m - k
       pos2 := sort.Search(n, func(i int) bool { return r[i] >= bi })
       Rbot[k] = prefixR[pos2] + int64(n-pos2)*int64(bi)
       // Ltop[k]: sum max(0, l_i - (m-k))
       th := m - k
       posl := sort.Search(n, func(i int) bool { return l[i] > th })
       cnt := n - posl
       sumVals := prefixL[n] - prefixL[posl]
       Ltop[k] = sumVals - int64(th)*int64(cnt)
       // Lbot[k]: sum max(0, l_i - k)
       th2 := k
       posl2 := sort.Search(n, func(i int) bool { return l[i] > th2 })
       cnt2 := n - posl2
       sumVals2 := prefixL[n] - prefixL[posl2]
       Lbot[k] = sumVals2 - int64(th2)*int64(cnt2)
   }
   // search best k
   bestK, bestS := 0, int64(-1)
   for k := m; k >= 1; k-- {
       if !prefixOk[k] {
           continue
       }
       // compute LB and UB for k
       LB := Ltop[k]
       if t-Rbot[k] > LB {
           LB = t - Rbot[k]
       }
       UB := Rtop[k]
       if t-Lbot[k] < UB {
           UB = t - Lbot[k]
       }
       if LB > UB {
           continue
       }
       // compute minimal S and maximal S
       lowS := ceilDiv(LB, int64(k))
       highS := UB / int64(k)
       if lowS < 0 {
           lowS = 0
       }
       if highS > int64(n) {
           highS = int64(n)
       }
       if lowS > highS {
           continue
       }
       // known constraints
       equal := prefixCommon[k]
       bottomMax := suffixMax[k]
       if equal >= 0 {
           S := int64(equal)
           if S < lowS || S > highS || S <= bottomMax {
               continue
           }
           bestK = k
           bestS = S
           break
       }
       // no equal constraint, S > bottomMax
       if bottomMax+1 > highS {
           continue
       }
       S := highS
       if S <= bottomMax {
           S = bottomMax + 1
       }
       if S < lowS {
           S = lowS
       }
       if S < lowS || S > highS {
           continue
       }
       bestK = k
       bestS = S
       break
   }
   if bestK == 0 {
       fmt.Fprintln(writer, "-1 -1")
   } else {
       fmt.Fprintf(writer, "%d %d", bestK, bestS)
   }
}
