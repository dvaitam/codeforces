package main

import (
   "bufio"
   "fmt"
   "os"
   "strconv"
)

// Segment represents a constant-speed segment of a car
type Segment struct {
   speed int64
   time  int64
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   // read n and s
   var line string
   var err error
   line, err = reader.ReadString('\n')
   if err != nil {
       return
   }
   parts := splitFields(line)
   n, _ := strconv.Atoi(parts[0])
   // s not used explicitly
   // s, _ := strconv.ParseInt(parts[1], 10, 64)

   cars := make([][]Segment, n)
   for i := 0; i < n; i++ {
       // read line for car i
       line, err = reader.ReadString('\n')
       if err != nil {
           return
       }
       parts = splitFields(line)
       k, _ := strconv.Atoi(parts[0])
       segs := make([]Segment, k)
       idx := 1
       for j := 0; j < k; j++ {
           v, _ := strconv.ParseInt(parts[idx], 10, 64)
           t, _ := strconv.ParseInt(parts[idx+1], 10, 64)
           segs[j] = Segment{speed: v, time: t}
           idx += 2
       }
       cars[i] = segs
   }

   var totalCross int64 = 0
   // for each unordered pair of cars
   for i := 0; i < n; i++ {
       for j := i + 1; j < n; j++ {
           totalCross += countCrosses(cars[i], cars[j])
       }
   }
   fmt.Println(totalCross)
}

// countCrosses counts overtakes between two cars' segment schedules
func countCrosses(a, b []Segment) int64 {
   var i, j int
   var remA, remB int64
   var vA, vB int64
   var posA, posB int64
   // initialize first segments
   if len(a) > 0 {
       remA = a[0].time
       vA = a[0].speed
   }
   if len(b) > 0 {
       remB = b[0].time
       vB = b[0].speed
   }
   var crosses int64 = 0
   // iterate until one runs out of segments
   for i < len(a) && j < len(b) {
       // current overlapping duration
       dur := remA
       if remB < dur {
           dur = remB
       }
       // difference at start of interval
       delta0 := posA - posB
       dv := vA - vB
       // check if crossing inside (0,dur)
       if dv > 0 {
           if delta0 < 0 && delta0+dv*dur > 0 {
               crosses++
           }
       } else if dv < 0 {
           if delta0 > 0 && delta0+dv*dur < 0 {
               crosses++
           }
       }
       // advance positions
       posA += vA * dur
       posB += vB * dur
       remA -= dur
       remB -= dur
       if remA == 0 {
           i++
           if i < len(a) {
               remA = a[i].time
               vA = a[i].speed
           }
       }
       if remB == 0 {
           j++
           if j < len(b) {
               remB = b[j].time
               vB = b[j].speed
           }
       }
   }
   return crosses
}

// splitFields splits a line into fields, handling spaces and tabs
func splitFields(s string) []string {
   var res []string
   field := ""
   for i := 0; i < len(s); i++ {
       c := s[i]
       if c == ' ' || c == '\t' || c == '\r' || c == '\n' {
           if field != "" {
               res = append(res, field)
               field = ""
           }
       } else {
           field += string(c)
       }
   }
   if field != "" {
       res = append(res, field)
   }
   return res
}
