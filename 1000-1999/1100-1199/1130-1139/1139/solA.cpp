//هوالحق
#include <bits/stdc++.h>
using namespace std;
typedef long long ll;

int main() {
    ios::sync_with_stdio(false);
    int n;
    ll r = 0;
    cin >> n;
    for (int i = 0; i < n; ++i) {
        char e;
        cin >> e;
        if (!((int(e) - int(0)) % 2)) r += i + 1;
    }
    cout << r;
}