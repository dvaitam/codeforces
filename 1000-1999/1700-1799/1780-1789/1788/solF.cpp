#include <bits/stdc++.h>



using namespace std;

using i64 = long long;



int main() {

    ios::sync_with_stdio(false);

    cin.tie(nullptr);



    int n, q;

    cin >> n >> q;



    vector<int> d(n);

    vector<array<int, 2>> edges(n - 1);



    for (int i = 0; i < n - 1; i++) {

        int u, v;

        cin >> u >> v;

        d[--u] ^= 1;

        d[--v] ^= 1;

        edges[i] = {u, v};

    }



    vector<vector<array<int, 2>>> adj(n);

    

    for (int i = 0; i < q; i++) {

        int u, v, w;

        cin >> u >> v >> w;

        u--;

        v--;

        adj[u].push_back({v, w});

        adj[v].push_back({u, w});

    }



    vector<int> a(n, -1), b(n);



    for (int i = 0; i < n; i++) {

        if (a[i] != -1) continue;



        queue<int> q;

        q.push(i);

        a[i] = 0;



        while (!q.empty()) {

            int u = q.front();

            q.pop();



            b[u] = i;



            for (auto [v, w] : adj[u]) {

                if (a[v] == -1) {

                    a[v] = a[u] ^ w;

                    q.push(v);

                }

                if (a[v] != (a[u] ^ w)) {

                    cout << "No\n";

                    return 0;

                }

            }

        }

    }

    

    int ans = 0;

    vector<int> cnt(n);

    

    for (int i = 0; i < n; i++) {

        if (d[i]) {

            ans ^= a[i];

            cnt[b[i]]++;

        }

    }



    for (int i = 0; i < n; i++) {

        if (cnt[i] % 2 == 1) {

            for (int j = 0; j < n; j++) {

                if (b[j] == i) {

                    a[j] ^= ans;

                }

            }

            break;

        }

    }

    

    cout << "Yes\n";

    for (int i = 0; i < n - 1; i++) {

        auto [u, v] = edges[i];

        cout << (a[u] ^ a[v]) << " \n"[i == n - 1];

    }



    return 6/22;

}