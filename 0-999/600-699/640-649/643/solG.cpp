#include<cstdio>

using namespace std;

const int N=150015,D=5;

int i,j,k,n,m,p,ch,o,l,r,x;

int a[N];

struct tree

{

	int num,same,a[D],b[D];

	tree operator + (const tree &n) const

	{

		tree t;

		t.same=0;

		t.num=num;

		int i,j,k;

		for (i=0;i<num;i++) t.a[i]=a[i],t.b[i]=b[i];

		for (i=0;i<n.num;i++)

		{

			for (j=0;j<t.num;j++) if (t.a[j]==n.a[i])

			{

				t.b[j]+=n.b[i];

				break;

			}

			if (j<t.num) continue;

			if (t.num<p)

			{

				t.a[t.num]=n.a[i];

				t.b[t.num]=n.b[i];

				t.num++;

				continue;

			}

			k=0;

			for (j=1;j<t.num;j++) if (t.b[j]<t.b[k]) k=j;

			if (n.b[i]<t.b[k]) 

			{

				for (j=0;j<t.num;j++) t.b[j]-=n.b[i];

			}

			else

			{

				int tmp=t.b[k];

				t.a[k]=n.a[i];

				t.b[k]=n.b[i];

				for (j=0;j<t.num;j++) t.b[j]-=tmp;

			}

		}

		return t;

	}

} ans,T[N<<2];

void R(int &x)

{

	x=0;ch=getchar();

	while (ch<'0' || '9'<ch) ch=getchar();

	while ('0'<=ch && ch<='9') x=x*10+ch-'0',ch=getchar();

}

void T_build(int L,int R,int k)

{

	if (L==R)

	{

		T[k].num=1;

		T[k].a[0]=a[L];

		T[k].b[0]=1;

		return;

	}

	int mid=(L+R)>>1;

	T_build(L,mid,k<<1);

	T_build(mid+1,R,k<<1|1);

	T[k]=T[k<<1]+T[k<<1|1];

}

void down(int L,int R,int k)

{

	if (T[k].same)

	{

		int l=k<<1,r=l|1,mid=(L+R)>>1;

		T[l].same=T[r].same=T[k].same;

		T[l].num=T[r].num=1;

		T[l].a[0]=T[r].a[0]=T[k].same;

		T[l].b[0]=mid-L+1;

		T[r].b[0]=R-mid;

		T[k].same=0;

	}

}

void T_same(int L,int R,int l,int r,int x,int k)

{

	if (L==l && R==r)

	{

		T[k].same=x;

		T[k].num=1;

		T[k].a[0]=x;

		T[k].b[0]=R-L+1;

		return;

	}

	down(L,R,k);

	int mid=(L+R)>>1;

	if (r<=mid) T_same(L,mid,l,r,x,k<<1);

	else

	{

		if (l>mid) T_same(mid+1,R,l,r,x,k<<1|1);

		else T_same(L,mid,l,mid,x,k<<1),T_same(mid+1,R,mid+1,r,x,k<<1|1);

	}

	T[k]=T[k<<1]+T[k<<1|1];

}

tree T_query(int L,int R,int l,int r,int k)

{

	if (L==l && R==r) return T[k];

	down(L,R,k);

	int mid=(L+R)>>1;

	if (r<=mid) return T_query(L,mid,l,r,k<<1);

	if (l>mid) return T_query(mid+1,R,l,r,k<<1|1);

	return T_query(L,mid,l,mid,k<<1)+T_query(mid+1,R,mid+1,r,k<<1|1);

}

int main()

{

	R(n);R(m);R(p);p=100/p;

	for (i=1;i<=n;i++) R(a[i]);

	T_build(1,n,1);

	for (i=1;i<=m;i++)

	{

		R(o);R(l);R(r);

		if (l==85 && r==85)

		{

			i=i;

		}

		if (o==1)

		{

			R(x);

			T_same(1,n,l,r,x,1);

		}

		else

		{

			ans=T_query(1,n,l,r,1);

			printf("%d",ans.num);

			for (j=0;j<ans.num;j++) printf(" %d",ans.a[j]);

			puts("");

		}

	}

}