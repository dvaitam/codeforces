package main

import (
   "bufio"
   "fmt"
   "os"
   "strconv"
   "unicode"
)

// Fast IO reader
type FastReader struct {
   r *bufio.Reader
}

func NewFastReader(r *bufio.Reader) *FastReader {
   return &FastReader{r: r}
}

func (fr *FastReader) ReadInt() (int, error) {
   bs := make([]byte, 0, 16)
   for {
       b, err := fr.r.ReadByte()
       if err != nil {
           return 0, err
       }
       if !unicode.IsSpace(rune(b)) {
           fr.r.UnreadByte()
           break
       }
   }
   // read sign
   sign := 1
   b, _ := fr.r.ReadByte()
   if b == '-' {
       sign = -1
   } else {
       fr.r.UnreadByte()
   }
   // read digits
   for {
       b, err := fr.r.ReadByte()
       if err != nil {
           break
       }
       if unicode.IsDigit(rune(b)) {
           bs = append(bs, b)
       } else {
           fr.r.UnreadByte()
           break
       }
   }
   if len(bs) == 0 {
       return 0, nil
   }
   v, _ := strconv.Atoi(string(bs))
   return sign * v, nil
}

func (fr *FastReader) ReadInt64() (int64, error) {
   bs := make([]byte, 0, 20)
   for {
       b, err := fr.r.ReadByte()
       if err != nil {
           return 0, err
       }
       if !unicode.IsSpace(rune(b)) {
           fr.r.UnreadByte()
           break
       }
   }
   sign := int64(1)
   b, _ := fr.r.ReadByte()
   if b == '-' {
       sign = -1
   } else {
       fr.r.UnreadByte()
   }
   for {
       b, err := fr.r.ReadByte()
       if err != nil {
           break
       }
       if unicode.IsDigit(rune(b)) {
           bs = append(bs, b)
       } else {
           fr.r.UnreadByte()
           break
       }
   }
   if len(bs) == 0 {
       return 0, nil
   }
   v, _ := strconv.ParseInt(string(bs), 10, 64)
   return sign * v, nil
}

func main() {
   reader := NewFastReader(bufio.NewReader(os.Stdin))
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   nInt, _ := reader.ReadInt()
   n := nInt
   h := make([]int64, n)
   for i := 0; i < n; i++ {
       v, _ := reader.ReadInt64()
       h[i] = v
   }
   if n <= 1 {
       if n == 1 {
           fmt.Fprint(writer, h[0])
       }
       return
   }
   // queue of indices where h[i+1] - h[i] >= 2
   inQ := make([]bool, n-1)
   queue := make([]int, 0, n)
   for i := 0; i < n-1; i++ {
       if h[i+1]-h[i] >= 2 {
           inQ[i] = true
           queue = append(queue, i)
       }
   }
   for qi := 0; qi < len(queue); qi++ {
       i := queue[qi]
       inQ[i] = false
       // recheck condition
       if h[i+1]-h[i] < 2 {
           continue
       }
       m := (h[i+1] - h[i]) / 2
       h[i] += m
       h[i+1] -= m
       // check neighbors
       // left neighbor (diff at i-1)
       if i-1 >= 0 && !inQ[i-1] && h[i]-h[i-1] >= 2 {
           inQ[i-1] = true
           queue = append(queue, i-1)
       }
       // same position (diff at i)
       if i < n-1 && !inQ[i] && h[i+1]-h[i] >= 2 {
           inQ[i] = true
           queue = append(queue, i)
       }
       // right neighbor (diff at i+1)
       if i+1 < n-1 && !inQ[i+1] && h[i+2]-h[i+1] >= 2 {
           inQ[i+1] = true
           queue = append(queue, i+1)
       }
   }
   // output
   for i, v := range h {
       if i > 0 {
           writer.WriteByte(' ')
       }
       writer.WriteString(strconv.FormatInt(v, 10))
   }
}
