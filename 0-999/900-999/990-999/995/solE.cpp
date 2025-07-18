#include <bits/stdc++.h>

using namespace std;

int inv(int a, int p) {
    int r = 1, b = p - 2;
    while (b) {
        if (b & 1) r = 1ll * r * a % p;
        a = 1ll * a * a % p;
        b >>= 1;
    }
    return r;
}

int random(int p) {
    long long x = 1ll * rand() * rand() + rand();
    return x % p;
}

int step(int a, int b) {
    int ret = 0;
    while (b > 0) {
        ret += a / b + 1;
        int tmp = a % b;
        a = b; b = tmp;
    }
    return ret;
}

vector <int> solve(int u, int p) {
    while (1) {
        int k = random(p);
        if (k == 0) continue;
        int a = 1ll * k * u % p, b = k;
        if (step(a, b) >= 100) continue;
        vector <int> ret;
        while (a > 0) {
            if (a >= b) {
                ret.push_back(0);
                a -= b;
            }
            else {
                ret.push_back(1);
                swap(a, b);
            }
        }
        return ret;
    }
}

int p, u, v;

int main(void) {
    srand(time(NULL));
    ios_base::sync_with_stdio(0); cin.tie(0); cout.tie(0);
    cin >> u >> v >> p;
    vector <int> a1 = solve(u, p);
    vector <int> a2 = solve(v, p);
    reverse(a2.begin(), a2.end());
    cout << a1.size() + a2.size() << endl;
    for (auto x: a1) cout << (x ? 3 : 2) << ' ';
    for (auto x: a2) cout << (x ? 3 : 1) << ' ';
    cout << endl;
    return 0;
}