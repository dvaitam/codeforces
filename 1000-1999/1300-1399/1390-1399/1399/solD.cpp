#include<bits/stdc++.h>

#define LF "\n"

using namespace std;

using ll = long long;

constexpr int N = (int)1e5 + 10;



void solve() {

    int n;

    cin >> n;

    string s;

    cin >> s;

    queue<int> pos[2];

    vector<int> cpos(n);

    int cnt = 0;

    for (int i = 0; i < n; i++) {

        int now = s[i] - '0';

        if (!pos[1-now].empty()) {

            pos[now].push(pos[1 - now].front());

            cpos[i] = pos[1 - now].front();

            pos[1 - now].pop();

        }

        else {

            pos[now].push(++cnt);

            cpos[i] = cnt;

        }

    }

    cout << cnt << LF;

    for (int i = 0; i < n; i++) {

        cout << cpos[i] << " ";

    }

    cout << LF;

}



int main() {

    ios::sync_with_stdio(false);

    cin.tie(nullptr);

    cout.tie(nullptr);

    

    int t = 1;

    cin >> t;

    

    while (t--) {

        solve();

    }

    

    return 0;

}