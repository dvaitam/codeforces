#include <bits/stdc++.h>

using namespace std;

int main() {

    int n;

    cin >> n;

    if (n == 1 || n == 2) {

        cout << "1\n1 1";

        return 0;

    }

    if (n % 3 == 1) {

        n++;

    }

    int orig_n = n;

    if (n % 3 == 0) {

        n --;

    }

    if (n % 3 == 2) {

        vector<pair<int,int> > vec;

        for (int i = 0; i < n; i++) {

            int j = n/3 - i;

            if (j >= 0 && j < n) {

                //cout << i + 1 << " " << j + 1 << '\n';

                vec.push_back(make_pair(i + 1, j + 1));

            }

            j = n - 1 - i;

            if (abs(i - j) <= n/3) {

                vec.push_back(make_pair(i + 1, j + 1));

            }

        }

        if (orig_n % 3 == 0) {

            vec.push_back(make_pair(orig_n, orig_n));

        }

        cout << vec.size() << '\n';

        for (auto& p: vec) {

            cout << p.first << " " << p.second << '\n';

        }

    }



}