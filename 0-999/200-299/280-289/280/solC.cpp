#include <queue>
#include <cstdio>
#include <algorithm>

const int MAXN = 100001;
const long double ONE = 1.;
struct Edge
{
  int v;
  Edge *next;
} g[MAXN*2], *header[MAXN];
void AddEdge(const int x, const int y)
{
  static int LinkSize = 0;
  Edge* const node = g+(LinkSize++);
  node->v = y;
  node->next = header[x];
  header[x] = node;
}
int n, parent[MAXN], f[MAXN];
void BFS()
{
  std::queue<int> Q;
  for (Q.push(1); !Q.empty(); Q.pop())
  {
    const int u = Q.front();
    f[u] = f[parent[u]] + 1;
    for (Edge *e = header[u]; e; e = e->next)
      if (e->v != parent[u])
        parent[e->v] = u, Q.push(e->v);
  }
}
int main()
{
  scanf("%d", &n);
  for (int i = 0, x, y; i < n; ++i)
  {
    scanf("%d%d", &x, &y);
    AddEdge(x, y);
    AddEdge(y, x);
  }
  BFS();
  long double ans = 0;
  for (int i = 1; i <= n; ++i)
    ans += ONE / f[i];
  printf("%.10f", static_cast<double>(ans));
}