#include <bits/stdc++.h>

using namespace std;



template<typename T> struct manacher {

    vector<int> lengths;

    manacher(T s) { combined_manacher(s); }

    template<typename T_string> static vector<int> odd_manacher(const T_string &pattern) {

        int n = int(pattern.size());

        vector<int> radius(n, 0);

        int loc = 0;

        for (int i = 1; i < n; i++) {

            if (i <= loc + radius[loc]) radius[i] = min(radius[loc - (i - loc)], loc + radius[loc] - i);

            while (i - radius[i] > 0 && i + radius[i] < n - 1 && 

                pattern[i - radius[i] - 1] == pattern[i + radius[i] + 1]) radius[i]++;

            if (i + radius[i] > loc + radius[loc]) loc = i;

        }

        return radius;

    }

    template<typename T_string> void combined_manacher(const T_string &pattern) {

        int n = int(pattern.size());

        T_string extended(2 * n + 1, 0);

        for (int i = 0; i < n; i++) extended[2 * i + 1] = pattern[i];

        lengths = odd_manacher(extended);

    }

    // whether s[start, end) is a palindrome

    inline bool is_palindrome(int start, int end) const {

        return lengths[start + end] >= end - start;

    }

};



int main() {

    ios::sync_with_stdio(false);

    cin.tie(nullptr);



    int test_cases;

    cin >> test_cases;

    while (test_cases--) {

        string s;

        cin >> s;

        int n = int(s.size());

        int l = 0, r = n - 1;

        while (l <= r && s[l] == s[r]) l++, r--;

        if (r < l) {

            cout << s << '\n';

            continue;

        }

        pair<int, string> ans;

        for (int iter = 0; iter < 2; ++iter) {

            manacher M(s);

            int i = -1;

            for (int j = l; j <= r; ++j) if (M.is_palindrome(l, j + 1)) i = j;

            string t = s.substr(0, i + 1) + s.substr(r + 1);

            ans = max(ans, make_pair(int(t.size()), t));

            l = n - 1 - r, r = n - 1 - l;

            reverse(s.begin(), s.end());

        }

        cout << ans.second << '\n';

    }

}