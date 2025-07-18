#include <bits/stdc++.h>

#define FAST ios::sync_with_stdio(0); cin.tie(0)
#define FILEIN freopen("input.txt", "r", stdin)
#define FILEOUT freopen("output.txt", "w", stdout)
#define SP << " "
#define endl "\n"
#define PI (3.1415926535897932384626433832795)
#define INF (2147483647)
#define LLINF (9223372036854775807LL)
#define EPS (1e-9)

#define ALL(a) a.begin(), a.end()
#define ALLR(a) a.rbegin(), a.rend()
#define TYPFOR(i, n) for (int i = 0; i < n; ++i)

typedef long long LL;
typedef long double LD;

using namespace std;

int main() {
    FAST;
    int xta, xtb, xtc, yta, ytb, ytc;
    int xa, xb, xc, ya, yb, yc;
    vector<pair<int, int>> x(3);
    for (int i = 0; i < 3; ++i) {
        int a1, a2;
        cin >> a1 >> a2;
        x[i].first = a1;
        x[i].second = a2;
    }
    sort(ALL(x));
    xa = x[0].first;
    xb = x[1].first;
    xc = x[2].first;
    ya = x[0].second;
    yb = x[1].second;
    yc = x[2].second;
    vector<pair<int, int>> ans(0);

    for (int i = xa; i < xb; ++i) {
        ans.push_back({i, ya});
    }

    for (int i = min(min(yc, yb), ya); i <= max(max(yc, yb), ya); ++i) {
        ans.push_back({xb, i});
    }

    for (int i = xc; i > xb; --i) {
        ans.push_back({i, yc});
    }

    cout << ans.size() << endl;

    for (int i = 0; i < ans.size(); ++i) {
        cout << ans[i].first SP << ans[i].second << endl;
    }




    return 0;
}