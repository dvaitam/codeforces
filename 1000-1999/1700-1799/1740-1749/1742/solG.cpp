/**

 *      Author:  nicksms

 *      Created: 02.01.2023 22:52:06

**/



#include <bits/stdc++.h>

using namespace std;

#define ll long long



int main(){

    ios_base::sync_with_stdio(false);

    cin.tie(nullptr);

    int t; cin >> t;

    while (t--) {

        int n; cin >> n;

        vector<int> v(n);

        for (auto &&p : v) cin >> p;

        int cur = 0;

        for (int i = 0; i < min(32, (int)v.size()); i++) {

            int best = cur;

            for (int j = i; j < n; j++) {

                if ((cur | v[j]) > best) {

                    best = cur | v[j];

                    swap(v[i], v[j]);

                }

            }

            cur = best;

        }

        for (auto &&p : v) cout << p << " ";

        cout << "\n";

    }

    return 0;

}