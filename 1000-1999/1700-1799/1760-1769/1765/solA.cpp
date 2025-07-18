#include <bits/stdc++.h>

using namespace std;

using ll = long long;

using ld = long double;

using pii = pair<int, int>;



int main() {

    ios_base::sync_with_stdio(false);

    cin.tie(nullptr);



    int n, m;

    cin >> n >> m;

    vector<string> v(n);

    for (auto& s : v) cin >> s;



    constexpr int MAXN = 500;

    vector<bitset<MAXN>> bs(m);

    for (int i = 0; i < n; ++i) {

        for (int j = 0; j < m; ++j) {

            if (v[i][j] == '1')

                bs[j].set(i);

        }

    }



    vector<vector<int>> g(m);

    for (int i = 0; i < m; ++i) {

        for (int j = 0; j < m; ++j) {

            if (i == j) continue;

            if ((bs[i] == bs[j] && i < j) || (bs[i] != bs[j] && (bs[i] & bs[j]) == bs[j])) {

                g[i].push_back(j);

            }

        }

    }



    vector<int> mt(m, -1);

    vector<char> used;



    function<bool(int)> dfs = [&](int v) {

        if (used[v]) return false;

        used[v] = 1;

        for (int to : g[v]) {

            if (mt[to] == -1 || dfs(mt[to])) {

                mt[to] = v;

                return true;

            }

        }

        return false;

    };



    for (int i = 0; i < m; ++i) {

        used.assign(m, 0);

        dfs(i);

    }



    vector<int> revMt(m, -1);

    for (int i = 0; i < m; ++i) {

        if (mt[i] != -1) {

            revMt[mt[i]] = i;

        }

    }



    vector<vector<int>> paths;

    for (int i = 0; i < m; ++i) {

        if (mt[i] == -1) {

            vector<int> path;

            int cur = i;

            while (cur != -1) {

                path.push_back(cur);

                cur = revMt[cur];

            }

            reverse(path.begin(), path.end());

            paths.push_back(path);

        }

    }



    int k = (int)paths.size();

    vector<int> which(m), access(m);

    vector<vector<int>> matrix(n, vector<int>(k, 1));

    for (int i = 0; i < k; ++i) {

        for (int x : paths[i]) {

            which[x] = i + 1;

            int who = 2;

            for (int j = 0; j < n; ++j)

                who += v[j][x] == '0';

            access[x] = who;

            for (int j = 0; j < n; ++j) {

                if (v[j][x] == '1') {

                    matrix[j][i] = max(matrix[j][i], who);

                }

            }

        }

    }



    cout << k << '\n';

    for (int i = 0; i < m; ++i) cout << which[i] << ' ';

    cout << '\n';

    for (int i = 0; i < m; ++i) cout << access[i] << ' ';

    cout << '\n';

    for (int i = 0; i < n; ++i) {

        for (int x : matrix[i]) cout << x << ' ';

        cout << '\n';

    }

}