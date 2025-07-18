#include <bits/stdc++.h>
#ifndef ONLINE_JUDGE
#include "debug.h"
#else
#define deb(x...)
#endif
#define time cerr << "Time Taken: " << (double)clock() / CLOCKS_PER_SEC * 1000 << " ms\n";
#define meme(s, x) memset(s, x, sizeof(s))
#define all(v) ((v).begin()), ((v).end())
#define line cout << "\n------------\n";
#define here cout << "Here\n";
#define pb push_back
#define el '\n'
#define int long long
using namespace std;
void judge();
// abcdefghijklmnopqrstuvwxyz
// primes 2,3,5,7,11,13,17,19,23,29,31,37,41,43,47,53,59,61,67,71,73,79,83,89,97,101
//--------------------------------------------------------------
void solve()
{
    int n;
    cin >> n;
    string s;
    cin >> s;
    int a = s[0] - '0';
    for (int i = 1; i < n; i++) {
        int b = s[i] - '0';
        // deb(a, b);
        if (b == 1) {
            if (a == 1) {
                cout << '-';
                a = 0;
            } else {
                cout << '+';
                a = 1;
            }
        } else {
            if (a == 1) {
                cout << '+';
                a = 1;
            } else {
                cout << '-';
                a = 0;
            }
        }
    }
    cout << el;
}

int32_t main()
{
    judge();
    int tt = 1;
    cin >> tt;
    while (tt--) {
        solve();
    }
    // time;
}

//

void judge()
{
#ifdef ONLINE_JUDGE
    ios_base::sync_with_stdio(false), cin.tie(NULL), cout.tie(0);
#endif
}