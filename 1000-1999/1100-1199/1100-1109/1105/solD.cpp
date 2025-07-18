#pragma comment(linker, "/STACK: 40000000")

#include <iostream>
#include <algorithm>
#include <string>
#include <cstring>
#include <complex>
#include <queue>
#include <vector>
#include <tuple>
#include <set>
#include <map>
#include <assert.h>
#include <stdlib.h>
#include <stdio.h>

#define M_PI 3.14159265358979323846
#define INF 1000000000000000099LL
#define SIZE 1005
#define SEP -77
#define ENCODE 1024
#define MOD 1000000007LL

#define DBG false
#define lli long long int
#define mp make_pair
#define cpxd complex<double>
#define compCount 1000001

#define DLOG(x) cout << #x << ": " << x << endl;

//#define int long long int
#define double long double

#define EPS 0.000000000001

using namespace std;

string mtx[SIZE];
queue<int> bfs_queue[10]; /// player => queue with delimiter SEP
int stages_per_turn[10];
int territory[10];
bool player_locked[10];
int locked_players;
int n, m, p;

void fail() {
    cout << "0" << endl;
    exit(0);
}

void mtxlog() {
    return;

    cout << endl;
    for (int i = 0; i < n; i++) {
        cout << mtx[i] << endl;
    }
    cout << endl;
}

void expand_once(int pidx) {
    if (bfs_queue[pidx].front() == SEP) {
        player_locked[pidx] = true;
        locked_players++;
        return;
    }

    while (bfs_queue[pidx].front() != SEP) {
        int cur = bfs_queue[pidx].front();
        int tmp;
        bfs_queue[pidx].pop();

        int row = cur / ENCODE;
        int col = cur % ENCODE;

        for (int sj = 0, srow = 1, scol = 0; sj < 4; tmp = srow, srow = scol, scol = -tmp, sj++) {
            int nrow = row+srow;
            int ncol = col+scol;

            if (ncol < 0 || nrow < 0) continue;
            if (ncol >= m || nrow >= n) continue;
            if (mtx[nrow][ncol] != '.') continue;

            mtx[nrow][ncol] = ('0' + pidx);
            bfs_queue[pidx].push(nrow * ENCODE + ncol);
        }
    }

    /// reset delimiter
    bfs_queue[pidx].pop();
    bfs_queue[pidx].push(SEP);

    mtxlog();
}

void solve() {
    cin >> n >> m >> p;
    locked_players = 0;

    for (int i = 0; i < 10; i++ ) {
        player_locked[i] = true;
    }

    for (int i = 1; i <= p; i++ ) {
        cin >> stages_per_turn[i];
    }

    for (int i = 0; i < n; i++) {
        cin >> mtx[i];
    }

    for (int row = 0; row < n; row++) {
        for (int col = 0; col < m; col++) {
            if ('1' <= mtx[row][col] && mtx[row][col] <= '9') {
                int playernum = mtx[row][col] - '0';
                bfs_queue[playernum].push(row * ENCODE + col);
            }
        }
    }

    for (int i = 1; i <= p; i++ ) {
        player_locked[i] = false;
        bfs_queue[i].push(SEP);
    }

    /// solve

    while (locked_players < p) {
        for (int i = 1; i <= p; i++) {
            for (int stg = 0; (!player_locked[i]) && (stg < stages_per_turn[i]); stg++) {
                expand_once(i);
            }
        }
    }

    /// collect ans
    for (int row = 0; row < n; row++) {
        for (int col = 0; col < m; col++) {
            if ('1' <= mtx[row][col] && mtx[row][col] <= '9') {
                int playernum = mtx[row][col] - '0';
                territory[playernum]++;
            }
        }
    }

    for (int i = 1; i <= p; i++) {
        cout << territory[i] << " ";
    }
    cout << endl;
}

#undef int
int main()
{
    //freopen("file.in", "r", stdin);
    //freopen("file.out", "w", stdout);

    cin.sync_with_stdio(false);
    cin.tie(0);
    cout.tie(0);
    cout.precision(20);

    solve();

    return 0;
}


/*

1 2 3 4 5
0 1 2 3 0




*/