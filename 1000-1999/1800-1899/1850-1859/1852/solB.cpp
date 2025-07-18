#include<bits/stdc++.h>

using namespace std;

int const maxn = 1e5 + 5;
int a[maxn];

main() {
#ifdef HOME
    freopen("input.txt", "r", stdin);
#endif // HOME
    ios_base::sync_with_stdio(0);
    cin.tie(0);
    cout.tie(0);
    int t;
    cin >> t;
    while (t--) {
        int n, x;
        cin >> n;
        vector < pair < int, int > > b;
        for (int i = 1; i <= n; i++) {
            cin >> x;
            b.push_back({x, i});
        }
        sort(b.begin(), b.end());
        int l = 0, r = b.size() - 1, value = n, cnt_p = 0, cnt_m = 0, flag = 1;
        while (l <= r) {
            if (b[l].first == cnt_p) {
                int cur = l;
                while (cur <= r && b[l].first == b[cur].first) {
                    a[b[l].second] = -value, cnt_m++, l++;
                }
                value--;
            } else if (b[r].first == n - cnt_m) {
                int cur = r;
                while (cur >= l && b[r].first == b[cur].first) {
                    a[b[r].second] = value, cnt_p++, r--;
                }
                value--;
            } else {
                flag = 0;
                break;
            }
        }
        if (!flag) {
            cout << "NO\n";
            continue;
        }
        cout << "YES\n";
        for (int i = 1; i <= n; i++) cout << a[i] << " ";
        cout << '\n';
    }
    return 0;
}