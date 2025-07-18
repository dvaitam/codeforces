#include <iostream>
#include <cstdio>
#include <string>

using namespace std;

int min(int a, int b) {
    if (a < b)
        return a;
    return b;
}

int main() {
    int n, k = -1;
    cin >> n;
    if (n <= 2) {
        cout << "No";
        return 0;
    }
    if (n % 2 == 0) {
        cout << "Yes" << endl;
        cout << "2 1 " << n << endl;
        cout << n - 2;
        for (int i = 2; i < n; i++)
            cout << " " << i;
        cout << endl;
        return 0;
    }
    k = (n + 1) / 2;
    cout << "Yes" << endl;
    cout << "1 " << k << endl;
    cout << n - 1;
    for (int i = 1; i <= n; i++)
        if (i != k)
            cout << " " << i;
    cout << endl;
	return 0;
}