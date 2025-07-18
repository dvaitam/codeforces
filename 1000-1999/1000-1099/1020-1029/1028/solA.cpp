#include<bits/stdc++.h>
using namespace std;
const int inf = 1e9 + 100;
const int N = 120;
int n;
int m;
char grid[N][N];
int row;
int col;
int maxi_r;
int maxi_c;
int mini_r;
int mini_c;
int main() {
    mini_c = mini_r = inf;
    maxi_c = maxi_r = -inf;
    cin >> n >> m;
    for (int i = 0 ; i < n ; i++) {
        for (int j = 0 ; j < m ; j++) {
            cin >> grid[i][j];
        }
    }
    for (int i = 0 ; i < n ; i++) {
        for (int j = 0 ; j < m ; j++) {
            if(grid[i][j] == 'B') {
                maxi_r = max(i, maxi_r);
                maxi_c = max(j, maxi_c);
                mini_r = min(i, mini_r);
                mini_c = min(j, mini_c);
            }
        }
    }
    row = (mini_r + maxi_r) >> 1;
    col = (mini_c + maxi_c) >> 1;
    cout << row + 1 << " " << col + 1 << endl;
return 0;
}