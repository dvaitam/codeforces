#include <bits/stdc++.h>
#include<ext/pb_ds/assoc_container.hpp>
#include<ext/pb_ds/tree_policy.hpp>

using namespace std;
using namespace __gnu_pbds;
// find_by_order, order_of_key
typedef tree<int, null_type, less<int>, rb_tree_tag, tree_order_statistics_node_update> oset;
typedef tree<int, null_type, less_equal<int>, rb_tree_tag, tree_order_statistics_node_update> omset;
#define int long long
#define endl '\n'
#define pi pair<int,int>
#define adjs(name, type, size) vector<vector<type>>name(size)
#define adjpass(name, type) vector<vector<type>>&name
#define rest(name, val) memset(name,val,sizeof(name))
#define all(x) x.begin(),x.end()
#define killua ios_base::sync_with_stdio(false), cin.tie(NULL), cout.tie(0)
//changes in dir:
int dx[] = {0, 0, 1, -1, -1, 1, -1, 1};
int dy[] = {1, -1, 0, 0, 1, 1, -1, -1};
int cases = 01;

/***إِلا أَنْ يَشَاءَ اللَّهُ وَاذْكُرْ رَبَّكَ إِذَا نَسِيتَ وَقُلْ عَسَى أَنْ يَهْدِيَنِي رَبِّي لأَقْرَبَ مِنْ هَذَا رَشَدًا ***/
class DSU {
private:
    vector<int> rank, par, size;
public:
    DSU(int n = 0) {
        rank.resize(n + 1);
        par.resize(n + 1);
        size.resize(n + 1, 1);
        for (int i = 0; i <= n; i++)
            par[i] = i;
    }

    int findUpar(int node) {
        if (node == par[node])
            return node;
        return par[node] = findUpar(par[node]);//path comp
    }

    void unionbyRank(int u, int v) {
        int UparU = findUpar(u);
        int UparV = findUpar(v);
        if (UparU == UparV)
            return;
        if (rank[UparU] < rank[UparV]) {
            par[UparU] = UparV;
        } else {
            par[UparV] = UparU;
            rank[UparU] += (rank[UparU] == rank[UparV]);
        }
    }

    void unionbySize(int u, int v) {
        int UparU = findUpar(u);
        int UparV = findUpar(v);
        if (UparU == UparV)
            return;
        if (size[UparU] < size[UparV]) {
            par[UparU] = UparV;
            size[UparV] += size[UparU];
        } else {
            par[UparV] = UparU;
            size[UparU] += size[UparV];
        }
    }
};


void gon() {
    int n, m;
    cin >> n >> m;
    int mi = 1e9;
    vector<vector<pi>> g(n + 5);
    vector<pair<int, pi>> v(m), non;
    for (auto &i: v) cin >> i.second.first >> i.second.second >> i.first, i.first *= -1;
    std::sort(v.begin(), v.end());
    DSU d(n);
    vector<pi > trys;
    int cost = 1e9;
    int st, en;
    for (auto i: v) {
        g[i.second.first].push_back({i.second.second, i.first});
        g[i.second.second].push_back({i.second.first, i.first});
        if (d.findUpar(i.second.first) != d.findUpar(i.second.second)) {
        } else {
            st = i.second.first;
            en = i.second.second;
            cost = abs(i.first);
        }
        d.unionbySize(i.second.first, i.second.second);
    }
    vector<int> vis(n + 5), finalnodes;
    vector<int> dads(n + 5, -1);
    queue<int> q;
    q.push(st);
    vis[st] = 1;
    while (q.size()) {
        int node = q.front();
        q.pop();
        if (node == en)
            break;
        for (auto i: g[node]) {
            if (i.first == en && node == st)
                continue;
            if (!vis[i.first]) {
                q.push(i.first);
                vis[i.first] = 1;
                dads[i.first] = node;
            }
        }
    }
    int cur = en;
    while (1) {
        finalnodes.push_back(cur);
        if (dads[cur] == -1) {
            break;
        }
        cur = dads[cur];
    }
    cout << cost << " " << finalnodes.size() << endl;
    for (auto i: finalnodes)
        cout << i << " ";
    cout << endl;
}

int32_t main() {
#ifndef ONLINE_JUDGE
    freopen("Input.txt", "r", stdin);
    freopen("Output.txt", "w", stdout);
#endif
    killua;
    int t = 1;
    if (cases) cin >> t;
    while (t--) {
        gon();
    }
}