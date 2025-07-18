#include <cstdio>
using namespace std;

int l[22][22], n;
bool u[1111];

bool hujak(int i, int c) {
}

int main(void) {

    scanf("%d", &n);
    for (int i = 0; i < n; i++) {
        l[i][i] = 0;
    }
    for (int i = 0; i < 1001; i++) {
        u[i] = false;
    }

    l[0][1] = l[1][0] = 1; u[1] = true;
    l[0][2] = l[2][0] = 2; u[2] = true;
    l[1][2] = l[2][1] = 3; u[3] = true;

    for (int i = 3; i < n; i++) {

        for (int c = 1; c < 1001; c++) {

            int j;
            for (j = 0; j < i; j++) {
                if (u[c+l[0][j]]) break;
            }

            if (j < i) continue;

            for (j = 0; j < i; j++) {
                u[c+l[0][j]] = true;
                l[j][i] = l[i][j] = c+l[0][j];
            }
            break;
        }

    }

    for (int i = 0; i < n; i++) {
        for (int j = 0; j < n; j++) {
            printf("%d ", l[i][j]);
        }
        printf("\n");
    }

    return 48-48;
}