#include<bits/stdc++.h>

using namespace std;

using namespace std::chrono;

#define int long long

#define fastio() ios_base::sync_with_stdio(false); cin.tie(NULL); cout.tie(NULL);

#define f(i,a,n) for(int i=a; i<n; i++)

#define rev_f(i, n, a) for (int i=n; i>a; i--)

#define ai(n, arr) f(i,0,n) cin>>arr[i]

#define matin(n, m, mat) f(i,0,n) f(j,0,m) cin>>mat[i][j]

#define all(x) x.begin(), x.end()

#define endl '\n'

#define c(x) cout << (x)

#define cty cout << "YES" << endl

#define ctn cout << "NO" << endl

#define csp(x) cout << (x) << " "

#define c1(x) cout << (x) << endl

#define c2(x, y) cout << (x) << " " << (y) << endl

#define c3(x, y, z) cout << (x) << " " << (y) << " " << (z) << endl

#define c4(a, b, c, d) cout << (a) << " " << (b) << " " << (c) << " " << (d) << endl

#define c5(a, b, c, d, e) cout << (a) << " " << (b) << " " << (c) << " " << (d) << " " << (e) << endl

#define c6(a, b, c, d, e, f) cout << (a) << " " << (b) << " " << (c) << " " << (d) << " " << (e) << " " << (f) << endl

#define cn(n, arr) f(i,0,n) cout<<arr[i]<<' '; cout<<endl



#ifndef ONLINE_JUDGE

#define debug(x) cout << #x << ' '; _print(x); cout << endl

#else

#define debug(x)

#endif



void _print(int t) {cout << t;}

void _print(bool t) {cout << t;}

void _print(double t) {cout << t;}

void _print(string t) {cout << t;}

template<class T, class V> void _print(pair<T, V> p);

template<class T> void _print(vector<T> v);

template <class T> void _print(set <T> v);

template <class T> void _print(multiset <T> v);

template <class T, class V> void _print(map <T, V> v);

template<class T, class V> void _print(pair<T, V> p) {cout<<'{'; _print(p.first); cout<<","; _print(p.second); cout <<'}';}

template<class T> void _print(vector<T> v) {cout << "[ "; for(T i: v) {_print(i); cout << " ";} cout << ']';}

template <class T> void _print(set <T> v) {cout << "[ "; for (T i : v) {_print(i); cout << " ";} cout << "]";}

template <class T> void _print(multiset <T> v) {cout << "[ "; for (T i : v) {_print(i); cout << " ";} cout << "]";}

template <class T, class V> void _print(map <T, V> v) {cout << "[ "; for (auto i : v) {_print(i); cout << " ";} cout << "]";}

template<typename typC,typename typD> ostream &operator<<(ostream &cout,const pair<typC,typD> &a) { return cout<<a.first<<' '<<a.second; }

template<typename typC,typename typD> ostream &operator<<(ostream &cout,const vector<pair<typC,typD>> &a) { for (auto &x:a) cout<<x<<'\n'; return cout; }

template<typename typC> ostream &operator<<(ostream &cout,const vector<typC> &a) { int n=a.size(); if (!n) return cout; cout<<a[0]; for (int i=1; i<n; i++) cout<<' '<<a[i]; return cout; }



int findx(int n) {

    int l=1, h=1e6;

    while (l<=h) {

        int x = (l+h)/2;

        if ((x+n)*(x+n-1) - x*(x-1) >= 2*n*(2*n+1)) h = x-1;

        else l = x+1;

    }

    if ((l+n)*(l+n-1) - l*(l-1) == 2*n*(2*n+1)) return l;

    return -1;

}



void solve() {

    int n;

    cin >> n;

    int x = findx(n);

    if (x==-1) ctn;

    else {

        cty;

        int i=1, j=x-1;

        for (; j<=2*n; i++, j++) c2(i, j);

        for (j=x+1-i; j<x-1; i++, j++) c2(i, j);

    }

}



signed main(){

    fastio();

    // cout << fixed << setprecision(15); // activate if the answers are in decimal

    auto start = high_resolution_clock::now();

    int t=0;

    cin>>t;

    while (t--) solve();

    auto stop = high_resolution_clock::now();

    auto duration = duration_cast<milliseconds>(stop-start);

    // cout << duration.count() << " milliseconds\n";



    // know case

    // int tvalue=10000, caseno=86;

    // f(i, 0, t) {

    //     if (t==tvalue) {

    //         int n;

    //         cin >> n;

    //         int arr[n];

    //         ai(n, arr);

    //         if (i==caseno-1) {

    //             c1(n);

    //             cn(n, arr);

    //         }

    //     }

    //     else solve();

    // }

}