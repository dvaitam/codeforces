#include <iostream>
#include <vector>
#include <cassert>
#include <map>
#include <set>
#include <queue>
#include <array>
#include <algorithm>
#include <functional>
#include <numeric>
#include <limits>

using namespace std;

typedef long long ll;
template<typename T> using minPQ = priority_queue<T, vector<T>, greater<T>>;

int main()
{

    ios_base::sync_with_stdio(0);
    cin.tie(0);

    int tttt=1;
    cin>>tttt;
    while(tttt--)
    {
        int n;
        cin >> n;
        n <<= 1;
        string s;
        cin >> s;
        int cnt = 0;
        for(int i = 0; i < n; i++) {
            cnt += (s[i] == '1'? 1: 0);
        }
        if(s.front() != s.back() || cnt%2) {
            cout << -1 << '\n';
            continue;
        }
        vector<int> flips[2];
        for(int i = 1; i < n - 2; i++) {
            if(s[i]=='1') {
                flips[i%2].push_back(i);
                s[i] = '0';
                s[i+1] = s[i+1] == '0'? '1': '0';
            }
        }
        string base_seq = "(";
        for(int i = 0; i < n/2-1; i++) {
            base_seq.push_back('(');
            base_seq.push_back(')');
        }
        base_seq.push_back(')');
        vector<string> ans;
        if(s.front() == '1') {
            ans.push_back(base_seq);
        }
        for(int i = 0; i < 2; i++) {
            string tmp = base_seq;
            for(int j = 0; j < flips[i].size(); j++) {
                swap(tmp[flips[i][j]], tmp[flips[i][j]+1]);
            }
            ans.push_back(tmp);
        }
        cout << ans.size() << '\n';
        for(int i = 0; i < ans.size(); i++) {
            cout << ans[i] << '\n';
        }
    }
    return 0;
}