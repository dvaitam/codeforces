package main

import (
    "bufio"
    "fmt"
    "os"
    "sort"
)

type Event struct {
    time int64
    typ  int // +1 start, -1 end
    id   int
    d    int64
}

type User struct {
    rate int64
}

func processSegment(duration int64, b int64, active []*User) int64 {
    var ans int64
    m := int64(len(active))
    if m == 0 {
        return 0
    }
    for duration > 0 {
        var sum int64
        for _, u := range active {
            sum += u.rate
        }
        if sum <= b {
            t := (b - sum)/m + 1
            if t > duration {
                t = duration
            }
            ans += t*sum + m*t*(t-1)/2
            for _, u := range active {
                u.rate += t
            }
            duration -= t
        } else {
            for _, u := range active {
                u.rate /= 2
            }
            duration--
        }
    }
    return ans
}

func main() {
    in := bufio.NewReader(os.Stdin)
    out := bufio.NewWriter(os.Stdout)
    defer out.Flush()

    var n int
    var b int64
    if _, err := fmt.Fscan(in, &n, &b); err != nil {
        return
    }
    events := make([]Event, 0, 2*n)
    for i := 0; i < n; i++ {
        var s, f, d int64
        fmt.Fscan(in, &s, &f, &d)
        events = append(events, Event{time: s, typ: 1, id: i, d: d})
        events = append(events, Event{time: f + 1, typ: -1, id: i})
    }
    sort.Slice(events, func(i, j int) bool {
        if events[i].time == events[j].time {
            return events[i].typ > events[j].typ
        }
        return events[i].time < events[j].time
    })

    active := make(map[int]*User)
    var activeList []*User
    var curTime int64
    var ans int64
    if len(events) > 0 {
        curTime = events[0].time
    }
    idx := 0
    for idx < len(events) {
        t := events[idx].time
        ans += processSegment(t-curTime, b, activeList)
        curTime = t
        for idx < len(events) && events[idx].time == t {
            e := events[idx]
            if e.typ == 1 {
                u := &User{rate: e.d}
                active[e.id] = u
                activeList = append(activeList, u)
            } else {
                if u, ok := active[e.id]; ok {
                    delete(active, e.id)
                    for i := range activeList {
                        if activeList[i] == u {
                            activeList = append(activeList[:i], activeList[i+1:]...)
                            break
                        }
                    }
                }
            }
            idx++
        }
    }
    // process remaining time if any
    ans += processSegment(0, b, activeList)
    fmt.Fprintln(out, ans)
}

