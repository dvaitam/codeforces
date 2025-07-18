#include <iostream>
#include <vector>
#include <algorithm>
#include <set>
#include <map>
#include <string>
using namespace std;
long long p = 998244353;
map<long long, long long> ans;

template <typename T>
vector<T> fetch(vector<T> a, long long n) {
    for (long long i = 0; i < n; ++i) {
        T z;
        cin >> z;
        a.push_back(z);
    }
    return a;
}

vector<pair<int, int>> operations;

void solve_recursively(int n, int offset);

void permute_recursively(int n, int offset) {
    if (n == 1) return;
    solve_recursively(n-1, offset);
    if (n > 2) {
        pair<int, int> p(1 + offset, n - 2 + offset);
        operations.push_back(p);
    }
    permute_recursively(n-1, offset + 1);
}


void solve_recursively(int n, int offset) {
    permute_recursively(n, offset);
    pair<int, int> p(0 + offset, n - 1 + offset);
    operations.push_back(p);
}

void solve() {
    int n;
    cin >> n;
    vector<int> arr;
    arr = fetch(arr, n);
    vector<int> mask(n, 0);
    int max_sum = 0;
    vector<int> max_mask = mask;
    while(mask[0] < 2) {
        int sum = 0;
        int streak = 0;
        for (int i = 0; i < n; ++i) {
            if (mask[i] == 0) {
                ++streak;
            } else {
                sum += streak * streak;
                sum += arr[i];
                streak = 0;
            }
        }
        sum += streak * streak;
        if (sum > max_sum) {
            max_sum = sum;
            max_mask = mask;
        }
        int last = mask.size() - 1;
        ++mask[last];
        while(last > 0 && mask[last] == 2) {
            mask[last] = 0;
            --last;
            ++mask[last];
        }
    }

    cout << max_sum << " ";
    int prev = -1;
    for (int i = 0; i < n; ++i) {
        if (max_mask[i] == 0 && arr[i] != 0) {
            pair<int, int> p(i, i);
            operations.push_back(p);
        }

        if (max_mask[i] == 0 && prev == -1) {
            prev = i;
        } else if (max_mask[i] == 1) {
            if (prev != -1) {
                solve_recursively(i - prev, prev);
            }
            prev = -1;
        }
    }
    if (prev != -1) {
        solve_recursively(n - prev, prev);
    }

    cout << operations.size() << "\n";
    for(auto elem: operations) {
        cout << elem.first + 1 << " " << elem.second + 1 << "\n";
    }
}

void clear_all() {
    operations.clear();
}

int main() {
    int cases = 1;
    //cin >> cases;
    for (int i = 0; i < cases; ++i) {
        clear_all();
        solve();
    }
    return 0;
}