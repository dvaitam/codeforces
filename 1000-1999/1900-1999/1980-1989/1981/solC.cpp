#include<bits/stdc++.h>

using namespace std;

#define int long long
int const maxn = 2e5 + 5;
int a[maxn], b[maxn], inf = 1e9;

int get(int x) {
    for (int i = 30; i >= 0; i--) {
        if ((x>>i)&1) return i;
    }
    return -1;
}

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
        int n;
        cin >> n;
        for (int i = 1; i <= n; i++) cin >> a[i];
        int flag = 1, l = 1;
        while (l <= n) {
            if (a[l] != -1) {
                b[l] = a[l];
                l++;
            } else {
                int r = l;
                while (r <= n && a[r] == -1) r++;
                if (l == 1) {
                    b[r] = a[r];
                    if (r == n + 1) b[r] = 1;
                    for (int i = r - 1; i >= 1; i--) {
                        if (2 * b[i + 1] <= inf)
                            b[i] = 2 * b[i + 1];
                        else
                            b[i] = b[i + 1] / 2;
                    }
                } else if (r == n + 1) {
                    for (int i = l; i < r; i++) {
                        if (b[i - 1] * 2 <= inf) b[i] = 2 * b[i - 1];
                        else b[i] = b[i - 1] / 2;
                    }
                } else {
                    int L = a[l - 1], R = a[r], can = 0;
                    for (int suff = 0; suff <= min(30ll, r - l + 1); suff++) {
                        int cur = L / (1ll << suff);
                        if (!cur) continue;
                        int tmp = r - l - suff + 1;
                        int v1 = get(cur), v2 = get(R);
                        //if (suff == 1) cout << v1 << " " << v2 << " " << cur << " " << R << endl;
                        if (v1 > v2) continue;
                        int value = (R>>(v2 - v1));
                        int need = v2 - v1;
                        if (cur == value) {
                            if (need <= tmp && (tmp - need) % 2 == 0) {
                                can = 1;
                                for (int i = l; i < l + suff; i++) {
                                    b[i] = b[i - 1] / 2;
                                }
                                int mask = (R&((1ll << need) - 1)), x = need;
                                for (int i = l + suff; i < l + suff + need; i++) {
                                    x--;
                                    b[i] = b[i - 1] * 2 + ((R>>x)&1);
                                }
                                for (int i = l + suff + need; i < r; i++) {
                                    if (i % 2 == (l + suff + need) % 2) b[i] = b[i - 1] * 2;
                                    else b[i] = b[i - 1] / 2;
                                }
                                //cout << suff << " " << v1 << " " << v2 << " " << cur << " " << R << endl;
                                break;
                            }
                        }
                    }
                    flag &= can;
                }
                l = r;
            }
        }
        for (int i = 1; i < n; i++) {
            if (b[i] != b[i + 1] / 2 && b[i + 1] != b[i] / 2) flag = 0;
        }
        if (!flag) cout << -1 << '\n';
        else {
            for (int i = 1; i <= n; i++) cout << b[i] << " ";
            cout << '\n';
        }
    }
    return 0;
}