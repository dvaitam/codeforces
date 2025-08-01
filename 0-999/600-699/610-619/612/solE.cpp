#include <cstdio>
void gi(int &x){char ch = getchar(); x = 0; while (ch < '0' || ch > '9') ch = getchar(); while (ch >= '0' && ch <= '9') x = x * 10 + ch - 48, ch = getchar();}
void pi(int x){if (x > 9) pi(x / 10); putchar(x % 10 + 48);}
int n;
#define sz 1002000
int p[sz], v[sz], ans[sz];
int cycle[sz];
int tmp[sz], t2[sz];

int main()
{
    int i, j, k, l, r;
    gi(n);
    for (i = 1; i <= n; i++)gi(p[i]);
    for (i = 1; i <= n; i++)
        if(!v[i])
        {
            l = 0;
            for (j = i; !v[j]; j = p[j]) v[tmp[++l] = j] = 1;
            if(l&1)
            {
                for (j = 1; j <= l; j++){
                  k = (j << 1) - 1; if (k > l) k -= l;
                  t2[k] = tmp[j];
                }
            }
            else
            {
                if (cycle[l])
                {
                  j = cycle[l]; r = i;
                  for (k = 1; k <= l + l; k++)
                    if (k & 1) t2[k] = r, r = p[r];
                    else t2[k] = j, j = p[j];
                  cycle[l] = 0; l <<= 1;
                }
                else
                {
                    cycle[l] = i;
                    goto over;
                }
            }
            for (j = 1; j < l; j++) ans[t2[j]] = t2[j + 1];
            ans[t2[l]] = t2[1];
            over:;
        }
        for (i = 1; i <= n; i++) if (ans[i] == 0){printf("-1\n"); return 0;}
        for (i = 1; i <= n; i++) pi(ans[i]), putchar(' ');
    return 0;
}