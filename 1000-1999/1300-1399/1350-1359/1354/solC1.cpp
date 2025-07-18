/*

* @Author: tsumondai

* @Target: VOI2023

*/

#include <bits/stdc++.h>

using namespace std;



#define int long long

#define fi first

#define se second

#define pb push_back

#define mp make_pair

#define foru(i, l, r) for(int i = l; i <= r; i++)

#define ford(i, r, l) for(int i = r; i >= l; i--)



typedef pair<int, int> ii;

typedef pair<ii, int> iii;

typedef pair<ii, ii> iiii;



const int N = 1e6 + 5;



const int oo = 1e9 + 7, mod = 1e9 + 7;



double n, m;

vector<int> arr;



void read() {



}



void init() {



}



void process() {

    double pi=atan(1)*4;

    cin >> n;

    cout << fixed << setprecision(6) <<  1/tan(pi/2/n) << '\n';

}



signed main() {

    cin.tie(0)->sync_with_stdio(false);

    //freopen(".inp", "r", stdin);

    //freopen(".out", "w", stdout);

    read();

    init();

    int t;

    cin >> t;

    while (t--) process();

    return 0;

}