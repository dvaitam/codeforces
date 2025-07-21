package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, d, S int64
   if _, err := fmt.Fscan(in, &n, &d, &S); err != nil {
       return
   }
   type robot struct{ f, c int64 }
   var epos []robot
   var negCount int64
   var totalCposE int64
   for i := int64(0); i < n; i++ {
       var ci, fi, li int64
       fmt.Fscan(in, &ci, &fi, &li)
       if li >= d && ci > 0 {
           epos = append(epos, robot{f: fi, c: ci})
           totalCposE += ci
       } else {
           negCount++
       }
   }
   A := int64(len(epos))
   if A == 0 {
       fmt.Println("0 0")
       return
   }
   // B = totalCposE - A
   B := totalCposE - A
   Cneg := negCount
   // sort by fuel ascending
   sort.Slice(epos, func(i, j int) bool { return epos[i].f < epos[j].f })
   // prefix fuel
   costF := make([]int64, A+1)
   for i := int64(1); i <= A; i++ {
       costF[i] = costF[i-1] + epos[i-1].f
   }
   // find max r with costF[r] <= S
   rMax := int64(0)
   for i := int64(1); i <= A; i++ {
       if costF[i] <= S {
           rMax = i
       } else {
           break
       }
   }
   if rMax == 0 {
       fmt.Println("0 0")
       return
   }
   // compute kMax
   var kMax int64
   if B + rMax >= Cneg {
       // can load all negatives
       kMax = A + Cneg
   } else {
       kMax = A + B + rMax
   }
   // determine minimal r needed for kMax
   var rNeed int64
   if kMax == A + Cneg {
       // can reach full
       need := Cneg - B
       if need < 1 {
           need = 1
       }
       rNeed = need
   } else {
       rNeed = rMax
   }
   if rNeed > rMax {
       rNeed = rMax
   }
   fuel := costF[rNeed]
   fmt.Printf("%d %d", kMax, fuel)
}
