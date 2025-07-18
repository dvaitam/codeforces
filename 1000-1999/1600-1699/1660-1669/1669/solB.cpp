#include <bits/stdc++.h>
#include <iostream>
#include <vector>
#include <algorithm>
#include <array>
#include <set>
#include <map>
#include <numeric>
#include <string>
#include <functional>

using namespace std;

typedef long long ll;

void solve25() {
    ll n;
    cin >> n;

    ll ans;
    ans = (n / 3) + 1;
    if (n % 3 == 0) {
        ans = n / 3;
    } else if (n == 1) {
        ans = 2;
    } else if (n % 3 == 1 || n % 2 == 0) {
        if (n / 2 <= (n / 3) + 1) {
            ans = n / 2;
        }
    }
    cout << ans << '\n';
}

void solve26() {
    int n;
    cin >> n;

    vector<int> v(n);
    for (int i = 0; i < n; i++) {
        v[i] = i + 1;
    }
    cout << n << '\n';
    for (int j = 0; j < n; j++) {
        for (auto &x: v) {
            cout << x << ' ';
        }
        cout << '\n';

        swap(v[j], v[j + 1]);

    }
}

void solve27() {
    ll n;
    cin >> n;

    vector<ll> a(n);
    for (auto &x: a) {
        cin >> x;
    }
    for (int i = 0; i < n; i++) {
        if (a[i] % 10 == 0) {
            continue;
        }
        if (a[i] % 5 == 0) {
            a[i] += (a[i] % 10);
            continue;
        } else {
            while (a[i] % 10 != 2) {
                a[i] += (a[i] % 10);
            }
        }
        a[i] %= 20;
    }
    for (int i = 1; i < n; i++) {
        if (a[0] != a[i]) {
            cout << "No\n";
            return;
        }
    }
    cout << "Yes\n";
}


void solve28() {
    int n;
    cin >> n;

    vector<int> v(n);
    for (int i = 0; i < n; i++) {
        cin >> v[i];

    }
    for (int i = 1; i <= n; i++) {
        int j = i;
        map<int, int> m;
        m[j]++;
        while (m[j] != 2) {
            m[v[j - 1]]++;
            j = v[j - 1];
        }
        cout << j << ' ';
    }
}


void solve29() {
    int n;
    cin >> n;

    int res = 0;
    int mx = INT_MIN;
    map<int, int> m;
    vector<int> a(n);
    for (auto &x: a) {
        cin >> x;

        m[x]++;
        mx = max(m[x], mx);
        if (m[x] > 1) {
            res++;
        }
    }
    cout << mx << ' ' << n - res;
}

void solve30() {
    int n;
    cin >> n;

    string s;
    cin >> s;

    map<char, int> mp;
    int res = 0;
    for (int i = 0; i < n; i++) {
        if (mp[s[i]] == 0) {
            res += ((count(s.begin(), s.end(), s[i])) - 1);
        }
        mp[s[i]]++;
    }
    if (n <= 26) {
        cout << res << '\n';
    } else {
        cout << "-1\n";
    }
}

void solve31() {
    int n;
    cin >> n;

    string Herloc;
    string Mar;
    cin >> Herloc >> Mar;

    int mnCount = n;
    int i = 0;
    int j = 0;
    sort(Herloc.begin(), Herloc.end());
    sort(Mar.begin(), Mar.end());
    while (i < n && j < n) {
        if (Herloc[i] <= Mar[j]) {
            i++;
            j++;
            mnCount--;
        } else {
            j++;
        }
    }
    j = 0;
    i = 0;
    int mxCount = 0;
    while (i < n && j < n) {
        if (Herloc[i] < Mar[j]) {
            i++;
            j++;
            mxCount++;
        } else {
            j++;
        }


    }
    cout << mnCount << '\n' << mxCount << '\n';

}


void solve32() {
    ll n;
    cin >> n;

    vector<ll> a(n);
    for (auto &x: a) {
        cin >> x;

    }
    map<ll, ll> mp;
    ll res = 0;
    for (int i = 0; i < n; i++) {
        mp[a[i]]++;
    }
    for (int i = 0; i < n; i++) {
        if (mp[a[i]] != -1 && mp[-a[i]] != -1) {
            if (a[i] != 0) {
                res += (mp[-a[i]] * mp[a[i]]);
            } else if (mp[a[i]] >= 2) {
                res += (mp[a[i]] * (mp[a[i]] - 1)) / 2;
            }
            mp[a[i]] = -1;
            mp[-a[i]] = -1;
        }
    }
    cout << res << '\n';
}


int capsLook(char &a) {
    return a >= 'A' && a <= 'Z' ? 1 : 0;
}

void solve33() {
    string s, ans;
    cin >> s;

    ans = s;
    int flag = 1;
    for (int i = 1; i < s.size(); i++) {
        if (capsLook(s[i])) {
            s[i] = (int) s[i] + 32;
        } else {
            flag = 0;
            break;
        }
    }
    if (flag == 1) {
        s[0] = !capsLook(s[0]) ? (int) s[0] - 32 : (int) s[0] + 32;
    }

    cout << (flag == 1 ? s : ans) << '\n';
}


void solve1() {
    ll n;
    cin >> n;

    int flag = 0;
    map<int, int> mp;
    for (int i = 0; i < n; i++) {
        int x;
        cin >> x;

        mp[x]++;
        if (mp[x] >= 3 && flag != 1) {
            cout << x << ' ' << '\n';
            flag = 1;
        }
    }
    if (flag == 0){
        cout << "-1\n";
    }

}

int main() {
    ios_base::sync_with_stdio(false);
    cin.tie(nullptr);
    cout.tie(nullptr);


    int t;
    cin >> t;
    while (t--) {
        //char a = 'a';
        //char b = 'A';
        //cout << (int)a  <<'\n';
        //cout << (int)b;
        solve1();
    }


}