package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var N, M int
   if _, err := fmt.Fscan(reader, &N, &M); err != nil {
       return
   }
   a := make([]int, N+2)
   for i := 1; i <= N; i++ {
       fmt.Fscan(reader, &a[i])
   }
   // Block size
   B := int(math.Sqrt(float64(N))) + 1
   blocks := (N + B - 1) / B
   nxt := make([]int, N+2)
   cnt := make([]int, N+2)
   last := make([]int, N+2)
   // rebuild block b
   rebuild := func(b int) {
       L := b*B + 1
       R := (b + 1) * B
       if R > N {
           R = N
       }
       for i := R; i >= L; i-- {
           j := i + a[i]
           if j > N {
               nxt[i] = N + 1
               cnt[i] = 1
               last[i] = i
           } else if j >= L && j <= R {
               nxt[i] = nxt[j]
               cnt[i] = cnt[j] + 1
               last[i] = last[j]
           } else {
               nxt[i] = j
               cnt[i] = 1
               last[i] = i
           }
       }
   }
   // initial build
   for b := 0; b < blocks; b++ {
       rebuild(b)
   }
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   for m := 0; m < M; m++ {
       var typ int
       fmt.Fscan(reader, &typ)
       if typ == 0 {
           var idx, val int
           fmt.Fscan(reader, &idx, &val)
           a[idx] = val
           rebuild((idx - 1) / B)
       } else {
           var pos int
           fmt.Fscan(reader, &pos)
           total := 0
           lastPos := pos
           for pos <= N {
               total += cnt[pos]
               lastPos = last[pos]
               pos = nxt[pos]
           }
           fmt.Fprintln(writer, lastPos, total)
       }
   }
}
