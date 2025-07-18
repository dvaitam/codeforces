#include<bits/stdc++.h>

using namespace std;

deque<string>k;

int main() {

    ios_base::sync_with_stdio(NULL);

    cin.tie(0);

    cout.tie(0);

    int n, m, ans = 0;

    string q = "";

    cin >> n >> m;

    map<string, int>mp;

    for(int i = 1; i <= n;i++) {

        string s, t;

        cin >> s;

        t = s;

        reverse(t.begin(), t.end());

        if(mp[t] != 0) {

            mp[t]--;

            k.push_front(t);

            k.push_back(s);

            ans += s.size() * 2;

        }

        else if(s == t) {

            q = s;

        }

        else {

            mp[s]++;

        }

    }

    cout << ans + q.size() << '\n';

    for(int i = 0; i < k.size() / 2;i++) {

        cout << k[i];

    }

    cout << q;

    for(int i = k.size() / 2; i < k.size();i++) {

        cout << k[i];

    }

}