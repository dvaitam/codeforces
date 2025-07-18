#include <bits/stdc++.h>

#define large long long int
#define xlarge unsigned large
#define regpa pair<large,large>

using namespace std;

void go(int start, int limit, int n) {
    for (int i = start; i < limit; i += 2) {
        cout << i << " " << i + 1 << "\n";
        for (int j = 0; j < n - 1; j++)
            cout << i << " " << limit + j << "\n" << i + 1 << " " << limit + j << "\n";
    }
    for (int i = limit; i < limit + n - 1; i++)
        cout << i << " " << limit + n - 1 << "\n";
}

int main() {
    ios::sync_with_stdio(false);
    int n;
    cin >> n;
    if (n % 2 == 0) {
        cout << "NO";
        return 0;
    }
    cout << "YES\n" << 4 * n - 2 << " " << 2 * ((n - 1) + (n - 1) * (n - 1) + ((n - 1) / 2)) + 1 << "\n";
    go(1, n, n);
    go(2 * n, 3 * n - 1, n);
    cout << (2 * n) - 1 << " " << (4 * n) - 2;
    return 0;
}