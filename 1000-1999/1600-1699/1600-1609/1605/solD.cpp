#include <bits/stdc++.h>
namespace IO {
  #define iL (1 << 20)
  char ibuf[iL], *iS = ibuf + iL, *iT = ibuf + iL;
  #define gc() ((iS == iT) ? (iT = (iS = ibuf) + fread(ibuf, 1, iL, stdin), iS == iT ? EOF : *iS ++) : *iS ++)
  template<class T> inline void read(T &x) {
    x = 0;int f = 0;char ch = gc();
    for (; !isdigit(ch); f |= ch == '-', ch = gc());
    for (; isdigit(ch); x = (x << 1) + (x << 3) + (ch ^ 48), ch = gc());
    x = (f == 1 ? ~ x + 1 : x);
  }
  char Out[iL], *iter = Out;
  #define flush() fwrite(Out, 1, iter - Out, stdout), iter = Out
  template<class T> inline void write(T x, char ch = '\n') {
    T l, c[35];
    if (x < 0) *iter ++ = '-', x = ~ x + 1;
    for (l = 0; !l || x; c[l] = x % 10, l++, x /= 10);
    for (; l; -- l, *iter ++ = c[l] + '0');*iter ++ = ch;
    flush();
  }
}
using namespace IO;
#define N 200005
#define ll long long
#define DEBUG
using namespace std;
struct Edge { int x, y; }e[N * 2];
int t, n, head[N], nxt[N * 2], cnt, val[N], dep[N];
vector<int> v2[2];
void Add(int x, int y) { e[++ cnt] = (Edge){x, y}, nxt[cnt] = head[x], head[x] = cnt; }
void Dfs(int x, int fa) {
  dep[x] = dep[fa] + 1; v2[!(dep[x] & 1)].push_back(x);
  for (int i = head[x]; i; i = nxt[i]) {
    int y = e[i].y; if (y == fa) continue;
    Dfs(y, x);
  }
}
int main() {
#ifndef ONLINE_JUDGE
  freopen("test.in", "r", stdin);
  freopen("test.out", "w", stdout);
#endif
  read(t);
  while (t --) {
    read(n); cnt = 0; for (int i = 1; i <= n; i++) head[i] = 0;
    for (int i = 1, x, y; i < n; i++) read(x), read(y), Add(x, y), Add(y, x);
    v2[0].clear(), v2[1].clear(); vector<int> v[35];
    for (int i = 1; i <= n; i++) v[((int)log2(i))].push_back(i);
    int cnt[2]; cnt[0] = cnt[1] = 0;
    Dfs(1, 0); if (v2[0].size() > v2[1].size()) swap(v2[0], v2[1]); int num = v2[0].size();
    for (int i = 0; i < 30; i++) {
      int tmp = !((num >> i) & 1);
      for (int j = 0; j < v[i].size(); j++) val[v2[tmp][cnt[tmp] ++]] = v[i][j];
    }
    for (int i = 1; i <= n; i++) write(val[i], ' ');
    printf("\n");
  }
  return 0;
}