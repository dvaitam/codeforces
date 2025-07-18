#include <bits/stdc++.h>

using namespace std;



int main() {

    int n; cin >> n;

    vector<int> ops; ops.reserve(3*n*n/2);



    auto do_thing = [&ops,n](int offset, int m) {

        for (int i = 0; i < m-1; ++i) {

            ops.push_back((i+offset)%n);

        }

        for (int k = 1; m-2*k >= 1; ++k) {

            for (int i = m-k-2; i >= k-1; --i) {

                ops.push_back((i+offset)%n);

            }

            for (int i = k; i <= m-k-1; ++i) {

                ops.push_back((i+offset)%n);

            }

        }

    };



    do_thing(0, n);

    do_thing(n-n/2, n-n%2);

    do_thing(0, n);



    cout << ops.size() << '\n';

    for (int const i : ops) {

        cout << i << ' ';

    }

    cout << '\n';

}