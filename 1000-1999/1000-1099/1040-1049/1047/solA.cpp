#include<bits/stdc++.h>

using namespace std;

int main() {
    ios_base::sync_with_stdio(0);
    cin.tie(0);
    cout.tie(0);
    int n, i, j, k;
    cin >> n;
    int d = n / 3, r = n % 3;
    int a = d, b = d, c = d + r;
    if (a % 3 == 0 && c % 3 != 0) {
        a += 1;
        b -= 1;
    } else if (a % 3 == 0 && c % 3 == 0) {
        a -= 1;
        b -= 1;
        c += 2;
    } else if (a % 3 != 0 && c % 3 == 0) {
        if ((a + 1) % 3 == 0) {
            a += 2;
            c -= 2;
        } else {
            a += 1;
            c -= 1;
        }
    }
    cout << a << " " << b << " " << c;

    return 0;
}