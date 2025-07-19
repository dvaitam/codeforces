package main

import (
  "bufio"
  "fmt"
  "os"
)

const (
  MAXN = 100005
  LOGN = 18
)

var (
  fa, dep, siz, son, top, dfn [MAXN]int
  lg                  [MAXN]int
  mn, mx              [LOGN][MAXN]int
  children            [MAXN][]int
  dfnCount            int
)

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

func cmp(x, y int) int {
  if x == 0 || y == 0 {
    return x + y
  }
  if dfn[x] < dfn[y] {
    return x
  }
  return y
}

func cmp1(x, y int) int {
  if x == 0 || y == 0 {
    return x + y
  }
  if dfn[x] > dfn[y] {
    return x
  }
  return y
}

func dfs1(u int) {
  siz[u] = 1
  dfnCount++
  dfn[u] = dfnCount
  for _, v := range children[u] {
    dep[v] = dep[u] + 1
    dfs1(v)
    siz[u] += siz[v]
    if siz[v] > siz[son[u]] {
      son[u] = v
    }
  }
}

func dfs2(u, t int) {
  top[u] = t
  if son[u] != 0 {
    dfs2(son[u], t)
  }
  for _, v := range children[u] {
    if v != son[u] {
      dfs2(v, v)
    }
  }
}

func lca(x, y int) int {
  if x == 0 || y == 0 {
    return x + y
  }
  for top[x] != top[y] {
    if dep[top[x]] > dep[top[y]] {
      x, y = y, x
    }
    y = fa[top[y]]
  }
  if dep[x] < dep[y] {
    return x
  }
  return y
}

func getMin(l, r int) int {
  if l > r {
    return 0
  }
  k := lg[r-l+1]
  return cmp(mn[k][l], mn[k][r-(1<<k)+1])
}

func getMax(l, r int) int {
  if l > r {
    return 0
  }
  k := lg[r-l+1]
  return cmp1(mx[k][l], mx[k][r-(1<<k)+1])
}

func main() {
  reader := bufio.NewReader(os.Stdin)
  writer := bufio.NewWriter(os.Stdout)
  defer writer.Flush()

  n := readInt(reader)
  m := readInt(reader)
  for i := 1; i <= n; i++ {
    children[i] = children[i][:0]
  }
  for i := 2; i <= n; i++ {
    fa[i] = readInt(reader)
    children[fa[i]] = append(children[fa[i]], i)
  }
  dep[1] = 1
  dfs1(1)
  dfs2(1, 1)

  lg[1] = 0
  for i := 2; i <= n; i++ {
    lg[i] = lg[i>>1] + 1
  }
  for i := 1; i <= n; i++ {
    mn[0][i] = i
    mx[0][i] = i
  }
  for k := 1; k < LOGN; k++ {
    step := 1 << k
    if step > n {
      break
    }
    half := step >> 1
    for i := 1; i+step-1 <= n; i++ {
      mn[k][i] = cmp(mn[k-1][i], mn[k-1][i+half])
      mx[k][i] = cmp1(mx[k-1][i], mx[k-1][i+half])
    }
  }

  for ; m > 0; m-- {
    l := readInt(reader)
    r := readInt(reader)
    x := getMin(l, r)
    y := getMax(l, r)
    lx := getMin(l, x-1)
    rx := getMin(x+1, r)
    ans1 := lca(cmp(lx, rx), y)
    ly := getMax(l, y-1)
    ry := getMax(y+1, r)
    ans2 := lca(cmp1(ly, ry), x)
    if dep[ans1] > dep[ans2] {
      fmt.Fprintln(writer, x, dep[ans1])
    } else {
      fmt.Fprintln(writer, y, dep[ans2])
    }
  }
}
