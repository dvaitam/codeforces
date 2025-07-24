package main

import (
   "bufio"
   "io"
   "os"
)

func main() {
   rdr := bufio.NewReader(os.Stdin)
   // fast read helpers
   nextNonSpace := func() (byte, error) {
       for {
           b, err := rdr.ReadByte()
           if err != nil {
               return 0, err
           }
           if b > ' ' {
               return b, nil
           }
       }
   }
   readInt := func(first byte) (int, error) {
       neg := false
       b := first
       if b == '-' {
           neg = true
           var err error
           b, err = rdr.ReadByte()
           if err != nil {
               return 0, err
           }
       }
       var v int
       for ; b >= '0' && b <= '9'; b, _ = rdr.ReadByte() {
           v = v*10 + int(b-'0')
       }
       if neg {
           v = -v
       }
       return v, nil
   }
   // read n, m, k
   b, _ := nextNonSpace()
   n, _ := readInt(b)
   b, _ = nextNonSpace()
   m, _ := readInt(b)
   b, _ = nextNonSpace()
   k, _ := readInt(b)
   // read garlands
   type bulb struct { x, y, w int }
   bulbs := make([][]bulb, k)
   for i := 0; i < k; i++ {
       b, _ = nextNonSpace()
       li, _ := readInt(b)
       arr := make([]bulb, li)
       for j := 0; j < li; j++ {
           b, _ = nextNonSpace()
           xi, _ := readInt(b)
           b, _ = nextNonSpace()
           yi, _ := readInt(b)
           b, _ = nextNonSpace()
           wi, _ := readInt(b)
           arr[j] = bulb{xi - 1, yi - 1, wi}
       }
       bulbs[i] = arr
   }
   // read events
   b, _ = nextNonSpace()
   q, _ := readInt(b)
   eventType := make([]byte, q)
   eventArg := make([]int, q)
   // store ask rectangles
   x1 := make([]int, 0, 2048)
   y1 := make([]int, 0, 2048)
   x2 := make([]int, 0, 2048)
   y2 := make([]int, 0, 2048)
   askCount := 0
   for ei := 0; ei < q; ei++ {
       b, err := nextNonSpace()
       if err != nil {
           break
       }
       if b == 'S' { // SWITCH
           eventType[ei] = 'S'
           // skip rest of word
           for {
               c, _ := rdr.ReadByte()
               if c <= ' ' {
                   break
               }
           }
           b, _ = nextNonSpace()
           gi, _ := readInt(b)
           eventArg[ei] = gi - 1
       } else { // ASK
           eventType[ei] = 'A'
           // skip rest of word
           for i := 0; i < 2; i++ {
               // skip 'S','K'
               _ = rdr.ReadByte()
           }
           // now skip till space
           for {
               c, _ := rdr.ReadByte()
               if c <= ' ' {
                   break
               }
           }
           // read x1,y1,x2,y2
           b, _ = nextNonSpace()
           xi1, _ := readInt(b)
           b, _ = nextNonSpace()
           yi1, _ := readInt(b)
           b, _ = nextNonSpace()
           xi2, _ := readInt(b)
           b, _ = nextNonSpace()
           yi2, _ := readInt(b)
           x1 = append(x1, xi1-1)
           y1 = append(y1, yi1-1)
           x2 = append(x2, xi2-1)
           y2 = append(y2, yi2-1)
           eventArg[ei] = askCount
           askCount++
       }
   }
   // precompute S[i][j]
   S := make([][]int64, k)
   for i := 0; i < k; i++ {
       si := make([]int64, askCount)
       for _, b := range bulbs[i] {
           bx, by, bw := b.x, b.y, int64(b.w)
           for j := 0; j < askCount; j++ {
               if bx >= x1[j] && bx <= x2[j] && by >= y1[j] && by <= y2[j] {
                   si[j] += bw
               }
           }
       }
       S[i] = si
   }
   // process events
   state := make([]bool, k)
   for i := range state {
       state[i] = true
   }
   w := bufio.NewWriter(os.Stdout)
   defer w.Flush()
   for ei := 0; ei < q; ei++ {
       if eventType[ei] == 'S' {
           gi := eventArg[ei]
           state[gi] = !state[gi]
       } else {
           aj := eventArg[ei]
           var ans int64
           for i := 0; i < k; i++ {
               if state[i] {
                   ans += S[i][aj]
               }
           }
           // write answer
           w.WriteString(int64ToStr(ans))
           w.WriteByte('\n')
       }
   }
}

// int64 to string
func int64ToStr(v int64) string {
   if v == 0 {
       return "0"
   }
   buf := make([]byte, 0, 20)
   neg := false
   if v < 0 {
       neg = true
       v = -v
   }
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
