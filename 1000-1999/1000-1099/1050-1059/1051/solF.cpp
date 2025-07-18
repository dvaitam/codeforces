#include <bits/stdc++.h>

using namespace std;

struct edge {
  int u, v;
  long long w;
  edge () {}
  edge (int _u, int _v, long long _w) {
    u = _u, v = _v, w = _w;
  }
};

const int C = 45;
const int LG = 20;
const int N = 100005;
const long long INF = 1e17 + 5;

bitset <N> inside;
vector <edge> g[N], bad;
int t, cs, n, m, q, p[N][LG], h[N];
long long d[N], dist[C][N];

int find (int x) {
  return x == p[x][0] ? x : p[x][0] = find(p[x][0]);
}

void dfs (int u, int from = -1, int depth = 0, long long far = 0) {
  d[u] = far;
  h[u] = depth;
  p[u][0] = from;
  for (edge e : g[u]) {
    int v = e.v;
    long long w = e.w;
    if (v == from) continue;
    dfs(v, u, depth + 1, far + w);
  }
}

void spfa (int s, long long _d[]) {
  for (int i = 1; i <= n; ++i) {
    _d[i] = INF;
  }
  queue <int> q;
  inside.reset();
  _d[s] = 0, q.push(s), inside[s] = 1;
  while (not q.empty()) {
    int u = q.front();
    q.pop(), inside[u] = 0;
    for (edge e : g[u]) {
      int v = e.v;
      long long w = e.w;
      if (_d[u] + w < _d[v]) {
        _d[v] = _d[u] + w;
        if (!inside[v]) {
          inside[v] = 1, q.push(v);
        }
      }
    }
  }
}

int lca (int u, int v) {
  if (h[u] < h[v]) swap(u, v);
  for (int i = LG - 1; i >= 0; --i) {
    if (h[u] - (1 << i) >= h[v]) {
      u = p[u][i];
    }
  }
  if (u == v) return u;
  for (int i = LG - 1; i >= 0; --i) {
    if (p[u][i] != -1 and p[u][i] != p[v][i]) {
      u = p[u][i], v = p[v][i];
    }
  }
  return p[u][0];
}

long long treeDist (int u, int v) {
  return d[u] + d[v] - 2 * d[lca(u, v)];
}

int main() {
  scanf("%d %d", &n, &m);
  for (int i = 1; i <= n; ++i) {
    p[i][0] = i;
  }
  while (m--) {
    int u, v;
    long long w;
    scanf("%d %d %lld", &u, &v, &w);
    if (find(u) == find(v)) {
      bad.emplace_back(u, v, w);
    } else {
      p[find(u)][0] = find(v);
      g[u].emplace_back(u, v, w);
      g[v].emplace_back(v, u, w);
    }
  }
  for (int i = 1; i <= n; ++i) {
    for (int j = 0; j < LG; ++j) {
      p[i][j] = -1;
    }
  }
  dfs(1);
  for (int j = 1; j < LG; ++j) {
    for (int i = 1; i <= n; ++i) {
      if (p[i][j - 1] != -1) {
        p[i][j] = p[p[i][j - 1]][j - 1];
      }
    }
  }
  for (edge e : bad) {
    g[e.u].emplace_back(e.u, e.v, e.w);
    g[e.v].emplace_back(e.v, e.u, e.w);
  }
  for (int i = 0; i < bad.size(); ++i) {
    spfa(bad[i].u, dist[i << 1]);
    spfa(bad[i].v, dist[i << 1 | 1]);
  }
  scanf("%d", &q);
  while (q--) {
    int u, v;
    scanf("%d %d", &u, &v);
    long long ans = treeDist(u, v);
    for (int i = 0; i < bad.size(); ++i) {
      ans = min(ans, dist[i << 1][u] + bad[i].w + dist[i << 1 | 1][v]);
      ans = min(ans, dist[i << 1][v] + bad[i].w + dist[i << 1 | 1][u]);
    }
    printf("%lld\n", ans);
  }
  return 0;
}