package main

import (
  "bufio"
  "fmt"
  "os"
)

const mod = 1000000007

func readInt(r *bufio.Reader) int {
  x, sign := 0, 1
  b, err := r.ReadByte()
  if err != nil {
    return 0
  }
  for (b < '0' || b > '9') && b != '-' {
    b, _ = r.ReadByte()
  }
  if b == '-' {
    sign = -1
    b, _ = r.ReadByte()
  }
  for b >= '0' && b <= '9' {
    x = x*10 + int(b-'0')
    b, _ = r.ReadByte()
  }
  return x * sign
}

func main() {
  reader := bufio.NewReader(os.Stdin)
  writer := bufio.NewWriter(os.Stdout)
  defer writer.Flush()

  n := readInt(reader)
  q := readInt(reader)

  // read binary string of length n
  s := make([]byte, n)
  i := 0
  for i < n {
    b, err := reader.ReadByte()
    if err != nil {
      break
    }
    if b == '\n' || b == '\r' {
      continue
    }
    s[i] = b
    i++
  }

  // prefix sum of ones
  sum := make([]int, n+1)
  for i = 1; i <= n; i++ {
    sum[i] = sum[i-1]
    if s[i-1] == '1' {
      sum[i]++
    }
  }

  // precompute powers of 2 and prefix sums
  g := make([]int, n+1)
  f := make([]int, n+1)
  g[0] = 1
  for i = 1; i <= n; i++ {
    g[i] = g[i-1] * 2 % mod
  }
  for i = 1; i <= n; i++ {
    f[i] = f[i-1] + g[i-1]
    if f[i] >= mod {
      f[i] -= mod
    }
  }

  for ; q > 0; q-- {
    l := readInt(reader)
    r2 := readInt(reader)
    cnt := sum[r2] - sum[l-1]
    zeros := (r2 - l + 1) - cnt
    ans := int64(f[cnt]) * int64(g[zeros]) % mod
    fmt.Fprintln(writer, ans)
  }
}
