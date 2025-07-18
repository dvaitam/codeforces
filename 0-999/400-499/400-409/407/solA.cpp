#include <stdio.h>
#include <iostream>

using namespace std;

int gcd(int a, int b) {
    if (a == 0)
        return b;
    return gcd(b%a, a);
}

int main() {
    int a, b, ans_x, ans_y, flag = 0;
    scanf("%d%d", &a, &b);
    int x = gcd(a, b);
    if (a > b)
        swap(a, b);
    int num = b/x;
    int den = a/x;
    if (x == 1) {
        cout << "NO" << endl;
    }
    else {
        for (int i = 1; i < a; i++) {
            for (int j = i; j < a; j++) {
                if (j*j+i*i == a*a) {
                    ans_x = i;
                    ans_y = j;
                    if ((ans_x*num) % den || (ans_y*num)%den) {
                        continue;
                    }
                    flag = true;
                    break;
                }
            }
            if (flag)
                break;
        }
        if (!flag) {
            cout << "NO" << endl;
//            cout << (ans_x) << " " << den << " " << num << " " << ans_y << endl;
        }
        else {
            int ans1_x = -(num*ans_y)/den;
            int ans1_y = (num*ans_x)/den;
            if (ans1_x == ans_x || ans1_y == ans_y){
                ans1_x = -ans1_x;
                ans1_y = -ans1_y;
            }
            if (ans_x == ans1_x || ans_y == ans1_y) {
                cout << "NO" << endl;
            }
            else {
                cout << "YES" << endl;
                printf("0 0\n");
                printf("%d %d\n", ans_x, ans_y);
                printf("%d %d\n", ans1_x, ans1_y);       
            }
        }
    }
    return 0;    
}