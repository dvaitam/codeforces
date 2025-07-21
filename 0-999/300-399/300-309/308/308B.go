package main

import (
   "bufio"
   "bytes"
   "io"
   "os"
   "strconv"
   "strings"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   // read n, r, c
   line, _, err := reader.ReadLine()
   if err != nil && err != io.EOF {
       panic(err)
   }
   parts := strings.Fields(string(line))
   if len(parts) < 3 {
       return
   }
   n, _ := strconv.Atoi(parts[0])
   r, _ := strconv.Atoi(parts[1])
   c, _ := strconv.Atoi(parts[2])
   // read text line
   text, err := reader.ReadString('\n')
   if err != nil && err != io.EOF {
       panic(err)
   }
   text = strings.TrimRight(text, "\n")
   words := strings.Fields(text)
   if len(words) > n {
       words = words[:n]
   }
   // lengths
   lens := make([]int, len(words))
   for i, w := range words {
       lens[i] = len(w)
   }
   // next one line
   N := len(words)
   nxt := make([]uint32, N+1)
   sum := 0
   var j int
   for i := 0; i < N; i++ {
       if j < i {
           j = i
           sum = 0
       }
       // try extend j
       for j < N && sum + lens[j] + (j - i) <= c {
           sum += lens[j]
           j++
       }
       nxt[i] = uint32(j)
       sum -= lens[i]
   }
   nxt[N] = uint32(N)
   // binary lifting
   // max bits for r
   maxB := 0
   for tmp := r; tmp > 0; tmp >>= 1 {
       maxB++
   }
   dp := make([][]uint32, maxB)
   dp[0] = nxt
   for b := 1; b < maxB; b++ {
       dp[b] = make([]uint32, N+1)
       prev := dp[b-1]
       cur := dp[b]
       for i := 0; i <= N; i++ {
           cur[i] = prev[prev[i]]
       }
   }
   // find best segment
   bestLen := 0
   bestStart := 0
   var bestEnd uint32
   for i := 0; i < N; i++ {
       pos := uint32(i)
       rem := r
       for b := 0; rem > 0; b++ {
           if rem&1 != 0 {
               pos = dp[b][pos]
           }
           rem >>= 1
       }
       length := int(pos) - i
       if length > bestLen {
           bestLen = length
           bestStart = i
           bestEnd = pos
       }
       // early exit if all words
       if bestLen == N {
           break
       }
   }
   // output
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   cur := uint32(bestStart)
   lines := 0
   for lines < r && cur < bestEnd {
       nxtPos := nxt[cur]
       if nxtPos > bestEnd {
           nxtPos = bestEnd
       }
       // join words[cur:nxtPos]
       var buf bytes.Buffer
       for k := cur; k < nxtPos; k++ {
           if k > cur {
               buf.WriteByte(' ')
           }
           buf.WriteString(words[k])
       }
       writer.Write(buf.Bytes())
       writer.WriteByte('\n')
       cur = nxtPos
       lines++
   }
}
