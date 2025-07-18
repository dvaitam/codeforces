#include <bits/stdc++.h>

using namespace std;

int n,m;

vector<int>G[105];

long long dist1[105], dist2[105], memo1[105], memo2[105];



void bfs(int from, int to, long long dist[], long long memo[]) {

    for(int i=1; i<=n; i++) dist[i]=1e6;

    queue<pair<int,int> > q; 

    q.push(make_pair(0,from));

    memo[from]=1;

    while (!q.empty()) {

        int d = q.front().first, u = q.front().second; q.pop(); 

        dist[u] = d; 

        

        for(int i=0; i<(int)G[u].size(); i++) {

            int v = G[u][i]; 

            

            if (dist[v]+1 == dist[u]) memo[u]+=memo[v];

            if (dist[v] > d+1) {

                dist[v] = d+1; 

                q.push(make_pair(dist[v], v));

            }

        }

    }

}

int main() {

    scanf("%d%d", &n,&m);

    int x,y;

    for(int i=0; i<m;i++) scanf("%d%d",&x,&y), G[x].push_back(y), G[y].push_back(x);

    bfs(1, n, dist1, memo1);

    bfs(n, 1, dist2, memo2);

    

    double ans = 0;

    for(int i=1; i<=n; i++) {

        if (dist1[i]+dist2[i]!=dist1[n]) continue;

        long long safe = memo1[i] * memo2[i];

        long long total = memo1[n];

        double cur = (double)safe / total;

        if (i != n && i != 1) cur*=2;

        ans = max(ans,cur);

    }

    printf("%.9f\n", ans);

    return 0;

}