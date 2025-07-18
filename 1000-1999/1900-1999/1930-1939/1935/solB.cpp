#include <bits/stdc++.h>

using namespace std;
const  int c=200005;
int n, t[c], pref[c], suf[c], db[c];
void solve() {
    cin >> n;
    for (int i=1; i<=n; i++) cin >> t[i];
    int mex=0;
    for (int i=1; i<=n; i++) {
        db[t[i]]++;
        while (db[mex]) mex++;
        pref[i]=mex;
    }
    mex=0;
    for (int i=0; i<=n; i++) db[i]=0;
    for (int i=n; i>=1; i--) {
        db[t[i]]++;
        while (db[mex]) mex++;
        suf[i]=mex;
    }
    for (int i=0; i<=n; i++) db[i]=0;

    int pos=0;
    for (int i=1; i<n; i++) {
        if (pref[i]==suf[i+1]) pos=i;
    }

    if (!pos) cout << -1 << "\n";
    else {
        cout << 2 << "\n";
        cout << 1 << " " << pos << "\n";
        cout << pos+1 << " " << n << "\n";
    }
}
int main()
{
    ios_base::sync_with_stdio(false);
    int w;
    cin >> w;
    while (w--) {
        solve();
    }
    return 0;
}