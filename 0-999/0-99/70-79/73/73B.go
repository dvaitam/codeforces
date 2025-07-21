package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   type Racer struct {
       name string
       a    int
   }
   racers := make([]Racer, 0, n)
   var vasyaName string
   var vasyaA int
   for i := 0; i < n; i++ {
       var s string
       var ai int
       fmt.Fscan(in, &s, &ai)
       racers = append(racers, Racer{name: s, a: ai})
   }
   var m int
   fmt.Fscan(in, &m)
   b := make([]int, m)
   for i := 0; i < m; i++ {
       fmt.Fscan(in, &b[i])
   }
   fmt.Fscan(in, &vasyaName)
   // extract Vasya
   others := make([]Racer, 0, n-1)
   for _, r := range racers {
       if r.name == vasyaName {
           vasyaA = r.a
       } else {
           others = append(others, r)
       }
   }
   // build full point list for positions
   total := n
   // m awarded positions, zeros for rest
   B := make([]int, 0, total)
   for i := 0; i < m; i++ {
       B = append(B, b[i])
   }
   for i := 0; i < total-m; i++ {
       B = append(B, 0)
   }
   // sort B descending
   sort.Slice(B, func(i, j int) bool { return B[i] > B[j] })
   // best-case: Vasya takes max B[0]
   // remove B[0]
   bestB := make([]int, len(B)-1)
   copy(bestB, B[1:])
   threshBest := vasyaA + B[0]
   // worst-case: Vasya takes smallest (last)
   worstB := make([]int, len(B)-1)
   copy(worstB, B[:len(B)-1])
   threshWorst := vasyaA + B[len(B)-1]
   // compute best-case threats
   // assign bestB sorted desc to others sorted by a asc
   sort.Slice(others, func(i, j int) bool {
       if others[i].a != others[j].a {
           return others[i].a < others[j].a
       }
       return others[i].name < others[j].name
   })
   // bestB already sorted desc
   threatsBest := 0
   for i, r := range others {
       bi := bestB[i]
       sum := r.a + bi
       if sum > threshBest || (sum == threshBest && r.name < vasyaName) {
           threatsBest++
       }
   }
   // compute worst-case threats
   // assign worstB sorted desc to others sorted by a desc
   sort.Slice(others, func(i, j int) bool {
       if others[i].a != others[j].a {
           return others[i].a > others[j].a
       }
       return others[i].name > others[j].name
   })
   // worstB is sorted desc as prefix of B
   threatsWorst := 0
   for i, r := range others {
       bi := worstB[i]
       sum := r.a + bi
       if sum > threshWorst || (sum == threshWorst && r.name < vasyaName) {
           threatsWorst++
       }
   }
   // best rank is 1+threatsBest, worst is 1+threatsWorst
   out := bufio.NewWriter(os.Stdout)
   fmt.Fprintf(out, "%d %d", 1+threatsBest, 1+threatsWorst)
   out.Flush()
}
