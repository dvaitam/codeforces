#include <bits/stdc++.h>

#include <ext/pb_ds/assoc_container.hpp>

#include <random>



using namespace std;

using namespace __gnu_pbds;



typedef long long ll;

typedef unsigned long long ull;

typedef tree<int, null_type, less<>, rb_tree_tag, tree_order_statistics_node_update> ordered_set;

#define AboTaha_on_da_code ios_base::sync_with_stdio(false); cin.tie(NULL); cout.tie(NULL);

#define X first

#define Y second



const int dx[8]={0, 0, 1, -1, 1, -1, -1, 1}, dy[8]={1, -1, 0, 0, 1, -1, 1, -1};

const int M = 1e9+7, M2 = 998244353;

const double EPS = 1e-8;



void burn(int tc)

{

    int n; cin >> n;

    vector <int> b(n);

    for (auto &i : b) cin >> i;

    for (int i = 0; i < n; i++) {

        for (int j = 0; j < n; j++) {

            cout << ((i-j)*j%n+b[i]+n)%n << ' ';

        }

        cout << '\n';

    }

}



int main()

{

    // I live for this shit

    AboTaha_on_da_code



//    freopen("zeros.in", "r", stdin);

//    freopen("Aout.txt", "w", stdout);



    int T = 1; // cin >> T;



    for (int i = 1; i <= T; i++) {

//        cout << "Case " << i << ": ";

        burn(i);

//        cout << '\n';

    }

    return 0;

}