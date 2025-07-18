/**
 * Created by 5cm/s on 2024/05/20 22:50:08.
 * 诸天神佛，佑我上分！
 **/
#include <bits/stdc++.h>

using namespace std;

#define itr(it) begin(it), end(it)
#define endl '\n'
#define YES() void(cout << "YES\n")
#define NO() void(cout << "NO\n")

void elysia() {
    int n;
    cin >> n;
    string s;
    cin >> s;
    map<char, vector<int>> mp;
    for (int i = 0; i < n; ++i) {
        mp[s[i]].emplace_back(i);
    }
    int x = int(mp['E'].size()) - int(mp['W'].size());
    int y = int(mp['N'].size()) - int(mp['S'].size());
    if (x % 2 || y % 2) return NO();
    string ans(n, 'R');
    if (x || y) {
        if (x > 0) {
            for (int i = 0; i < x / 2; ++i) {
                ans[mp['E'][i]] = 'H';
            }
        } else {
            for (int i = 0; i < abs(x) / 2; ++i) {
                ans[mp['W'][i]] = 'H';
            }
        }
        if (y > 0) {
            for (int i = 0; i < y / 2; ++i) {
                ans[mp['N'][i]] = 'H';
            }
        } else {
            for (int i = 0; i < abs(y) / 2; ++i) {
                ans[mp['S'][i]] = 'H';
            }
        }
    } else {
        if (n == 2) return NO();
        if (mp['E'].size()) {
            ans[mp['E'][0]] = ans[mp['W'][0]] = 'H';
        } else {
            ans[mp['N'][0]] = ans[mp['S'][0]] = 'H';
        }
    }
    cout << ans << endl;
}

int main() {
#ifdef MEGURINE
    freopen("../input.txt", "r", stdin);
    freopen("../output.txt", "w", stdout);
    clock_t start = clock();
#endif
    ios::sync_with_stdio(false);
    cin.tie(nullptr);
    cout.tie(nullptr);
    int T = 1;
    cin >> T;
    cout << fixed << setprecision(12);
    while (T--) elysia();
#ifdef MEGURINE
    cout << "\nRunning Time: " << double(clock() - start) / CLOCKS_PER_SEC * 1000 << "ms" << endl;
#endif
    return 0;
}