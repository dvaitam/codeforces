#include <bits/stdc++.h>

using namespace std;



void solve() {

    int n;

    cin >> n;



    vector<int> h(n), w(n);

    for (int i = 0; i < n; i++) {

        cin >> h[i] >> w[i];

        if (h[i] > w[i]) swap(h[i], w[i]);

    }



    vector<int> p(n);

    iota(p.begin(), p.end(), 0);

    sort(p.begin(), p.end(), [&](int i, int j) {

        return h[i] < h[j];

    });



    int tmp = -1;

    vector<int> ans(n, -1);



    for (int i = 0; i < n; ) {

        int j = i;

        while (j < n && h[p[i]] == h[p[j]]) {

            j++;

        }



        for (int k = i; k < j; k++) {

            if (tmp != -1 && w[tmp] < w[p[k]]) {

                ans[p[k]] = tmp;

            }

        }



        for (int k = i; k < j; k++) {

            if (tmp == -1 || (w[tmp] > w[p[k]])) {

                tmp = p[k];

            }

        }



        i = j;

    }



    for (int num: ans) {

        cout << num + (num >= 0) << ' ';

    }

    cout << '\n';

}

 

int main() {

    ios::sync_with_stdio(0);

    cin.tie(0), cout.tie(0);



    int t = 1;

    cin >> t;

    while (t--) {

        solve();

    }

}