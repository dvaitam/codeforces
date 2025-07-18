// http://codeforces.com/contest/1099/problem/B

#include <bits/stdc++.h>
using namespace std;
int main() {
    cin.sync_with_stdio(0);
    cin.tie(0);

    int N;
    cin >> N;
    double aaa = sqrt(N);
    cout << ceil(floor(aaa) + N/floor(aaa)) << endl;
}