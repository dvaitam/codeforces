#include <algorithm>
#include <array>
#include <bitset>
#include <cassert>
#include <chrono>
#include <cmath>
#include <cstring>
#include <functional>
#include <iomanip>
#include <iostream>
#include <limits>
#include <map>
#include <numeric>
#include <queue>
#include <random>
#include <set>
#include <vector>

using namespace std;

#ifdef imtiyazrasool92
#include "algos/debug.h"
#else
#define dbg(...) 92
#endif

void run_case() {
    int N, M;
    cin >> N >> M;

    vector<int> A(M);
    for (auto &a : A)
        cin >> a;

    if (accumulate(A.begin(), A.end(), 0LL) < N) {
        cout << -1;
        return;
    }

    for (int i = 0; i < M; i++)
        if (i + A[i] > N) {
            cout << -1;
            return;
        }

    vector<int> out(M + 1);
    iota(out.begin(), out.end(), 1);
    out[M] = N + 1;

    for (int i = M - 1; i >= 0; i--) {
        int position = out[i] + A[i] - 1;
        int move = max(0, out[i + 1] - position - 1);
        out[i] += move;
    }

    for (int i = 0; i < M; i++)
        cout << out[i] << ' ';
}

int32_t main() {
    ios::sync_with_stdio(false);
#ifndef imtiyazrasool92
    cin.tie(nullptr);
#endif

    int tests = 1;
    // cin >> tests;

    while (tests--) {
        run_case();
        cout << '\n';
    }

    return 0;
}