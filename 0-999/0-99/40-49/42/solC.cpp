#include <cstdio>
#include <cstring>
#include <cstdlib>

int a[4];

int main() {
    for (int i = 0; i < 4; ++ i) scanf("%d", &a[i]);
    while (1) {
        if (a[0] == 1 && a[1] == 1 && a[2] == 1 && a[3] == 1) break;
        int flag = 0;
        for (int i = 0; i < 4; ++ i) {
            if (a[i] % 2 == 0 && a[(i + 1) % 4] % 2 == 0) {
                a[i] /= 2;
                a[(i + 1) % 4] /= 2;
                printf("/%d\n", i + 1);
                flag = 1;
            }
        }
        if (flag) continue;
        int x = rand() % 4;
        ++ a[x];
        ++ a[(x + 1) % 4];
        printf("+%d\n", x + 1);
    }
    return 0;
}