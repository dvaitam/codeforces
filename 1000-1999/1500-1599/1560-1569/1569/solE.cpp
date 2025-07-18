#include<bits/stdc++.h>

#define int long long

using namespace std;

const int md=998244353;

int k,A,h,w[40],a[40],b[40];

map <int,int> mp1,mp2;

bool flag=0;

int mid;



inline int qsm(int a,int b)

{

	int ret=1;

	for (;b;b>>=1)

	{

		if (b&1) ret=ret*a%md;

		a=a*a%md;

	}

	return ret;

}



void dfs(int t,int lc,int s,int num,int nn)

{

	while (a[lc])

	{

		lc++;

		if (lc>mid) lc=1,s=s/2+1;

	}

	if (s==2)

	{

//		for (int i=1;i<=mid;i++)

//			cout<<a[i]<<" ";

//		cout<<endl;

		mp1[(num+lc*w[1])%md]=nn+(1<<mid);

		mp2[(num+lc*w[2])%md]=nn+(1<<mid);

//		cout<<(num+lc*w[1])%md<<endl;

		return;

	}

	int cnt=lc;

	a[cnt]=s;

	int pp=s,kk=num+cnt*w[s];

	while (a[cnt])

	{

		cnt++;

	}

	cnt++;

	if (cnt>mid) cnt=1,pp=pp/2+1;

	dfs(t+1,cnt,pp,kk%md,nn|(1<<t-1));

	

	

	

	cnt=lc;

	a[cnt]=0;

	cnt++;

	while (a[cnt])

	{

		cnt++;

	}

	int q=cnt;

	a[cnt]=s;

	num+=cnt*w[s];

	cnt++;

	if (cnt>mid) cnt=1,s=s/2+1;

	dfs(t+1,cnt,s,num%md,nn);

	a[q]=0;

}



void work(int x)

{

	int kk=x;

//	cout<<x<<endl;

	int cnt=1,cc=mid+1;

	for (int i=0;i<mid-1;i++)

	{

//		cout<<i<<endl;

		while (b[cnt])

		{

			cnt++;

			if (cnt>mid) cnt=1,cc=cc/2+1;

		}

		if (kk&(1<<i))

		{

			b[cnt]=cc;

			while (b[cnt]) cnt++;

			cnt++;

			if (cnt>mid) cnt=1,cc=cc/2+1;

		}

		else

		{

			cnt++;

			while (b[cnt]) cnt++;

			b[cnt]=cc;

		}

	}

//	for (int i=1;i<=mid;i++)

//		cout<<b[i]<<" ";

//	cout<<endl;

	

	

	

	for (int i=mid+1;i<=mid<<1;i++)

	{

		b[i]=a[i];

	}

}



void ds(int t,int lc,int s,int num)

{

	if (flag) return;

	while (a[lc])

	{

		lc++;

		if (lc>mid<<1) lc=mid+1,s=s/2+1;

	}

	if (s==2)

	{

//		for (int i=mid+1;i<=mid<<1;i++)

//		{

//			cout<<a[i]<<" ";

//		}

//		cout<<endl;

		if (mp2[(h-(num+lc*w[1])%md+md)%md])

		{

			flag=1;

//			cout<<"!!!!!!!!"<<endl;

			work(mp2[(h-(num+lc*w[1])%md+md)%md]);

//			cout<<"!!!!!!!!"<<endl;

			bool bo=0;

			for (int i=1;i<=mid<<1;i++)

			{

//				cout<<b[i]<<" ";

				if (!b[i]&&!bo)

				{

					bo=1;

					b[i]=2;

				}

				else if (!b[i])

				{

					b[i]=1;

				}

			}

//			cout<<endl;

			return;

		}

//		cout<<(num+lc*w[2])%md<<endl;

		if (mp1[(h-(num+lc*w[2])%md+md)%md])

		{

			flag=1;

			work(mp1[(h-(num+lc*w[2])%md+md)%md]);

			bool bo=0;

			for (int i=1;i<=mid<<1;i++)

			{

				if (!b[i]&&!bo)

				{

					bo=1;

					b[i]=1;

				}

				else if (!b[i])

				{

					b[i]=2;

				}

			}

			return;

		}

		return;

	}

	int cnt=lc;

	a[cnt]=s;

	int pp=s,kk=num+cnt*w[s];

	while (a[cnt])

	{

		cnt++;

	}

	cnt++;

	if (cnt>mid<<1) cnt=mid+1,pp=pp/2+1;

	ds(t+1,cnt,pp,kk%md);

	

	if (flag) return;

	

	cnt=lc;

	a[cnt]=0;

	cnt++;

	while (a[cnt])

	{

		cnt++;

	}

	int q=cnt;

	a[cnt]=s;

	num+=cnt*w[s];

	cnt++;

	if (cnt>mid<<1) cnt=mid+1,s=s/2+1;

	ds(t+1,cnt,s,num%md);

	a[q]=0;

}



signed main()

{

	scanf("%lld%lld%lld",&k,&A,&h);

	mid=qsm(2,k-1);

	w[0]=1;

	for (int i=1;i<=mid<<1;i++)

	{

		w[i]=w[i-1]*A%md;

	}

	if (k==1)

	{

		if ((w[1]+2*w[2])%md==h)

		{

			puts("1 2");

		}

		else if ((2*w[1]+w[2])%md==h)

		{

			puts("2 1");

		}

		else puts("-1");

		return 0;

	}

	dfs(1,1,mid+1,0,0);

//	cout<<"$%^&*("<<endl;

	ds(1,mid+1,mid+1,0);

	if (flag)

	{

		for (int i=1;i<=mid<<1;i++)

		{

			printf("%d ",b[i]);

		}

	}

	else puts("-1");

	return 0;

}