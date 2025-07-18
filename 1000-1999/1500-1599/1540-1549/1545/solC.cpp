#include<bits/stdc++.h>
using namespace std;

const int N = 1005, M = 998244353;
int a[N][N];
bool rowdone[N];
bool must[N];
bool done[N][N];
int cnt[N];
vector<int> occ[N];
int vis[N];

vector<int> adj[N];

void AddEdge(int u, int v) {
//    cout<<"   "<<u<<" "<<v<<endl;
    adj[u].push_back(v);
    adj[v].push_back(u);
}

void dfs(int u, bool c) {
    if (vis[u] != -1) {
        assert(vis[u] == c);
        return;
    }
    vis[u] = c;
    if (rowdone[u])    assert(c == must[u]);
    for (int v: adj[u]) {
        dfs(v, !c);
    }
}

int main() {
    ios::sync_with_stdio(0);
    cin.tie(0);

    int t;
    cin>>t;

    while (t--) {
        int n;
        cin>>n;

        vector<int> ans;
        long long ways = 1;

        for (int i=0; i<2*n; i++) {
            for (int j=0; j<n; j++) {
                cin>>a[i][j];
            }
            rowdone[i] = 0;
            must[i] = 0;
        }

        for (int c=0; c<n; c++) {
            for (int v=1; v<=n; v++) {
                done[c][v] = 0;
            }
        }

        for (int it=0; it<n; it++) {

            int col = -1, val = -1;
            for (int c=0; c<n; c++) {
                for (int i=1; i<=n; i++)    cnt[i] = 0;

                for (int r=0; r<2*n; r++)
                    if (!rowdone[r])
                        cnt[a[r][c]]++;

                int found = -1;
                for (int i=1; i<=n; i++) {
                    if (done[c][i])     continue;
                    assert(cnt[i] >= 1);
                    if (cnt[i] == 1) {
                        found = i;
                        break;
                    }
                }

                if (found != -1) {
                    col = c;
                    val = found;
                    break;
                }
            }

            if (col == -1) {
                for (int r=0; r<2*n; r++)   adj[r].clear();
                for (int c=0; c<n; c++) {
                    for (int i=1; i<=n; i++)    cnt[i] = 0, occ[i].clear();

                    for (int r=0; r<2*n; r++)
                        if (!rowdone[r]) {
                            cnt[a[r][c]]++;
                            occ[a[r][c]].push_back(r);
                        }

                    for (int i=1; i<=n; i++) {
                        if (done[c][i])     continue;
                        assert(cnt[i] == 2);
                        AddEdge(occ[i][0], occ[i][1]);
                    }
                }
                for (int r=0; r<2*n; r++)   vis[r] = -1;

                for (int r=0; r<2*n; r++) {
                    if (rowdone[r] && vis[r] == -1) {
                        dfs(r, must[r]);
                    }
                }

                for (int r=0; r<2*n; r++) {
                    if (vis[r] == -1) {
                        dfs(r, 0);
                        ways = (ways*2)%M;
                    }
                }

                ans.clear();
                for (int i=0; i<2*n; i++)
                    if (vis[i] == 1)
                        ans.push_back(i);
                break;
            }
            else {
                int mustrow = -1;
                for (int r=0; r<2*n; r++)
                    if (!rowdone[r] && a[r][col] == val)
                        mustrow = r;

                assert(mustrow != -1);
                ans.push_back(mustrow);
                must[mustrow] = 1;

                for (int i=0; i<n; i++)     done[i][a[mustrow][i]] = 1;

                for (int r=0; r<2*n; r++) {
                    if (rowdone[r]) continue;

                    bool found = false;
                    for (int c=0; c<n; c++) {
                        if (a[r][c] == a[mustrow][c]) {
                            found = true;
                            break;
                        }
                    }
                    if (found)  {
                        rowdone[r] = 1;
                    }
                }
            }
        }

        cout<<ways<<endl;
        for (int x: ans)    cout<<1+x<<" ";
        cout<<endl;
    }
}