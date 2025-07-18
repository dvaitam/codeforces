#include <cmath>
#include <cstdio>
#include <cstdlib>
#include <cstring>

#include <algorithm>
#include <bitset>
#include <functional>
#include <iostream>
#include <iterator>
#include <limits>
#include <map>
#include <numeric>
#include <queue>
#include <set>
#include <sstream>
#include <stack>
#include <string>

#include <ext/numeric>
#include <ext/functional>

using namespace std;
using namespace __gnu_cxx;

int d[1003];
int p[1003];

int a[1003];
int n;

vector<pair <int, int> > way;

int odd(int n);

int even(int n) {
    if (n <= 0) return 0;
    int& r = d[n];
    if (r) return r;
    if (n == 2) {
        r = max(a[0], a[1]);
        p[n] = 0;
        return r;
    }

    int smallest = even(n-2) + max(a[n-2], a[n-1]);
    int k = n-2;
    int sum = 0;
    for (int i = n-2; i >= 0; i-=2) {
        int cur = even(i) + sum + max(a[i], a[n-1]); 
        sum += max(a[i], a[i-1]); 
        if (cur < smallest) {
            smallest = cur;
            k = i;
        }
    }
    sum = 0;
    for (int i = n-3; i >= 0; i-=2) {
        int cur = odd(i) + sum + max(a[i], a[n-1]); 
        sum += max(a[i], a[i+1]); 
        if (cur < smallest) {
            smallest = cur;
            k = i;
        }
    }
    r = smallest;
    p[n] = k;
    return r;
}

int odd(int n) {
    if (n == 1) return max(a[0], a[2]);
    int& r = d[n];
    if (r) return r;

    int smallest = even(n-1) + max(a[n-1], a[n+1]);
    int k = n-1;
    int sum = 0;
    for (int i = n-1; i >= 0; i-=2) {
        int cur = even(i) + sum + max(a[i], a[n+1]); 
        sum += max(a[i], a[i-1]); 
        if (cur < smallest) {
            smallest = cur;
            k = i;
        }
    }
    sum = 0;
    for (int i = n-2; i >= 0; i-=2) {
        int cur = odd(i) + sum + max(a[i], a[n+1]); 
        sum += max(a[i], a[i+1]); 
        if (cur < smallest) {
            smallest = cur;
            k = i;
        }
    }
    r = smallest;
    p[n] = k;
    return r;
}

void print_even(int n);

void print_odd(int n) {
    if (n <= 0) return;
    int k = p[n];
    if (k % 2) {
        print_odd(k);
        for (int i = k+2; i+1 <= n-1; i+=2)
            cout << i+1 << " " << i+2 << endl;
    } else {
        print_even(k);
        for (int i = k+1; i+1 <= n-1; i+=2)
            cout << i+1 << " " << i+2 << endl;
    }
    if (a[n+1]) cout << k+1 << " " << n+1+1 << endl;
    else cout << k+1 << endl;
}

void print_even(int n) {
    if (n <= 0) return;
    int k = p[n];
    if (k % 2) {
        print_odd(k);
        for (int i = k+2; i+1 < n-1; i+=2)
            cout << i+1 << " " << i+2 << endl;
    } else {
        print_even(k);
        for (int i = k+1; i+1 < n-1; i+=2)
            cout << i+1 << " " << i+2 << endl;
    }
    cout << k+1 << " " << n << endl;
}
int main() {
    cin >> n;
    for (int i = 0; i < n; ++i)
        cin >> a[i];

    if (n % 2) { 
        cout << odd(n) << endl;
        print_odd(n);
    } else {
        cout << even(n) << endl;
        print_even(n);
    }

    return 0;
}