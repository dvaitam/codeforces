#include<bits/stdc++.h>
using namespace std;
#define ll long long
#define fastread()      (ios_base:: sync_with_stdio(false),cin.tie(NULL));
// LL, other ways to simplify it?, tricks & math
int main() {
    fastread();
    int t; cin >> t;
    while (t--) {
        int n, k; cin >> n >> k;
        if (n<=2*k) {
            cout << -1 << endl;
            continue;
        }
        int i = 1, j = n;
        while (k--) {
            cout << i++ << " ";
            cout << j-- << " ";
        }
        for (int l = i; l < j+1; ++l) {
            cout << l << " ";
        }
        cout << endl;
    }
    return 0;
}