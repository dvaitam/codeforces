#include <algorithm>
#include <cmath>
#include <fstream>
#include <iomanip>
#include <iostream>
#include <map>
#include <queue>
#include <random>
#include <set>
#include <stack>
#include <string>
#include <unordered_map>
#include <unordered_set>
#include <utility>
#include <vector>
#include <cctype>
#include <cstring>
#include <list>
#include <assert.h>

using namespace std;
using ll = long long;
using ld = long double;

ll inf = 99999999999;

int main() {
    ios::sync_with_stdio(0);
    cin.tie(0); cout.tie(0);
    ll t;
    cin >> t;
    while (t--) {
        ll n, k;
        cin >> n >> k;
        ll d = k - 3;
        n -= d; k = 3;
        if (n % 3 == 0) {
            cout << n / 3 << ' ' << n / 3 << ' ' << n / 3 << ' ';
        } else if (n % 2 == 0) {
            if (n / 2 % 2 == 1) {
                cout << (n - 2) / 2 << ' ' << (n - 2) / 2 << ' ' << 2 << ' ';
            } else {
                cout << n / 2 << ' ' << n / 4 << ' ' << n / 4 << ' ';
            }
        } else {
            cout << (n - 1) / 2 << ' ' << (n - 1) / 2 << ' ' << 1 << ' ';
        }
        for (ll i = 0; i < d; ++i) cout << 1 << ' ';
        cout << '\n';
    }
    return 0;
}