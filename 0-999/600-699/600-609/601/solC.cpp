#include <cstdio>
#include <cstring>
#include <algorithm>
#include <vector>
#include <iostream>
#include <cassert>
#include <cmath>
#include <string>
#include <queue>
#include <set>
#include <map>
#include <cstdlib>

using namespace std;

#define TASKNAME "C"

void solve(int test_number);

int main() {
    ios_base::sync_with_stdio(false);
    cin.tie(nullptr);
    cout.setf(ios::fixed);
    cout.precision(9);
    cerr.setf(ios::fixed);
    cerr.precision(3);
#ifdef LOCAL
    freopen("test.in", "r", stdin);
    freopen("test.out", "w", stdout);
#else
#endif
    int n = 1;
    for (int i = 0; i < n; i++) {
        solve(i);
    }
}

const int MAX_N = 105;
const int MAX_M = 1005;
const int MAX_V = MAX_N * MAX_M;

double p[MAX_N][MAX_V];

void solve(int test_number) {
    int n, m, r[MAX_N], s = 0;
    cin >> n >> m;
    if (m == 1) {
        cout << 1.0 << endl;
        return;
    }
    for (int i = 1; i <= n; i++) {
        cin >> r[i];
        s += r[i];
    }
    p[0][0] = 1.0;
    for (int i = 1; i <= n; i++) {
        double sum = p[i - 1][0];
        for (int j = 1; j <= i * m; j++) {
            if (j - m - 1 >= 0) {
                sum -= p[i - 1][j - m - 1];
            }
            p[i][j] = sum;
            if (j >= r[i]) {
                p[i][j] -= p[i - 1][j - r[i]];
            }
            p[i][j] *= 1.0 / (m - 1);
            sum += p[i - 1][j];
        }
    }
    double ans = 0;
    for (int i = 0; i < s; i++) {
        ans += p[n][i];
    }
    ans *= (m - 1);
    cout << ans + 1.0 << endl;
}