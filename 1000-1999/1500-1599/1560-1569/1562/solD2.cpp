#include<bits/stdc++.h>

#include <ext/pb_ds/assoc_container.hpp>

#include <ext/pb_ds/tree_policy.hpp>

#define int long long

#define sz(x) static_cast<int>((x).size())

using namespace std;

using namespace __gnu_pbds;



typedef tree<int, null_type, less<int>, rb_tree_tag, tree_order_statistics_node_update> ordered_set;

typedef tree<pair<int,int>, null_type, less<pair<int,int>>, rb_tree_tag, tree_order_statistics_node_update> ordered_multiset;







/* ----------------------------------------------------- GO DOWN ---------------------------------------------------------------------- */





const int INF = 1e18;







signed main(){



        ios::sync_with_stdio(0);

        cin.tie(0);

        cout.tie(0);

        

        int t;

        cin >> t;

        while (t--) {



                int n, q;

                cin >> n >> q;

                string s;

                cin >> s;

                int p[n + 1];

                p[0] = 0;

                vector<int> v1[n + 1];

                vector<int> v2[n + 1];

                for (int i = 0; i < n; i++) {

                        if (i % 2 == 0) {

                                if (s[i] == '+') p[i + 1] = p[i] + 1;

                                else p[i + 1] = p[i] - 1;

                        }

                        else {

                                if (s[i] == '+') p[i + 1] = p[i] - 1;

                                else p[i + 1] = p[i] + 1;  

                        }

                        if (p[i + 1] >= 0) {

                                v1[p[i + 1]].push_back(i + 1);

                        }

                        else {

                                v2[-p[i + 1]].push_back(i + 1);

                        }

                }





                for (int i = 0; i < q; i++) {



                        int l, r;

                        cin >> l >> r;

                        int c = p[r] - p[l - 1];

                        if (c == 0) {

                                cout << 0 << "\n";

                                continue;

                        }



                        if (c % 2 == 0) {

                                c--;

                                cout << 2 << "\n";

                                cout << r << " ";

                                r--;

                        }

                        else {

                                cout << 1 << "\n";

                        }

                        



                                

                        int f;

                        if (p[r] > p[l - 1]) {

                                f = p[l - 1] + (p[r] - p[l - 1] + 1) / 2;

                        }

                        else {

                                f = p[l - 1] - (p[l - 1] - p[r] + 1) / 2;

                        }

                        int pos = 0;

                        if (f >= 0) {

                                auto it = upper_bound(v1[f].begin(), v1[f].end(), l - 1) - v1[f].begin();

                                pos = v1[f][it];

                        }

                        else {

                                auto it = upper_bound(v2[-f].begin(), v2[-f].end(), l - 1) - v2[-f].begin();

                                pos = v2[-f][it];

                        }

                        cout << pos << "\n";







                }







        }





}