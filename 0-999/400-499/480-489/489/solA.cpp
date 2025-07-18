#include <iostream>

#include <string>

#include <vector>

#include <algorithm>

#define ll long long int

#define fr(n) for (int i = 0; i < n; i++)

#define fr2(n) for (int j = 0; j < n; j++)

#define min3(a, b, c) min(a, min(b, c))

#define max3(a, b, c) max(a, max(b, c))

using namespace std;



int main() {

    ios::sync_with_stdio(false);

    int n, a[3003];

    cin >> n;

    fr(n) cin >> a[i];

    cout << n << '\n';

    fr(n) {

        int x = 2147483647, y = 0;

        for (int j = i; j < n; j++) if (a[j] < x) x = a[j], y = j;

        swap(a[i], a[y]);

        cout << i << ' ' << y << '\n';

    }

}