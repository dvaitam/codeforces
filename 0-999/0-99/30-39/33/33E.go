package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
   "strconv"
   "strings"
)

func checkErr(err error) {
   if err != nil {
       panic(err)
   }
}

func parseHM(s string) int {
   // s = HH:MM
   parts := strings.Split(s, ":")
   h, _ := strconv.Atoi(parts[0])
   m, _ := strconv.Atoi(parts[1])
   return h*60 + m
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var m, n, k int
   _, err := fmt.Fscan(reader, &m, &n, &k)
   checkErr(err)
   subj := make([]string, m)
   for i := 0; i < m; i++ {
       fmt.Fscan(reader, &subj[i])
   }
   tsolve := make(map[string]int)
   times := make([]int, m)
   for i := 0; i < m; i++ {
       fmt.Fscan(reader, &times[i])
       tsolve[subj[i]] = times[i]
   }
   // break segments: sleep, breakfast, lunch, dinner
   type seg struct{ start, end, idx int }
   breaks := make([]seg, 4)
   for i := 0; i < 4; i++ {
       var s string
       fmt.Fscan(reader, &s)
       parts := strings.Split(s, "-")
       st := parseHM(parts[0])
       en := parseHM(parts[1])
       breaks[i] = seg{st, en, i + 1}
   }
   // build work minutes
   totalMin := k * 1440
   workTimes := make([]int, 0, totalMin)
   // per day, mark breaks
   dayBreaks := make([][2]bool, 1440)
   // Precompute break flags for one day
   fb := make([]bool, 1440)
   for _, b := range breaks {
       for t := b.start; t <= b.end; t++ {
           fb[t] = true
       }
   }
   // global workTimes
   for d := 0; d < k; d++ {
       base := d * 1440
       for t := 0; t < 1440; t++ {
           if !fb[t] {
               workTimes = append(workTimes, base+t)
           }
       }
   }
   capMax := len(workTimes)
   // read students/jobs
   type job struct{ p, c, cap, idx int }
   jobs := make([]job, 0, n)
   for i := 0; i < n; i++ {
       var ssub, tstr string
       var di, ci int
       fmt.Fscan(reader, &ssub, &di, &tstr, &ci)
       p, ok := tsolve[ssub]
       if !ok {
           continue
       }
       // deadline minute
       dd := parseHM(tstr)
       deadline := (di-1)*1440 + dd
       // must finish strictly before exam: finish <= deadline-1
       // compute cap = count of workTimes < deadline
       // binary search
       // find first index where workTimes[idx] >= deadline
       lo, hi := 0, len(workTimes)
       for lo < hi {
           mid := (lo + hi) / 2
           if workTimes[mid] < deadline {
               lo = mid + 1
           } else {
               hi = mid
           }
       }
       cap := lo
       jobs = append(jobs, job{p: p, c: ci, cap: cap, idx: i + 1})
   }
   // sort jobs by cap increasing
   sort.Slice(jobs, func(i, j int) bool {
       return jobs[i].cap < jobs[j].cap
   })
   // dp
   negInf := -1 << 60
   dp := make([]int, capMax+1)
   for i := 1; i <= capMax; i++ {
       dp[i] = negInf
   }
   // picks[i][t] = picked job i at t
   picks := make([][]bool, len(jobs))
   for i := range picks {
       picks[i] = make([]bool, capMax+1)
   }
   curCap := capMax
   for i, jb := range jobs {
       // new capacity = jb.cap
       // clear dp[t] for t > jb.cap
       for t := jb.cap + 1; t <= curCap; t++ {
           dp[t] = negInf
       }
       // knapsack step
       for t := jb.cap; t >= jb.p; t-- {
           if dp[t-jb.p] + jb.c > dp[t] {
               dp[t] = dp[t-jb.p] + jb.c
               picks[i][t] = true
           }
       }
       curCap = jb.cap
   }
   // find best
   bestProfit := 0
   bestT := 0
   for t := 0; t <= curCap; t++ {
       if dp[t] > bestProfit {
           bestProfit = dp[t]
           bestT = t
       }
   }
   // reconstruct
   sel := make([]job, 0)
   t := bestT
   for i := len(jobs) - 1; i >= 0; i-- {
       if t >= 0 && picks[i][t] {
           sel = append(sel, jobs[i])
           t -= jobs[i].p
       }
   }
   // reverse sel to chronological order (by cap)
   for i, j := 0, len(sel)-1; i < j; i, j = i+1, j-1 {
       sel[i], sel[j] = sel[j], sel[i]
   }
   // prepare break segments sorted by end
   sort.Slice(breaks, func(i, j int) bool { return breaks[i].end < breaks[j].end })
   // output
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fprintln(writer, bestProfit)
   fmt.Fprintln(writer, len(sel))
   // schedule
   offset := 0
   for _, jb := range sel {
       // start and finish
       startG := workTimes[offset]
       finishG := workTimes[offset+jb.p-1]
       offset += jb.p
       // compute time index
       minuteInDay := startG % 1440
       ti := 4
       prev := -1
       for _, b := range breaks {
           if b.end < minuteInDay && b.end > prev {
               ti = b.idx
               prev = b.end
           }
       }
       // start time HH:MM
       h0 := minuteInDay / 60
       m0 := minuteInDay % 60
       // finish day and time
       finDay := finishG/1440 + 1
       finMin := finishG % 1440
       h1 := finMin / 60
       m1 := finMin % 60
       fmt.Fprintf(writer, "%d %d %02d:%02d %d %02d:%02d\n", jb.idx, ti, h0, m0, finDay, h1, m1)
   }
