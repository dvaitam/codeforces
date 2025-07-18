#include <stdio.h>

int main() {
    char s[100];
    int res = 0;
    int num = 0;
    int op = 1;
    int i;
    
    scanf("%s", s);
    //printf("%d %d\n", '-', '0');
    for(i = 0; s[i] != 0; i++) {
        if(s[i] >= '0' && s[i] <= '9') num = num * 10 + s[i] - '0';
        else {
            res += num * op;
            num = 0;
            if(s[i] == '+') op = 1;
            else op = -1;
        }
    }
    res += num * op;
    if(res >= 100) {
        for(i = 0; i < res / 100; i++) printf("+");
        printf("\n++++++++++++++++++++++++++++++++++++++++++++++++.");
        printf(">");
    }
    if(res >= 10) {
        for(i = 0; i < (res % 100) / 10; i++) printf("+");
        printf("\n++++++++++++++++++++++++++++++++++++++++++++++++.");
        printf(">");
    }
    for(i = 0; i < res % 10; i++) printf("+");
    printf("\n++++++++++++++++++++++++++++++++++++++++++++++++.");
    return 0;
}