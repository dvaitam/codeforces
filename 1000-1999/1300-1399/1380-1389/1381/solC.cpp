#include <bits/stdc++.h>
using namespace std;
using ll = long long;
#define debug(x) cerr << "[" << __LINE__ << ' ' << #x << "]: " << (x) << endl

int main() {
    cin.tie(0)->sync_with_stdio(0);

    int T;
    cin >> T;
    while (T--) {
        int n, x, y;
        cin >> n >> x >> y;

        int a[n];
        vector<bool> c(n+2);
        for (int& e : a) cin >> e, c[e] = true;
        
        int unused = 0;
        for (int i=1; i<n+2; i++) if (!c[i]) unused = i;

        y -= x;
        vector<pair<int, int>> b;
        vector<int> cnt(n+2);
        vector<bool> used(n);
        for (int i=0; i<n && (int)b.size() < y; i++) {
            if (cnt[a[i]]+1 <= y/2) b.emplace_back(a[i], i), cnt[a[i]]++, used[i] = true;
        }

        if (y - b.size() == 1 && y % 2 == 1){
            int xx = -1, yy = -1;
            for (int i=0; i<n; i++) if (!used[i]) {
                xx = i;
                used[i] = true;
                break;
            }
            for (int i=0; i<n; i++) if (!used[i] && xx != -1 && a[xx] != a[i]) {
                yy = i;
                used[i] = true;
                break;
            }
            if (yy == -1) {
                cout << "NO\n";
                continue;
            }
            b.emplace_back(a[xx], xx);
            b.emplace_back(unused, yy);
        }

        if ((int)b.size() < y || n-(int)b.size() < x) {
            cout << "NO\n";
            continue;
        }

        /*for (auto [idk, idk2] : b) {
            debug(idk);
            debug(idk2);
        }*/

        vector<int> ans(n, -1);

        sort(b.begin(), b.end());
        int half = b.size() / 2;
        for (int i=0; i<(int)b.size(); i++) {
            ans[b[i].second] = b[(i+half)%b.size()].first;
        }
        for (int i=0; i<n; i++) {
            if (ans[i] == -1 && x > 0) x--, ans[i] = a[i];
            else if (ans[i] == -1) ans[i] = unused;
        }

        cout << "YES\n";
        for (int e : ans) cout << e << ' ';
        cout << '\n';
    }
}