#include <bits/stdc++.h>

using namespace std;

typedef long long ll;

ll gcd(ll a, ll b) {
    return b? gcd(b, a % b) : a;
}

int main() {
    ll n, m, k;
    cin >> n >> m >> k;
    ll tn = n, tm = m, tk = k;
    ll tgcd;
    tgcd = gcd(tn, tk);
    tn /= tgcd;
    tk /= tgcd;
    tgcd = gcd(tm, tk);
    tm /= tgcd;
    tk /= tgcd;
    if (tk != 1 && tk != 2) {
        cout << "NO";
        return 0;
    }
    if (tk == 2) {
        cout << "YES\n";
        cout << 0 << ' ' << 0 << endl;
        cout << tn << ' ' << 0 << endl;
        cout << 0 << ' ' << tm << endl;
        return 0;
    }
    if (tk == 1) {
        bool f = 0;
        if (!f && tn * 2 <= n) {
            tn *= 2;
            f = 1;
        }
        if (!f && tm * 2 <= m) {
            tm *= 2;
            f = 1;
        }
        if (!f) {
            cout << "NO";
            return 0;
        }
        cout << "YES\n";
        cout << 0 << ' ' << 0 << endl;
        cout << tn << ' ' << 0 << endl;
        cout << 0 << ' ' << tm << endl;
    }
    return 0;
}