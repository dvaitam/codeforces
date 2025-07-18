#include<cstdio>  

    #include<cstring>  

    #include<cctype>  

    #include<cstdlib>  

    #include<cmath>  

    #include<iostream>  

    #include<sstream>  

    #include<iterator>  

    #include<algorithm>  

    #include<string>  

    #include<vector>  

    #include<set>  

    #include<map>  

    #include<deque>  

    #include<queue>  

    #include<stack>  

    #include<list>  

    typedef long long ll;  

    typedef unsigned long long llu;  

    const int maxn = 100000 + 10;  

    const int inf = 0x3f3f3f3f;  

    const double pi = acos(-1.0);  

    const double eps = 1e-8;  

    using namespace std;  

      

      

    int vis[10];  

    int a;  

    int main(){  

        int n;  

        while(scanf("%d",&n) == 1){  

            memset(vis,0,sizeof vis);  

            for (int i = 0; i < n; ++i)  

            {  

                scanf("%d",&a);  

                vis[a]++;  

            }  

            if (vis[5] || vis[7]){  

                printf("-1\n");  

                continue;  

            }  

            if (vis[1] != n/3)printf("-1\n");  

            else {  

                if (vis[2] - vis[4] < 0 || vis[2] - vis[4] != vis[6] - vis[3] || vis[6] < vis[3])printf("-1\n");  

                else {  

                    for (int i = 0; i < vis[4]; ++i)printf("1 2 4\n");  

                    for (int i = 0; i < vis[6]-vis[3]; ++i)printf("1 2 6\n");  

                    for (int i = 0; i < vis[3]; ++i)printf("1 3 6\n");  

      

      

                }  

      

            }  

        }  

      

        return 0;  

    }