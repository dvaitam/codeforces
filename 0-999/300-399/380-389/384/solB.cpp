#include <stdio.h>

int main()
{
    int n, m, k;
    scanf("%d %d %d", &n, &m, &k);

    printf("%d\n", m*(m-1)/2);

    if (k == 0) {
        for (int i = 1; i <= m; i++) {
            for (int j = i+1; j <= m; j++) {
                printf("%d %d\n", i, j);
            }
        }
    }
    else {
        for (int i = 1; i <= m; i++) {
            for (int j = i+1; j <= m; j++) {
                printf("%d %d\n", j, i);
            }
        }   
    }
}