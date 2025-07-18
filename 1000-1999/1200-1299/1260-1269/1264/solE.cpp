#include <bits/stdc++.h>
using namespace std;

const int maxn = 5150 + 10;
const int maxg = 100 + 10;
const int INF = 1000000000;

int n;
int s, t;
int rel[maxg][maxg];
int wedge[maxg][maxg];
int indeg[maxn];

int edgeidx(int x, int y) {
  return (2*n-x) * (x-1) / 2 + (y-x) + n;
}

struct Edge {
  int from, to, cap, flow, cost;
  Edge(int from, int to, int cap, int flow, int cost):
    from(from), to(to), cap(cap), flow(flow), cost(cost) {}
};

template <int maxn> struct MCMF {
  int n, m, s, t;
  vector <Edge> edges;
  vector <int> G[maxn];
  int inq[maxn];
  int d[maxn];
  int p[maxn];
  int a[maxn];

  void init(int n) {
    this->n = n;
    for(int i = 0; i < n; i++) G[i].clear();
    edges.clear();
  }

  void AddEdge(int from, int to, int cap, int cost) {
    edges.push_back(Edge(from, to, cap, 0, cost));
    edges.push_back(Edge(to, from, 0, 0, -cost));
    m = edges.size();
    G[from].push_back(m-2);
    G[to].push_back(m-1);
  }

  bool BellmanFord(int s, int t, int& cost) {
    for(int i = 0; i < n; i++) d[i] = INF;
    memset(inq, 0, sizeof(inq));
    d[s] = 0; inq[s] = 1; p[s] = 0; a[s] = INF;

    queue <int> Q;
    Q.push(s);
    while(!Q.empty()) {
      int u = Q.front(); Q.pop();
      inq[u] = 0;
      for(int i = 0; i < G[u].size(); i++) {
        Edge& e = edges[G[u][i]];
        if(e.cap > e.flow && d[e.to] > d[u] + e.cost) {
          d[e.to] = d[u] + e.cost;
          p[e.to] = G[u][i];
          a[e.to] = min(a[u], e.cap - e.flow);
          if(!inq[e.to]) { Q.push(e.to); inq[e.to] = 1; }
        }
      }
    }
    if(d[t] == INF) return false;
    cost += d[t] * a[t];
    int u = t;
    while(u != s) {
      edges[p[u]].flow += a[t];
      edges[p[u]^1].flow -= a[t];
      u = edges[p[u]].from;      
    }
    return true;
  }

  int Mincost(int s, int t) {
    int cost = 0;
    while(BellmanFord(s, t, cost));
    return cost;
  }
};

MCMF <maxn> solver;

int main() {
  scanf("%d", &n);
  s = 0; t = (n+1)*n/2 + n + 1;
  solver.init(t+1);
  
  for(int i = 1; i <= n; i++)
    for(int j = 1; j <= n; j++) {
        rel[i][j]=2;
    }
    int m;
    scanf("%d",&m);
    while(m--)
    {
        int x,y;
        scanf("%d %d",&x,&y);
        rel[x][y]=1;
        rel[y][x]=0;
        indeg[y]++;
    }

  for(int i = 1; i <= n; i++)
    for(int j = i+1; j <= n; j++) {
      if(rel[i][j] < 2) continue;
      solver.AddEdge(s, edgeidx(i, j), 1, 0);
      solver.AddEdge(edgeidx(i, j), i, 1, 0);
      wedge[j][i] = solver.m - 2;
      solver.AddEdge(edgeidx(i, j), j, 1, 0);
      wedge[i][j] = solver.m - 2;
    }

  int down = 0;
  for(int i = 1; i <= n; i++) {
    down += indeg[i] * (indeg[i]-1) / 2;
    for(int j = indeg[i] + 1; j < n; j++)
      solver.AddEdge(i, t, 1, j-1);
  }

  down += solver.Mincost(s, t);

for(int i=1;i<=n;i++)rel[i][i]=0;
  for(int i = 1; i <= n; i++) {
    for(int j = 1; j <= n; j++) {
      if(rel[i][j] < 2) printf("%d", rel[i][j]);
      else printf("%d", solver.edges[wedge[i][j]].flow);
    }
    printf("\n");
  }

  return 0;
}