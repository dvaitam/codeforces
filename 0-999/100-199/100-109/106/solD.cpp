#include <cstdio>
#include <cstring>
#include <algorithm>

using namespace std;

char ss[1100][1100];
int sum[1100][1100];
int posi[300][2];

inline
int get(int x1, int y1, int x2, int y2)
{
    return sum[x2][y2] - sum[x1 - 1][y2] - sum[x2][y1 - 1] + sum[x1 - 1][y1 - 1];
}

int dir[4][2] = {-1, 0, 1, 0, 0, -1, 0, 1};
int d[300];

int main()
{
    int n, m;

    //freopen("input.txt", "r", stdin);

    d['N'] = 0;
    d['S'] = 1;
    d['W'] = 2;
    d['E'] = 3;

    scanf("%d%d", &n, &m);

    for (int i = 1; i <= n; i++) {
        scanf("%s", ss[i] + 1);
    }

    for (int i = 1; i <= n; i++) {
        sum[i][0] = 0;
        for (int j = 1; j <= m; j++) {
            sum[i][j] = sum[i][j - 1];
            if (ss[i][j] == '#')
                sum[i][j]++;
        }
    }

    bool ok[300];
    memset(ok, false, sizeof(ok));
    for (int i = 1; i <= n; i++) {
        for (int j = 1; j <= m; j++) {
            sum[i][j] += sum[i - 1][j];
            if (ss[i][j] >= 'A' && ss[i][j] <= 'Z') {
                posi[ss[i][j]][0] = i;
                posi[ss[i][j]][1] = j;
                ok[ss[i][j]] = true;
            }
        }
    }

    /*for (int i = 1; i <= n; i++) {
        for (int j = 1; j <= m; j++) {
            printf("%d ", sum[i][j]);
        }
        printf("\n");
    }*/

    int k;
    scanf("%d", &k);
    char cmd[200];
    int x;
    for (int i = 0; i < k; i++) {
        scanf("%s%d", cmd, &x);
        for (int j = 'A'; j <= 'Z'; j++) {
            if (ok[j]) {
                int nx = posi[j][0] + dir[d[cmd[0]]][0] * x,
                    ny = posi[j][1] + dir[d[cmd[0]]][1] * x;

                if (nx >= 1 && nx <= n && ny >= 1 && ny <= m) {
                    int ret = get(min(nx, posi[j][0]), min(ny, posi[j][1]),
                                max(nx, posi[j][0]), max(ny, posi[j][1]));
                    if (ret > 0)
                        ok[j] = false;
                } else {
                    ok[j] = false;
                }
                posi[j][0] = nx;
                posi[j][1] = ny;

            }
        }
    }
    int cc = 0;
    for (int i = 'A'; i <= 'Z'; i++) {
        if (ok[i]) {
            printf("%c", i);
            cc++;
        }
    }
    if (cc == 0)
        printf("no solution");
    printf("\n");

    return 0;
}