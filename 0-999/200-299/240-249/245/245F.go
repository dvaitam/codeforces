package main

import (
   "bufio"
   "fmt"
   "os"
   "strconv"
   "strings"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   // read first line: n and m
   line, err := reader.ReadString('\n')
   if err != nil && len(line) == 0 {
       return
   }
   parts := strings.Fields(line)
   if len(parts) < 2 {
       return
   }
   n, _ := strconv.Atoi(parts[0])
   m, _ := strconv.Atoi(parts[1])
   // prepare month prefix sums for 2012 (leap year)
   daysBefore := [13]int{}
   md := []int{31, 29, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}
   for i := 1; i <= 12; i++ {
       daysBefore[i] = daysBefore[i-1] + md[i-1]
   }
   var q []int
   for {
       rec, err := reader.ReadString('\n')
       if err != nil && len(rec) == 0 {
           break
       }
       // handle last line without newline
       if err != nil {
           rec = strings.TrimRight(rec, "\r\n")
       }
       // timestamp string "YYYY-MM-DD HH:MM:SS"
       if len(rec) < 19 {
           continue
       }
       ts := rec[:19]
       // parse components
       month, _ := strconv.Atoi(ts[5:7])
       day, _ := strconv.Atoi(ts[8:10])
       hour, _ := strconv.Atoi(ts[11:13])
       minute, _ := strconv.Atoi(ts[14:16])
       second, _ := strconv.Atoi(ts[17:19])
       // compute seconds since 2012-01-01 00:00:00
       days := daysBefore[month-1] + (day - 1)
       cur := days*86400 + hour*3600 + minute*60 + second
       // add to queue
       q = append(q, cur)
       // remove items older than cur-n+1 (keep >= cur-n+1)
       limit := cur - n + 1
       i := 0
       for i < len(q) && q[i] < limit {
           i++
       }
       if i > 0 {
           q = q[i:]
       }
       if len(q) >= m {
           fmt.Println(ts)
           return
       }
   }
   fmt.Println(-1)
}
