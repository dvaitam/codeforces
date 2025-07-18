/* LittleFall : Hello! */
#include <bits/stdc++.h>
using namespace std; typedef long long ll;
inline int read(); inline void write(int x);
const int M = 15000016, N = 300000;

//欧拉筛-开始
int minFactor[M];//phi[M];
int prime[M], primeNum;
void calPhi()
{
    for (int i = 2; i < M; i++)
    {
        if (!minFactor[i])
        {
            prime[primeNum++] = i;
            minFactor[i] = primeNum;
        }
        for (int j = 1; j <= minFactor[i]; j++)
        {
            ll t = i * prime[j - 1];
            if (t >= M)break;
            minFactor[t] = j;
            if (j == minFactor[i])
                break;
        }
    }
}
//欧拉筛-结束s
int save[N], cnt[M];
int main(void)
{
	#ifdef _LITTLEFALL_
	freopen("in.txt","r",stdin);
    #endif

	calPhi();
	for(int i=1;i<M;i++)
		minFactor[i] = prime[minFactor[i]-1];
	minFactor[1] = 0;
	//printf("%d\n",minFactor[1] );
	int n = read(), gcd = 0;
	for(int i=1;i<=n;i++) save[i] = read();
	for(int i=1;i<=n;i++) gcd = __gcd(gcd,save[i]);
	for(int i=1;i<=n;i++) save[i] /= gcd;

	int ans = 0;
	for(int i=1;i<=n;i++)
	{
		//分解质因数,cnt++
		int t = save[i];
		while(t != 1)
		{
			//printf("t=%d\n",t );
			int x = minFactor[t];
			//printf("x=%d\n",x );
			while(minFactor[t]==x)
				t/=x;
			//printf("xt=%d\n",t );
			ans = max(ans,++cnt[x]);
		}
	}
	printf("%d\n",ans?n-ans:-1 );
    return 0;
}


inline int read()
{
    int x=0,f=1;char ch=getchar();
    while(ch<'0'||ch>'9') {if(ch=='-')f=-1;ch=getchar();}
    while(ch>='0'&&ch<='9'){x=x*10+ch-'0';ch=getchar();}
    return x*f;
}
inline void write(int x)
{
     if(x<0) putchar('-'),x=-x;
     if(x>9) write(x/10);
     putchar(x%10+'0');
}