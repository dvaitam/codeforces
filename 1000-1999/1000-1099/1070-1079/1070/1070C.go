package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

type Event struct {
   x   int
   typ int // 0=start, 1=end
   idx int
}

var (
   N, K, M    int
   l, r, c, p []int
   events     []Event
   maxP       int
   treeCnt    []int64
   treeSum    []int64
)

func update(tidx, tl, tr, price int, cnt int64) {
   if tl == tr {
       treeCnt[tidx] += cnt
       treeSum[tidx] += cnt * int64(tl)
       return
   }
   mid := (tl + tr) >> 1
   if price <= mid {
       update(tidx<<1, tl, mid, price, cnt)
   } else {
       update(tidx<<1|1, mid+1, tr, price, cnt)
   }
   treeCnt[tidx] = treeCnt[tidx<<1] + treeCnt[tidx<<1|1]
   treeSum[tidx] = treeSum[tidx<<1] + treeSum[tidx<<1|1]
}

// queryF returns sum of smallest x items by price
func queryF(tidx, tl, tr int, x int64) int64 {
   if x <= 0 {
       return 0
   }
   if tl == tr {
       if treeCnt[tidx] < x {
           x = treeCnt[tidx]
       }
       return x * int64(tl)
   }
   leftCnt := treeCnt[tidx<<1]
   mid := (tl + tr) >> 1
   if leftCnt >= x {
       return queryF(tidx<<1, tl, mid, x)
   }
   // take all left, plus from right
   return treeSum[tidx<<1] + queryF(tidx<<1|1, mid+1, tr, x-leftCnt)
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   fmt.Fscan(reader, &N, &K, &M)
   l = make([]int, M)
   r = make([]int, M)
   c = make([]int, M)
   p = make([]int, M)
   events = make([]Event, 0, 2*M)
   for i := 0; i < M; i++ {
       fmt.Fscan(reader, &l[i], &r[i], &c[i], &p[i])
       events = append(events, Event{l[i], 0, i})
       events = append(events, Event{r[i], 1, i})
       if p[i] > maxP {
           maxP = p[i]
       }
   }
   // init segment tree
   size := 4 * (maxP + 5)
   treeCnt = make([]int64, size)
   treeSum = make([]int64, size)
   // sort events by x, then typ (start before end)
   sort.Slice(events, func(i, j int) bool {
       if events[i].x != events[j].x {
           return events[i].x < events[j].x
       }
       return events[i].typ < events[j].typ
   })
   it := 0
   var ans int64
   // sweep over days
   for day := 1; day <= N; day++ {
       // add starting segments
       for it < len(events) && events[it].x == day && events[it].typ == 0 {
           idx := events[it].idx
           update(1, 1, maxP, p[idx], int64(c[idx]))
           it++
       }
       // query sum of smallest K items
       ans += queryF(1, 1, maxP, int64(K))
       // remove ending segments
       for it < len(events) && events[it].x == day && events[it].typ == 1 {
           idx := events[it].idx
           update(1, 1, maxP, p[idx], -int64(c[idx]))
           it++
       }
   }
   fmt.Fprintln(writer, ans)
}
