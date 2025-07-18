#include <bits/stdc++.h>
using namespace std;

const int MAX_N = 800;

int N, A[MAX_N];
bool g[MAX_N][MAX_N];
bool lmk[MAX_N][MAX_N], rmk[MAX_N][MAX_N];
bool lt[MAX_N][MAX_N], rt[MAX_N][MAX_N];

bool rrecur(int, int);

int gcd(int a, int b) {
    return b == 0 ? a : gcd(b, a % b);
}

bool lrecur(int i, int j) {
    if (i > j) return 1;
    if (i == j) return g[i - 1][i];
    if (!lmk[i][j]) {
        for (int k = i; k <= j; k++) if (g[i - 1][k]) {
            lt[i][j] |= rrecur(i, k - 1) && lrecur(k + 1, j);
            if (lt[i][j]) break;
        }
        lmk[i][j] = 1;
    }
    return lt[i][j];
}

bool rrecur(int i, int j) {
    if (i > j) return 1;
    if (i == j) return g[j][j + 1];
    if (!rmk[i][j]) {
        for (int k = i; k <= j; k++) if (g[k][j + 1]) {
            rt[i][j] |= rrecur(i, k - 1) && lrecur(k + 1, j);
            if (rt[i][j]) break;
        }
        rmk[i][j] = 1;
    }
    return rt[i][j];
}

bool solve() {
    for (int i = 0; i < N; i++) for (int j = i; j < N; j++)
        g[i][j] = gcd(A[i], A[j]) > 1;

    if (lrecur(1, N - 1)) return 1;
    if (rrecur(0, N - 2)) return 1;
    for (int k = 1; k < N - 1; k++) {
        if (rrecur(0, k - 1) && lrecur(k + 1, N - 1))
            return 1;
    }

    return 0;
}

int main() {
    ios_base::sync_with_stdio(false); cin.tie(NULL);

    cin >> N;
    for (int i = 0; i < N; i++)
        cin >> A[i];

    cout << (solve() ? "Yes" : "No"); 

    return 0;
}