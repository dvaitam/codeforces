#include <algorithm>
#include <cmath>
#include <cstring>
#include <deque>
#include <iomanip>
#include <iostream>
#include <map>
#include <numeric>
#include <queue>
#include <set>
#include <string>
#include <unordered_map>
#include <vector>

using namespace std;
using ll = long long;

bool check(string &s, int curIndex, int lastIndex) {
    int sum = 0;
    for (int i = curIndex; i > lastIndex; i--) {
        sum += (s[i] - '0');
        if (sum % 3 == 0) {
            return true;
        }
    }
    return false;
}

int main() {
    cin.tie(nullptr);
    ios::sync_with_stdio(false);

    string s;
    cin >> s;

    int len = s.size();

    int lastIndex = -1;
    int ans = 0;
    for (int i = 0; i < len; i++) {
        if (check(s, i, lastIndex)) {
            ans++;
            lastIndex = i;
        }
    }

    cout << ans << endl;

    return 0;
}