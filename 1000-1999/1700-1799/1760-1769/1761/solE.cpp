#include "bits/stdc++.h"

using namespace std;

using i64 = long long;

//constexpr int mod = 998244353;

//constexpr int mod = 1000000007;

class UnionFind {  //并查集板子

public:

    int n;

    vector<int> f, sz;

    UnionFind(const int &n) : n(n), f(n + 1), sz(n + 1, 1) { 

        iota(f.begin(), f.end(), 0); 

    }

    int findfa(int x) {

        if(f[x] != x)

            f[x] = findfa(f[x]);

    	return f[x];

    }

    bool unite(int x, int y) {

        int fx = findfa(x), fy = findfa(y);

        if(fx == fy)

            return false;

        if(sz[fx] > sz[fy])

            swap(fx, fy);

        f[fx] = fy;

        sz[fy] += sz[fx];

        return true;

    }

    bool check(int x, int y) {

        return findfa(x) == findfa(y);

    }

    int& operator [](int x) {

        return f[x];

    }

};

void solve()

{

    int n;

    cin >> n;

    UnionFind u(n);

    vector<string> s(n + 1);

    vector<int> d(n + 1);

    for(int i = 1; i <= n; i ++) {

        cin >> s[i];

    }

    for(int i = 1; i <= n; i ++) {

        for(int j = 0; j < n; j ++) {

            if(s[i][j] == '1')

                u.unite(i, j + 1), d[i]++;

        }

    }

    if(u.sz[u.findfa(1)] == n) {

        cout << 0 << '\n';

        return;

    }

    vector<vector<int> > v(n + 1);

    vector<int> roots;

    for(int i = 1; i <= n; i ++) {

        v[u.findfa(i)].push_back(i);

        if(u.findfa(i) == i) {

            roots.push_back(i);

            if(u.sz[i] == 1) {

                cout << 1 << '\n' << i << '\n';

                return;

            }

        }

    }

    for(int root : roots) {

        for(int uu : v[root]) {

            if(d[uu] != u.sz[root] - 1){

                int mn = -1;

                for(int x : v[root])

                    if(mn == -1 || d[mn] > d[x])

                        mn = x;

                cout << 1 << '\n' << mn << '\n';

                return;

            }

        }

    }

    if(roots.size() == 2) {

        int u;

        if (v[roots[0]].size() <= v[roots[1]].size()) u = roots[0];

        else u = roots[1];

        cout << v[u].size() << '\n';

        for(auto x : v[u]) cout << x << ' ';

        cout << '\n';

        return;

    }

    cout << 2 << '\n';

    cout << v[roots[0]][0] << ' ' << v[roots[1]][0] << '\n';

}



signed main()

{

#ifdef Parry

    //freopen("in.txt", "r", stdin);

    //freopen("out.txt", "w", stdout);

#endif

    ios::sync_with_stdio(0); cin.tie(0);



    int t = 1;

    cin >> t;

    while (t --) {

        solve();

    }



    return 0;

}