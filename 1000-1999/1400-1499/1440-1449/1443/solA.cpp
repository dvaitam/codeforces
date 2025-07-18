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
    int N;
    cin >> N;
    for (int i = 0; i < N; i++)
        cout << 4 * N - (i + 1) * 2 << '\n';
}

int32_t main() {
    ios::sync_with_stdio(false);
#ifndef imtiyazrasool92
    cin.tie(nullptr);
#endif

    int tests = 1;
    cin >> tests;

    while (tests--) {
        run_case();
        cout << '\n';
    }

    return 0;
}