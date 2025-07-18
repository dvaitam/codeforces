#include <iostream>
#include <bits/stdc++.h>
using namespace std;
int t, x;

void Solve() {
    for (int i=1; i<=t; i++) {
        cin >> x;
        cout << 1 << ' ' << x - 1 << '\n';
    }
}

int main()
{
    ios_base::sync_with_stdio(0);
    cin.tie(0);
    cout.tie(0);
    cin >> t;
    Solve();
    return 0;
}