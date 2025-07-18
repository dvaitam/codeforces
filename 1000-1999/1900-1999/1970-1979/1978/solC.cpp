/**
 _  _   __  _ _ _  _  _ _
 |a  ||t  ||o    d | |o  |
| __    _| | _ | __|  _ |
| __ |/_  | __  /__\ / _\|

**/

#include <bits/stdc++.h>

using namespace std;

typedef long long ll;



int main () {
    ios_base::sync_with_stdio(false);
    cin.tie(0);

    int t;
    cin >> t;
    while (t--) {
        int n; ll k;
        cin >> n >> k;
        if (k % 2 != 0) {
            cout << "No\n";
            continue;
        }
        int p[n + 2];
        iota(p + 1, p + n + 1, 1);
        for (int i = 1; i < n - i + 1; i++) {
            int j = n - i + 1;
            if ((j - i) * 2 < k) {
                swap(p[i], p[j]); k -= (j - i) * 2;
            } else {
                swap(p[i], p[i + k / 2]);
                k = 0;
                break;
            }
        }
        if (k == 0) {
            cout << "Yes\n";
            for (int i = 1; i <= n; i++) {
                cout << p[i] << " ";
            }
            cout << "\n";
        } else {
            cout << "No\n";
        }
    }

    return 0;
}