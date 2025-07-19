package main

import (
   "bufio"
   "os"
   "sort"
)

// Event represents a change in coverage at position X by flag (+1 or -1)
type Event struct {
   x    int64
   flag int
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   // fast integer reader
   readInt := func() (int64, error) {
       var x int64
       var neg bool
       // read next non-space
       for {
           b, err := reader.ReadByte()
           if err != nil {
               return 0, err
           }
           if b == '-' {
               neg = true
               break
           }
           if b >= '0' && b <= '9' {
               x = int64(b - '0')
               break
           }
       }
       // read rest digits
       for {
           b, err := reader.ReadByte()
           if err != nil {
               break
           }
           if b < '0' || b > '9' {
               break
           }
           x = x*10 + int64(b-'0')
       }
       if neg {
           x = -x
       }
       return x, nil
   }

   // read n
   n64, err := readInt()
   if err != nil {
       return
   }
   n := int(n64)
   events := make([]Event, 0, 2*n)
   for i := 0; i < n; i++ {
       x, _ := readInt()
       y, _ := readInt()
       events = append(events, Event{x: x, flag: 1})
       events = append(events, Event{x: y + 1, flag: -1})
   }
   sort.Slice(events, func(i, j int) bool {
       return events[i].x < events[j].x
   })
   ans := make([]int64, n+1)
   var currX int64 = 1
   currCnt := 0
   m := len(events)
   for i := 0; i < m; i++ {
       e := events[i]
       // add length from currX to e.x
       length := e.x - currX
       if currCnt > 0 && currCnt <= n {
           ans[currCnt] += length
       }
       currCnt += e.flag
       // combine events with same x
       currX = e.x
       for i+1 < m && events[i+1].x == e.x {
           i++
           currCnt += events[i].flag
       }
       // next segment starts at this x
       // currX already set
   }
   // output
   w := bufio.NewWriter(os.Stdout)
   defer w.Flush()
   for i := 1; i <= n; i++ {
       // write ans[i]
       // use manual write to avoid imports
       w.WriteString(int64ToString(ans[i]))
       if i < n {
           w.WriteByte(' ')
       }
   }
   w.WriteByte('\n')
}

// int64ToString converts integer to string
func int64ToString(v int64) string {
   if v == 0 {
       return "0"
   }
   neg := v < 0
   if neg {
       v = -v
   }
   // max length for int64 is 20 digits
   buf := make([]byte, 0, 20)
   for v > 0 {
       buf = append(buf, byte('0'+v%10))
       v /= 10
   }
   if neg {
       buf = append(buf, '-')
   }
   // reverse
   for i, j := 0, len(buf)-1; i < j; i, j = i+1, j-1 {
       buf[i], buf[j] = buf[j], buf[i]
   }
   return string(buf)
}
