#include<cstdio>

#include<algorithm>

using namespace std;

const int N=4,M=1010,oo=1000000000;

int i,j,k,l,n,m,Xn,Yn,Dn,t,ans,nm,Fg,xx,yy,lx,ly,rx,ry;

int X[M],Y[M],D[M],b[24][4],c[4][2],Ans[4][2];

struct cc {int x,y;} A[N];

void pre()

{

	for (i=0;i<4;i++)

		for (j=0;j<4;j++) if (i!=j)

			for (k=0;k<4;k++) if (i!=k && j!=k)

				for (l=0;l<4;l++) if (i!=l && j!=l && k!=l)

				{

					b[nm][0]=i;

					b[nm][1]=j;

					b[nm][2]=k;

					b[nm][3]=l;

					nm++;

				}

	c[0][0]=0;c[0][1]=0;

	c[1][0]=1;c[1][1]=0;

	c[2][0]=0;c[2][1]=1;

	c[3][0]=1;c[3][1]=1;

}

int abs(int x)

{

	if (x<0) return -x;

	return x;

}

void Js(int d,int x,int y)

{

	int i,j,k,t,tt,fg,xx,yy;

	for (i=0;i<nm;i++)

	{

		fg=t=0;

		for (j=0;j<4;j++)

		{

			xx=x+c[b[i][j]][0]*d;

			yy=y+c[b[i][j]][1]*d;

			if (xx!=A[j].x && yy!=A[j].y)

			{

				fg=1;

				break;

			}

			tt=abs(xx-A[j].x)+abs(yy-A[j].y);

			if (tt>t) t=tt;

		}

		if (!fg && t<ans)

		{

			ans=t;

			for (j=0;j<4;j++)

			{

				Ans[j][0]=x+c[b[i][j]][0]*d;

				Ans[j][1]=y+c[b[i][j]][1]*d;

			}

		}

	}

}

int main()

{

	pre();

	scanf("%d",&t);

	while (t--)

	{

		ans=oo;

		Dn=0;

		for (i=0;i<4;i++) scanf("%d%d",&A[i].x,&A[i].y);

		for (i=0;i<4;i++)

			for (j=i+1;j<4;j++) D[++Dn]=abs(A[i].x-A[j].x),D[++Dn]=abs(A[i].y-A[j].y);

		sort(D+1,D+Dn+1);

		D[Dn+1]=-1;

		l=-1;

		for (i=0;i<=Dn;i++) if (D[i]!=D[i+1]) D[++l]=D[i];

		Dn=l;

		for (i=1;i<=Dn;i++)

		{

			Xn=Yn=0;

			for (j=0;j<4;j++)

			{

				X[++Xn]=A[j].x;

				X[++Xn]=A[j].x-D[i];

				X[++Xn]=A[j].x+D[i];

				Y[++Yn]=A[j].y;

				Y[++Yn]=A[j].y-D[i];

				Y[++Yn]=A[j].y+D[i];

			}

			for (j=0;j<nm;j++)

			{

				lx=ly=oo;rx=ry=-oo;

				for (k=0;k<4;k++)

				{

					xx=A[k].x-c[b[j][k]][0]*D[i];

					yy=A[k].y-c[b[j][k]][1]*D[i];

					if (xx<lx) lx=xx;

					if (xx>rx) rx=xx;

					if (yy<ly) ly=yy;

					if (yy>ry) ry=yy;

				}

				X[++Xn]=(lx+rx)/2;

				Y[++Yn]=(ly+ry)/2;

			}

			sort(X+1,X+Xn);

			X[Xn+1]=oo;

			l=0;

			for (j=1;j<=Xn;j++) if (X[j]!=X[j+1]) X[++l]=X[j];

			Xn=l;

			sort(Y+1,Y+Yn);

			Y[Yn+1]=oo;

			l=0;

			for (j=1;j<=Yn;j++) if (Y[j]!=Y[j+1]) Y[++l]=Y[j];

			Yn=l;

			for (j=1;j<=Xn;j++)

				for (k=1;k<=Yn;k++) Js(D[i],X[j],Y[k]);

		}

		if (ans<oo)

		{

			printf("%d\n",ans);

			for (i=0;i<4;i++) printf("%d %d\n",Ans[i][0],Ans[i][1]);

		}

		else puts("-1");

	}

}