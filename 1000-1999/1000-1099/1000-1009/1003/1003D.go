package main

import (
  "bufio"
  "bytes"
  "os"
  "sort"
  "strconv"
)

func main() {
  br := bufio.NewReader(os.Stdin)
  bw := bufio.NewWriter(os.Stdout)
  defer bw.Flush()
  // read n, m
  n, _ := readInt(br)
  m, _ := readInt(br)
  a := make([]int64, n)
  for i := int64(0); i < n; i++ {
    a[i], _ = readInt(br)
  }
  sort.Slice(a, func(i, j int) bool { return a[i] < a[j] })
  // compress unique values and counts
  var b []int64
  var c []int64
  for i := 0; i < len(a); i++ {
    if i == 0 || a[i] != a[i-1] {
      b = append(b, a[i])
      c = append(c, 1)
    } else {
      c[len(c)-1]++
    }
  }
  // process queries
  var buf bytes.Buffer
  for qi := int64(0); qi < m; qi++ {
    j, _ := readInt(br)
    var ans int64
    rem := j
    for i := len(b) - 1; i >= 0; i-- {
      if rem <= 0 {
        break
      }
      t := rem / b[i]
      if t > c[i] {
        t = c[i]
      }
      rem -= t * b[i]
      ans += t
    }
    if rem != 0 {
      buf.WriteString("-1\n")
    } else {
      buf.WriteString(strconv.FormatInt(ans, 10))
      buf.WriteByte('\n')
    }
  }
  bw.Write(buf.Bytes())
}

// readInt reads next integer (int64) from bufio.Reader
func readInt(br *bufio.Reader) (int64, error) {
  var x int64
  var sign int64 = 1
  var ch byte
  var err error
  // skip non-numeric characters
  for {
    ch, err = br.ReadByte()
    if err != nil {
      return 0, err
    }
    if (ch >= '0' && ch <= '9') || ch == '-' {
      break
    }
  }
  if ch == '-' {
    sign = -1
    ch, err = br.ReadByte()
    if err != nil {
      return 0, err
    }
  }
  for ; ch >= '0' && ch <= '9'; ch, err = br.ReadByte() {
    if err != nil {
      break
    }
    x = x*10 + int64(ch-'0')
  }
  return x * sign, nil
}
