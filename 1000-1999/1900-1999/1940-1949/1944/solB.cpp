#include <iostream>
#include <vector>
#include <string>
#include <cmath>
#include <algorithm>
#include <map>
#include <set>
typedef long long ll;
using namespace std;

int main() {
    ios::sync_with_stdio(0);
    cin.tie(0);
    int t;
    cin >> t;
    while (t--) {
        int n, k;
        cin >> n >> k;
        k *= 2;

        vector<ll> a(2*n);
        vector<ll> b(n + 1, 0); 

        for (int i = 0; i < 2 * n; i++) {
            cin >> a[i];
        }
        for (int i = 0; i < n; i++) {
            b[a[i]]++;
        }

        vector<int> d0, d1, d2;

        for (int i = 1; i <= n; i++) {
            if (b[i] == 0) d0.push_back(i);
            else if(b[i] == 1) d1.push_back(i);
            else d2.push_back(i);
        }

        int temp = 0;
        for (auto x : d2) {
            if (temp < k) {
                temp+=2;
                cout << x << ' ' << x << ' ';
            }
        }

        for (auto x : d1) {
            if (temp < k) {
                temp++;
                cout << x << ' ';
            }
        }
        cout << "\n";
        temp = 0;
        for (auto x : d0) {
            if (temp < k) {
                temp += 2;
                cout << x << ' ' << x << ' ';
            }
        }

        for (auto x : d1) {
            if (temp < k) {
                temp++;
                cout << x << ' ';
            }
        }
        cout << "\n";

    }
    return 0;
}