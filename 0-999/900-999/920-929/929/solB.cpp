//FIvanO
#include <iostream>
#include <fstream>
#include <sstream>
#include <cmath>
#include <math.h>
#include <algorithm>
#include <cstring>
#include <stdlib.h>
#include <stdio.h>
#include <vector>
#include <map>
#include <set>
#include <stack>
#include <queue>
#include <deque>

using namespace std;

typedef unsigned long long ULL;

int main()
{

    ios_base::sync_with_stdio(0);

    int n, k;
    cin >> n >> k;
    int ans = 0;
    int l = 12;

    string s[105];

    for (int i = 0; i < n; ++i) {
        cin >> s[i];
        if (s[i][0] == '.' && s[i][1] != 'S' && k > 0) {
            k--;
            s[i][0] = 'x';
        }
        if (s[i][11] == '.' && s[i][10] != 'S' && k > 0) {
            k--;
            s[i][11] = 'x';
        }
        for (int j = 1; j < 11; ++j) {
            if (s[i][j] == '.' && s[i][j + 1] != 'S' && s[i][j - 1] != 'S' && k > 0) {
                k--;
                s[i][j] = 'x';
            }
        }
    }

    for (int i = 0; i < n; ++i) {
        if (s[i][0] == '.' && k > 0) {
            k--;
            s[i][0] = 'x';
        }
        if (s[i][11] == '.' && k > 0) {
            k--;
            s[i][11] = 'x';
        }
        if (k == 0) break;
        for (int j = 1; j < 11; ++j) {
            if (s[i][j] == '.' && k > 0) {
                if (s[i][j - 1] == 'S' && s[i][j + 1] == 'S') continue;
                else {
                    k--;
                    s[i][j] = 'x';
                }
            }
        }
    }

    for (int i = 0; i < n; ++i) {
        for (int j = 0; j < l; ++j) {
            if (s[i][j] == '.' && k > 0) {
                k--;
                s[i][j] = 'x';
            }
        }
        if (k == 0) break;
    }

    for (int i = 0; i < n; ++i) {
        if (s[i][0] == 'S') {
            if (s[i][1] != '.' && s[i][1] != '-') ans++;
        }
        if (s[i][11] == 'S') {
            if (s[i][10] != '.' && s[i][10] != '-') ans++;
        }
        for (int j = 1; j < 11; ++j) {
            if (s[i][j] == 'S') {
                if (s[i][j - 1] != '.' && s[i][j - 1] != '-') ans++;
                if (s[i][j + 1] != '.' && s[i][j + 1] != '-') ans++;
            }
        }
    }
    cout << ans << endl;
    for (int i = 0; i < n; ++i) cout << s[i] << endl;

    return 0;
}