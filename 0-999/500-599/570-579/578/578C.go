package main

import (
   "bufio"
   "fmt"
   "os"
   "strconv"
)

type Line struct {
   a, b int64
}
type Info struct {
   line  Line
   start float64
}

func minFloat(a, b float64) float64 {
   if a < b {
       return a
   }
   return b
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   readInt := func() int64 {
       s, _ := reader.ReadString(' ')
       // handle newline
       if len(s) > 0 && (s[len(s)-1] == '\n' || s[len(s)-1] == ' ') {
           s = s[:len(s)-1]
       }
       v, _ := strconv.ParseInt(s, 10, 64)
       return v
   }
   // alternative: read whole tokens
   var n int64
   // first token is n
   fmt.Fscan(reader, &n)
   n++
   xu := make([]Line, n+2)
   xu[1] = Line{0, 0}
   for i := int64(2); i <= n; i++ {
       var d int64
       fmt.Fscan(reader, &d)
       xu[i].a = i - 1
       xu[i].b = xu[i-1].b + d
   }
   maxv := make([]Info, n+2)
   minv := make([]Info, n+2)
   var maxc, minc int
   // upper envelope
   for i := int64(1); i <= n; i++ {
       for maxc > 0 {
           prev := maxv[maxc]
           yuan := float64(prev.line.a)*prev.start + float64(prev.line.b)
           now := float64(xu[i].a)*prev.start + float64(xu[i].b)
           if yuan > now {
               break
           }
           maxc--
       }
       maxc++
       maxv[maxc].line = xu[i]
       if maxc == 1 {
           maxv[maxc].start = -20000
       } else {
           cur := maxv[maxc].line
           pre := maxv[maxc-1].line
           maxv[maxc].start = float64(cur.b-pre.b) / float64(pre.a-cur.a)
       }
   }
   // lower envelope
   for i := n; i >= 1; i-- {
       for minc > 0 {
           prev := minv[minc]
           yuan := float64(prev.line.a)*prev.start + float64(prev.line.b)
           now := float64(xu[i].a)*prev.start + float64(xu[i].b)
           if yuan < now {
               break
           }
           minc--
       }
       minc++
       minv[minc].line = xu[i]
       if minc == 1 {
           minv[minc].start = -20000
       } else {
           cur := minv[minc].line
           pre := minv[minc-1].line
           minv[minc].start = float64(cur.b-pre.b) / float64(pre.a-cur.a)
       }
   }
   // two pointers
   pa, pb := 1, 1
   ans := 1e50
   for pa <= maxc && pb <= minc {
       var wei float64
       if maxv[pa].start > minv[pb].start {
           wei = maxv[pa].start
       } else {
           wei = minv[pb].start
       }
       v1 := float64(maxv[pa].line.a)*wei + float64(maxv[pa].line.b)
       v2 := float64(minv[pb].line.a)*wei + float64(minv[pb].line.b)
       now := v1 - v2
       if now < ans {
           ans = now
       }
       if pa == maxc {
           pb++
       } else if pb == minc {
           pa++
       } else if maxv[pa+1].start > minv[pb+1].start {
           pb++
       } else {
           pa++
       }
   }
   fmt.Printf("%.10f\n", ans)
}
