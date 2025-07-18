#include <cstdio>
#include <cstdlib>
#include <algorithm>
#include <queue>
#include <vector>

using namespace std;

const int MAXN = 100009;

int t[MAXN][3], bij[MAXN][2];

bool vis[MAXN];
int d[MAXN], f[MAXN];
vector<int> g[MAXN];

int q[MAXN], p[MAXN];

int find(int u) {
    int v = u;
    while (f[v] != v) v = f[v];
    while (f[u] != u) {
        int w = f[u];
        f[u] = v;
        u = w;
    }
    return u;
}

void add(int u, int v) {
    for (int i = 0; i < 2; ++i)
        if (bij[u][i] == 0) {
            bij[u][i] = v;
            break;
        }
}

void connect(int u, int v) {
//printf("%d, %d\n", u, v);
    if (find(u) == find(v)) return ;
//printf("done\n");
    add(u, v);
    add(v, u);
    f[find(v)] = u;
}

int main() {
    int Cases;
    scanf("%d", &Cases);
    while (Cases--) {
        int n;
        scanf("%d", &n);
        for (int i = 1; i <= n; ++i) {
            g[i].clear();
            d[i] = 0;
            bij[i][0] = bij[i][1] = 0;
            vis[i] = 0;
            f[i] = i;
        }
        for (int i = 1; i <= n - 2; ++i) {
            for (int j = 0; j < 3; ++j) {
                scanf("%d", &t[i][j]);
                g[t[i][j]].push_back(i);
                ++d[t[i][j]];
            }
        }
        queue<int> que;
        for (int i = 1; i <= n; ++i)
            if (d[i] == 1) que.push(i);
        int qnt = 0;
        while (!que.empty()) {
            int u = que.front();
            que.pop();
            for (int i = 0; i < (int)g[u].size(); ++i) {
                int w = g[u][i];
                if (!vis[w]) {
                    q[++qnt] = w;
                    for (int j = 0; j < 3; ++j)
                        if (t[w][j] != u) {
                            connect(u, t[w][j]);
                            if (--d[t[w][j]] == 1) que.push(t[w][j]);
                        }
                    vis[w] = 1;
                }
            }
        }
        for (int i = 1; i <= n; ++i) vis[i] = 0;
        for (int i = 1; i <= n; ++i)
            if (bij[i][1] == 0) {
                p[1] = i;
                break;
            }
        vis[p[1]] = vis[0] = 1;
        for (int i = 1, k = 1; i <= n; ++i) {
            for (int j = 0; j < 2; ++j)
                if (!vis[bij[p[i]][j]]) {
                    vis[bij[p[i]][j]] = 1;
                    p[++k] = bij[p[i]][j];
                }
        }
        for (int i = 1; i <= n; ++i)
            printf("%d%c", p[i], i == n? '\n': ' ');
        for (int i = 1; i <= n - 2; ++i)
            printf("%d%c", q[i], i == n - 2? '\n': ' ');
    }
    return 0;
}