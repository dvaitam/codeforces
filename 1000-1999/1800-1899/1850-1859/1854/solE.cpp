#include <bits/stdc++.h>

using namespace std;

int main() {
    int64_t m;
    cin >> m;
    int i = 0;
    int64_t cnt[60]{1};
    int64_t cans = 0;
    vector<int> ans;
    for (int i = 1, r = 0, c = 0; i <= 60 && c <= 60;) {
        if (r < (1ull << i) && m >= cnt[60 - i]) {
            m -= cnt[60 - i];
            for (int j = 60; j-- > i;) {
                cnt[j] += cnt[j - i];
            }
            ++r;
            ++c;
            ans.push_back(i);
        } else {
            r = 0;
            ++i;
        }
    }
    cout << ans.size() << '\n';
    for (auto a: ans) cout << a << ' ';
}