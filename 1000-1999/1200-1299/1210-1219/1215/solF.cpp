#include <cstdio>
#include <stack>
#include <algorithm>
const int maxn = 800010;

int n, m, p, M, x, y;
int l[maxn], r[maxn];

struct node {
  int l, r, id;
} a[maxn];

bool cmpl(const node &a, const node &b) { return a.l < b.l; }
bool cmpr(const node &a, const node &b) { return a.r < b.r; }

struct Graph {
 private:
  static const int maxn = 4000010, maxm = 5000010;
  int fi[maxn], v[maxm], nxt[maxm], cnt;
  int vs[maxn], vis[maxn], low[maxn], nm, c[maxn], bh;
  std::stack<int> stk;
 public:
  void push(int u, int v) { this->v[++cnt] = v, nxt[cnt] = fi[u], fi[u] = cnt; }
  void merge(int nw){
    nm++;
    while (stk.top() != nw) {
      c[stk.top()] = nm, vs[stk.top()] = 0; stk.pop();
    }
    c[nw] = nm, vs[nw] = 0, stk.pop();
  }
  void dfs(int nw){
    stk.push(nw), vs[nw] = 1;
    vis[nw] = low[nw] = ++bh;
    for (int i = fi[nw]; i; i = nxt[i]) {
      if (!vis[v[i]]) dfs(v[i]), low[nw] = std::min(low[nw], low[v[i]]);
      else if (vs[v[i]]) low[nw] = std::min(low[nw], vis[v[i]]);
    }
    if (low[nw] == vis[nw])  merge(nw);
  }
  void solve() {
    bh = nm = 0;
    for (int i = 1; i <= n << 1; ++i) vs[i] = vis[i] = low[i] = 0;
    for (int i = 1; i <= n << 1; ++i) if (!vis[i]) dfs(i);
    int K = 0, L = 1, R = M;
    for (int i = 1; i <= n; ++i) {
      if (c[i] == c[n + i]) return (void) puts("-1");
      if (c[i] < c[n + i]) K++, L = std::max(l[i], L), R = std::min(r[i], R);
    }
    printf("%d %d\n", K, L);
    for (int i = 1; i <= n; ++i) if (c[i] < c[n + i]) printf("%d ", i);
    puts("");
  }
} graph;

int main() {
  scanf("%d%d%d%d", &p, &n, &M, &m);
  for (int i = 1; i <= p; ++i) {
    scanf("%d%d", &x, &y);
    graph.push(x + n, y);
    graph.push(y + n, x);
  }
  for (int i = 1; i <= n; ++i) {
    scanf("%d%d", &a[i].l, &a[i].r), a[i].id = i;
    l[i] = a[i].l, r[i] = a[i].r;
  }
  std::sort(a + 1, a + 1 + n, cmpl);
  for (int i = M; i; --i) if (i != M) graph.push(i + n + n, i + n + n + 1);
  for (int i = n; i; --i) graph.push(a[i].l + n + n, a[i].id + n);
  std::sort(a + 1, a + 1 + n, cmpr);
  for (int i = 1; i <= M; ++i) if (i != 1) graph.push(i + M + n + n, i + M + n + n - 1);
  for (int i = 1; i <= n; ++i) graph.push(a[i].r + n + n + M, a[i].id + n);
  for (int i = 1; i <= n; ++i) {
    if (a[i].r != M) graph.push(a[i].id, a[i].r + n + n + 1);
    if (a[i].l != 1) graph.push(a[i].id, a[i].l + n + n - 1 + M);
  }
  for (int i = 1; i <= m; ++i) {
    scanf("%d%d", &x, &y);
    graph.push(x, y + n);
    graph.push(y, x + n);
  }
  graph.solve();
  return 0;
}