#include <bits/stdc++.h>



using namespace std;

#define sz(s) (int)(s.size())

#define all(v) v.begin(),v.end()

#define clr(d, v) memset(d,v,sizeof(d))

#define ll long long





void file() {

    std::ios_base::sync_with_stdio(0);

    cin.tie(NULL);

    cout.tie(NULL);

}



int main() {

    file();

    int n;

    cin >> n;

    vector<int> v(n - 1);

    ll sum = 0;

    int pos = 0;

    bool ok = true;

    for (auto &it: v) {

        cin >> it;

        sum += it;

        if (sum > 0)

            pos++;



    }



    vector<int> ans = {n - pos};

    vector<bool> vis(n + 1);

    vis[n - pos] = true;

    for (int i = 1; i < n; i++) {

        ans.push_back(v[i - 1] + ans.back());

        if (ans.back() > n || ans.back() <= 0 || vis[ans.back()]) {

            cout << "-1";

            return 0;

        }

        vis[ans.back()] = true;

    }

    for (auto it: ans)cout << it << " ";





}