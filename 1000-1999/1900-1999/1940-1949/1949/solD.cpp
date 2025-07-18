#include <bits/stdc++.h>
using namespace std;

#define rep(i, a, b) for (int i = a; i < (b); ++i)
#define trav(a, x) for (auto& a : x)
#define all(x) x.begin(), x.end()
#define sz(x) (int)(x).size()
typedef long long ll;
typedef pair<int, int> pii;
typedef vector<int> vi;

int main() {
    ios::sync_with_stdio(false);
    cin.tie(NULL);

    int n;
    cin >> n;
    vector<string> str(n);
    vector<vi> c(2, vi(n));
    vector<pii> e[2];
    string s = "FS";
    rep(i, 0, n) {
        cin >> str[i];
        rep(j, 0, n) {
            if(str[i][j] == 'F') {
                c[0][i] |= 1 << j;
                if(i < j) e[0].push_back({i, j});
            }
            else if(str[i][j] == 'S') {
                c[1][i] |= 1 << j;
                if(i < j) e[1].push_back({i, j});
            }
        }
    }
    int full = ((1 << n) - 1);
    int cnt = 0;
    rep(mask, 0, 1 << n) {
        if(__builtin_popcount(mask) - 1 > (3 * n + 3) / 4) continue;
        if(__builtin_popcount(mask ^ full) * 2 > (3 * n + 3) / 4) continue;
        cnt++;
        rep(x, 0, 2) {
            bool bad = false;
            int c0 = 0;
            for (auto p : e[x]) {
                int a = (mask >> p.first) & 1, b = (mask >> p.second) & 1;
                if(a ^ b) {
                    bad = true;
                    break;
                }
            }
            if(bad) continue;
            for (auto p : e[!x]) {
                int a = (mask >> p.first) & 1, b = (mask >> p.second) & 1;
                if(a && b) c0++;
            }

            if(__builtin_popcount(mask ^ full) * 2 + c0 > (3 * n + 3) / 4) continue;
//            cout << x << endl;
//            cout << bitset<5>(mask) << endl;
//            cout << mask << endl;
            rep(i, 0, n) {
                rep(j, i + 1, n) {
                    if(str[i][j] == '?') {
                        int a = (mask >> i) & 1, b = (mask >> j) & 1;
                        if(a && b) str[i][j] = str[j][i] = s[x];
                        else if(a ^ b) str[i][j] = str[j][i] = s[!x];
                        else str[i][j] = str[j][i] = 'F';
                    }
                }
            }
            rep(i, 0, n) cout << str[i] << '\n';

            return 0;
        }
    }
    assert(false);
//    cout << cnt << endl;
}