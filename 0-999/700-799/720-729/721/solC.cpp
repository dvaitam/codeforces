#include <bits/stdc++.h>



#define MAXN 5010



using namespace std;



struct Edge {

    int v, w, next;

}edge[MAXN];

int e, ft[MAXN];

int n, m, T;

map<int, int> dp[MAXN];

bool flag[MAXN];



void add_edge(int u, int v, int w)

{

    edge[e].v = v;

    edge[e].w = w;

    edge[e].next = ft[u];

    ft[u] = e++;

}



void dfs(int u)

{

    if (flag[u])

        return;

    dp[u].clear();

    flag[u] = true;

    if (u == n) {

        dp[u][1] = 0;

        return;

    }

    for (int i = ft[u]; i != -1; i = edge[i].next) {

        int v = edge[i].v;

        int w = edge[i].w;

        dfs(v);

        for (auto pr : dp[v]) {

            int idx = pr.first;

            int value = pr.second;

            int x = value + w;

            if (x > T)

                continue;

            if (dp[u].find(idx + 1) == dp[u].end())

                dp[u][idx + 1] = x;

            else

                dp[u][idx + 1] = min(dp[u][idx + 1], x);

        }

    }

}



int main()

{

    while (cin >> n >> m >> T) {

        e = 0;

        memset(ft, -1, sizeof(ft));

        int u, v, w;

        for (int i = 0; i < m; i++) {

            scanf("%d%d%d", &u, &v, &w);

            add_edge(u, v, w);

        }

        for (int i = 1; i <= n; i++) {

            flag[i] = false;

            dp[i].clear();

        }

        dfs(1);

        int ans = dp[1].begin() -> first;

        u = 1;

        for (auto pr : dp[1]) {

            ans = max(ans, pr.first);

        }

        printf("%d\n", ans);

        while (u != n){

            printf("%d ", u);

            for (int i = ft[u]; i != -1; i = edge[i].next) {

                v = edge[i].v;

                w = edge[i].w;

                if (dp[v].find(ans - 1) != dp[v].end() && w + dp[v][ans - 1] == dp[u][ans]) {

                    u = v;

                    ans--;

                    break;

                }

            }

        };

        printf("%d\n", u);

    }

    return 0;

}