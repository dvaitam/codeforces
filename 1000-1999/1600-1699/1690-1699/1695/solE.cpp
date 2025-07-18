/*

 * File created on 07/07/2022 at 23:06:15.

 * Link to problem: 

 * Description: 

 * Time complexity: O()

 * Space complexity: O()

 * Status: ---

 * Copyright: Ⓒ 2022 Francois Vogel

*/



#include <bits/stdc++.h>



#include <ext/pb_ds/assoc_container.hpp>

#include <ext/pb_ds/tree_policy.hpp>



using namespace std;

using namespace __gnu_pbds;



#define pii pair<int, int>

#define f first

#define s second



#define vi vector<int>

#define all(x) x.begin(), x.end()

#define size(x) (int)((x).size())

#define pb push_back

#define ins insert

#define cls clear



#define ll long long

#define ld long double



typedef tree<int, null_type, less<int>, rb_tree_tag, tree_order_statistics_node_update> ordered_set;



const int siz = 6e5+40;



vector<pii> graph [siz];

bool vis [siz];

bool taken [siz];



void dfs(vi& c, int x) {

    c.pb(x);

    if (vis[x]) return;

    vis[x] = true;

    for (pii y : graph[x]) if (!taken[y.s]) {

        taken[y.s] = true;

        dfs(c, y.f);

        c.pb(x);

    }

}



signed main() {

    cin.tie(0);

    ios_base::sync_with_stdio(0);



    int n;

    cin >> n;

    n *= 2;

    for (int i = 0; i < n/2; i++) {

        int x, y;

        cin >> x >> y;

        x--;

        y--;

        graph[x].pb({y, i});

        if (x != y) graph[y].pb({x, i});

        taken[i] = false;

    }

    for (int i = 0; i < n; i++) vis[i] = false;

    vector<vi> cycles;

    for (int i = 0; i < n; i++) if (!vis[i]) {

        if (graph[i].empty()) continue;

        cycles.pb(vi());

        dfs(cycles.back(), i);

        cycles.back().pop_back();

    }

    vi r [2];

    string a [2] = {"", ""};

    string b [2] = {"", ""};

    for (const vi& i : cycles) {

        if (size(i) <= 2) {

            cout << -1;

            exit(0);

        }

        for (int j = 0; j < size(i)/2; j++) {

            r[0].pb(i[j]);

            r[1].pb(i[size(i)-1-j]);

            if (j%2 == 0) {

                a[0] += 'R';

                a[1] += 'R';

                b[0] += 'L';

                b[1] += 'L';

            }

            else {

                a[0] += 'L';

                a[1] += 'L';

                b[0] += 'R';

                b[1] += 'R';

            }

            if (j == 0) {

                if ((size(i)/2)%2 == 0) {

                    a[0].back() = 'U';

                    a[1].back() = 'D';

                }

                else {

                    a[0].back() = 'U';

                    a[1].back() = 'D';

                }

            }

            else if (j == size(i)/2-1) {

                if ((size(i)/2)%2 == 0) {

                    a[0].back() = 'U';

                    a[1].back() = 'D';

                }

                else {

                    b[0].back() = 'U';

                    b[1].back() = 'D';

                }

            }

        }

    }

    cout << 2 << ' ' << size(r[0]) << '\n';



    cout << '\n';



    for (int i : r[0]) cout << i+1 << ' ';

    cout << '\n';

    for (int i : r[1]) cout << i+1 << ' ';

    cout << '\n';



    cout << '\n';



    cout << a[0] << '\n';

    cout << a[1] << '\n';



    cout << '\n';



    cout << b[0] << '\n';

    cout << b[1] << '\n';



    cout.flush();

    int d = 0;

    d++;

}