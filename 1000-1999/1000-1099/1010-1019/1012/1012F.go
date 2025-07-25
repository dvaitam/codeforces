package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

type Trip struct {
   idx    int
   s, e   int64
   t, L   int64
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var N, P int
   if _, err := fmt.Fscan(in, &N, &P); err != nil {
       return
   }
   trips := make([]Trip, N)
   for i := 0; i < N; i++ {
       var s, length, t int64
       fmt.Fscan(in, &s, &length, &t)
       trips[i] = Trip{idx: i, s: s, e: s + length - 1, t: t, L: s - t - 1}
   }
   // sort by start for gaps
   tripsByS := make([]Trip, N)
   copy(tripsByS, trips)
   sort.Slice(tripsByS, func(i, j int) bool { return tripsByS[i].s < tripsByS[j].s })
   // compute max L
   var maxL int64 = 0
   for _, tr := range trips {
       if tr.L > maxL {
           maxL = tr.L
       }
   }
   // build allowed application gaps
   type gap struct{ l, r int64 }
   var gaps []gap
   prevEnd := int64(0)
   for _, tr := range tripsByS {
       if prevEnd+1 <= tr.s-1 {
           gaps = append(gaps, gap{l: prevEnd + 1, r: tr.s - 1})
       }
       prevEnd = tr.e
   }
   if prevEnd+1 <= maxL {
       gaps = append(gaps, gap{l: prevEnd + 1, r: maxL})
   }
   // sort jobs by deadline L increasing
   jobs := make([]Trip, N)
   copy(jobs, trips)
   sort.Slice(jobs, func(i, j int) bool { return jobs[i].L < jobs[j].L })
   // passport schedules: list of intervals [start, finish]
   schedules := make([][]gap, P)
   // assignment results
   assign := make([]int, N)
   applyDay := make([]int64, N)
   // function to find latest allowed start in [a, ub]
   findLatest := func(a, ub int64) (int64, bool) {
       // binary search first gap with l > ub
       i := sort.Search(len(gaps), func(i int) bool { return gaps[i].l > ub })
       for k := i - 1; k >= 0; k-- {
           g := gaps[k]
           if g.r < a {
               break
           }
           // possible start
           d := ub
           if d > g.r {
               d = g.r
           }
           if d >= a {
               return d, true
           }
       }
       return 0, false
   }
   // try schedule jobs
   for _, job := range jobs {
       placed := false
       for p := 0; p < P; p++ {
           // inspect free intervals in schedule[p]
           sch := schedules[p]
           // sort by start
           sort.Slice(sch, func(i, j int) bool { return sch[i].l < sch[j].l })
           // prev finish for start domain
           prevF := int64(1)
           for k, occ := range sch {
               // free start domain [prevF, occ.l - job.t]
               ub := job.L
               if occ.l-job.t < ub {
                   ub = occ.l - job.t
               }
               if ub >= prevF {
                   if d, ok := findLatest(prevF, ub); ok {
                       // place here
                       assign[job.idx] = p + 1
                       applyDay[job.idx] = d
                       // add processing interval [d, d+job.t]
                       schedules[p] = append(sch, gap{l: d, r: d + job.t})
                       placed = true
                       break
                   }
               }
               // update prevF to max(prevF, occ.r)
               if occ.r > prevF {
                   prevF = occ.r
               }
               // for last, continuation in next iteration
               if k == len(sch)-1 {
                   // after last
                   ub2 := job.L
                   if ub2 >= prevF {
                       if d, ok := findLatest(prevF, ub2); ok {
                           assign[job.idx] = p + 1
                           applyDay[job.idx] = d
                           schedules[p] = append(sch, gap{l: d, r: d + job.t})
                           placed = true
                       }
                   }
               }
               if placed {
                   break
               }
           }
           // if no existing occupied intervals, try whole
           if !placed && len(sch) == 0 {
               ub := job.L
               if ub >= 1 {
                   if d, ok := findLatest(1, ub); ok {
                       assign[job.idx] = p + 1
                       applyDay[job.idx] = d
                       schedules[p] = append(schedules[p], gap{l: d, r: d + job.t})
                       placed = true
                   }
               }
           }
           if placed {
               break
           }
       }
       if !placed {
           fmt.Println("NO")
           return
       }
   }
   // success
   fmt.Println("YES")
   w := bufio.NewWriter(os.Stdout)
   defer w.Flush()
   for i := 0; i < N; i++ {
       fmt.Fprintf(w, "%d %d\n", assign[i], applyDay[i])
   }
}
