#include <iostream>

#include<vector>

#include<string>

#include<set>

#include<algorithm>

#include<cmath>

#include<map>

#include<stack>

#include<iomanip>

#include <cstring>

#include<string>

#include<queue>

#include<deque>

#include <numeric>



using namespace std;

#define ll long long

#define forn(i,a, n) for(long long i = a; i<(n); i++)

#define pb push_back

#define all(v) v.begin(),v.end()

#define rall(v) v.rbegin(),v.rend()

#define vec(n)  vector <ll>a(n);

#define endl "\n"

#define sz size()

const int N = 2e5 + 20, M = 1000 + 7;

const long long mod = 1000000007;

int cr[] = { 0, 0, -1, 1 };

int cc[] = { -1, 1, 0, 0 };



void fast()

{

    ios_base::sync_with_stdio(false);

    cin.tie(NULL);

    cout.tie(NULL);

}

#define isON(n,k) ((((n)>>k))&1)

#define ll long long

#define no cout << "NO\n"

#define yes cout << "YES\n"

//log10(n) + 1 = size n

ll gcd(ll a, ll b) {

    return b == 0 ? a : gcd(b, a % b);

}

ll lcm(ll a, ll b) {

    return a / gcd(a, b) * b;

}

int n,k;

bool BS(vector <int>v,ll mid) {

    int m = n / 2 ;

    ll k2 = k;

    if (v[m] > mid)return true;

    forn(i, m, n) {

        ll x = mid - v[i];

        if (x <= k2)k2 -= x;

        else return false;

    }

    return true;



}





void solve() {





    string s; cin >> s;

    int zs = 0, ones = 0;

    forn(i, 0, s.sz) {

        if (s[i] == '0')zs++;

        else ones++;

    }

    if (!zs || !ones)cout << s;

    else {



        while (zs > 0) {

            cout << "01";

            --zs;

        }



        while (ones > 0) {

            cout << "01";

            --ones;

        }

    }

    



}









int main() {

    fast();

    int t = 1;

    cin >> t;

    while (t--) {

        solve();

        cout << endl;

    }

    return 0;

}

//a z: 97 122