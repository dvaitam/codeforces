#include <array>

#include <bits/stdc++.h>

using namespace std;

using LL = long long;

using ULL = unsigned long long;

#define INF 0x3f3f3f3f

#define endl '\n'

#define IO ios::sync_with_stdio(0), cin.tie(0), cout.tie(0);

//#define LOCAL



void solve() {

    int n;

    cin >> n;

    cout << 2 << " ";

    for(int i = 3; i < n; i ++ ) cout << i << " ";

    cout << n << " " << 1 << endl;

}



int main(){

    IO

    clock_t ccc1 = clock();

#ifdef LOCAL

    ifstream cin("1.in");

    ofstream cout("1.out");

#endif

    int T;

    cin >> T;

    while(T--) {

        solve();

    }

    

    

end:

    cerr << "Time : " << clock() - ccc1 << "ms" << endl;

    return 0;

}