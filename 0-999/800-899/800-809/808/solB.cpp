#include<bits/stdc++.h>
#define int long long

using namespace std;

const int maxN = 3e5;

int readInt ()
{
    bool minus = false;
    int result = 0;
    char ch;
    ch = getchar();
    while (true)
    {
        if (ch == '-') break;
        if (ch >= '0' && ch <= '9') break;
        ch = getchar();
    }
    if (ch == '-') minus = true;
    else result = ch-'0';
    while (true)
    {
        ch = getchar();
        if (ch < '0' || ch > '9') break;
        result = result*10 + (ch - '0');
    }
    if (minus)
        return -result;
    else
        return result;
}
int a[maxN + 1];
main() {
    int n = readInt(), k = readInt();
    int s = 0;
    for(int i = 1; i <= n; i++) {
        a[i] = readInt();
    }
    int res = 0;
    for(int i = 1; i <= k; i++)
        s += a[i];
    res = s;
    for(int i = k + 1; i <= n; i++) {
        s -= a[i - k];
        s += a[i];
        res += s;
    }
    printf("%.6f", 1.0 * res / (n - k + 1));
}