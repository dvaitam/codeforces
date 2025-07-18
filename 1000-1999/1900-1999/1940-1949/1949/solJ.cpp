#ifdef ONLINE_JUDGE
#include<bits/stdc++.h>
#else
#include <iostream>
#include <string>
#include <map>
#include <sstream>
#include <assert.h>
#include <set>
#include <vector>
#include <bit>
#include <cstdint>
#include <queue>
#include <stack>
#include <cstring>
#include <utility>
#endif
using namespace std;

#define int long long

using ii = pair<int, int>;
using ull = unsigned long long;
using iii = pair<int, ii>;

const int inf = 0x3f3f3f3f;
const int ms = 105;


int dx[4] = {1, -1, 0, 0};
int dy[4] = {0, 0, 1, -1};

string gs[ms];
string ge[ms];
int dist[ms][ms];
ii pre[ms][ms];
int prey[ms][ms];
int n, m;
stack<ii> st;
queue<ii> toAdd;
queue<ii> toRemove;
int vis[ms][ms];
bool added[ms][ms];

void dfsAddingToSt(int x, int y) {
    if(vis[x][y]) return;
    vis[x][y] = 1;
    for(int d = 0; d < 4; d++) {
        int a = x+dx[d], b = y+dy[d];
        if(a >= 0 && a < n && b >= 0 && b < m && gs[a][b] == '*') {
            dfsAddingToSt(a, b);
        }
    }
    toRemove.emplace(x, y);
}

void dfsAdding2(int x, int y) {
    if(vis[x][y]) return;
    vis[x][y] = 1;
    toAdd.emplace(x, y);
    for(int d = 0; d < 4; d++) {
        int a = x+dx[d], b = y+dy[d];
        if(a >= 0 && a < n && b >= 0 && b < m && ge[a][b] == '*') {
            dfsAdding2(a, b);
        }
    }
}

void moveTo(int ex, int ey) {
    int rex = ex, rey = ey;
    while(pre[ex][ey] != ii(ex, ey)) {
        st.emplace(ex, ey);
       // cout << "to em " << ex+1 << " - " << ey+1 << " e vim de ";
        tie(ex, ey) = pre[ex][ey];
       // cout << ex+1 << " " << ey+1 << endl;
    }
    // currently ex, ey, need to move through toAdd
    dfsAddingToSt(ex, ey); // ex, ey is one from the start
    while(st.size() > 1) {
        toAdd.emplace(st.top()); // eles sao os primeiros q eu insiro
        toRemove.emplace(st.top()); // sao os ultimos q eu removo
       // cout << "filinha " << st.top().first+1 << " - " << st.top().second+1 << endl;
        st.pop();
    }
    memset(vis, 0, sizeof vis);
    dfsAdding2(rex, rey); // rex, rey is one from the end
    memset(vis, 0, sizeof vis);
    vector<pair<ii, ii>> ans;
    //cout << "moving " << ex << " " << ey << " to -> " << rex << " " << rey << endl;
    while(!toAdd.empty()) {
        auto [x, y] = toAdd.front();
        if(gs[x][y] == '*') {
            added[x][y] = 1;
            toAdd.pop();
        } else {
            assert(!toRemove.empty());
            auto [a, b] = toRemove.front();
            toRemove.pop();
            //cout << "trying to add " << x+1 << " " << y+1 << endl;
            //cout << "lets go removing " << a+1 << " " << b+1 << endl;
            if(added[a][b] && ge[a][b] == '*') {
               // cout << "failed to remove " << endl;
                continue;
            }
            toAdd.pop();
            added[x][y] = 1;
            gs[a][b] = '.';
            gs[x][y] = '*';
            ans.emplace_back(ii(a, b), ii(x, y));
        }
    }
    cout << "YES" << endl;
    cout << ans.size() << endl;
    for(auto [a, b] : ans) {
        cout << a.first+1 << " " << a.second+1 << " " << b.first+1 << " " << b.second +1<< '\n';
    }
}

void solve() {
    cin >> n >> m;
    memset(dist, -1, sizeof dist);
    queue<ii> q;
    for(int i = 0; i < n; i++) {
        cin >> gs[i];
        for(int j = 0; j < m; j++) {
            if(gs[i][j] == '*') {
                dist[i][j] = 0;
                q.emplace(i, j);
                pre[i][j] = ii(i, j);
            }
        }
    }
    for(int i = 0; i < n; i++) {
        cin >> ge[i];
    }
    while(!q.empty()) {
        auto [x, y] = q.front();
        q.pop();
        if(ge[x][y] == '*') {
            moveTo(x, y);
            return;
        }
        for(int d = 0; d < 4; d++) {
            int a = x+dx[d];
            int b = y+dy[d];
            if(a >= 0 && a < n && b >= 0 && b < m && dist[a][b] == -1 && gs[a][b] != 'X') {
                dist[a][b] = dist[x][y]+1;
                pre[a][b] = ii(x, y);
                q.emplace(a, b);
            }
        }
    }
    cout << "NO" << endl;
}

int32_t main() {
    cin.tie(0); ios::sync_with_stdio(0);
    int t = 1;
    //cin >> t;
    while(t--) {
        solve();
    }
}