#include <bits/stdc++.h>

#define M 100010



using namespace std;



int X[M],Y[M],SymX[M],SymY[M],tmp[M],A[10],n,hsh[M],Xx[M],Yx[M],T[M],Tx[210][M],cnt=0;



bool cmpx(int x,int y){return X[x]<X[y];}

bool cmpy(int x,int y){return Y[x]<Y[y];}



int check(int x,int y,int goal)

{

	if(Xx[x]==-1||Yx[y]==-1)return 0;

	int s=hsh[x],v=Yx[y];

	for(int i=v;i;i -=i &(-i))goal -=Tx[s][i];

	return goal ==0;

}



int main()

{

	scanf("%d",&n);

	for(int i=1;i<=n;i++)scanf("%d%d",&X[i],&Y[i]);

	for(int i=1;i<=9;i++)scanf("%d",&A[i]);

	sort(A+1,A+10);

	for(int i=1;i<=9;i++)

		for(int j=i+1;j<=9;j++)

			for(int k=j+1;k<=9;k++){

				int tmp=A[i]+A[j]+A[k];

				if(!hsh[tmp])hsh[tmp]=++cnt;

				if(!hsh[n-tmp])hsh[n-tmp]=++cnt;

			}

	for(int i=1;i<=n;i++)tmp[i]=i;

	sort(tmp+1,tmp+n+1,cmpy);

	for(int i=1,t=1;i<=n;i++){

		SymY[t]=Y[tmp[i]];

		if(Y[tmp[i]]==Y[tmp[i+1]])Yx[i]=-1,Y[tmp[i]]=t;

		else Yx[i]=t,Y[tmp[i]]=t++;

	}

	for(int i=1;i<=n;i++)tmp[i]=i;

	sort(tmp+1,tmp+n+1,cmpx);

	for(int i=1,t=1;i<=n;i++){

		SymX[t]=X[tmp[i]];

		if(X[tmp[i]]==X[tmp[i+1]])Xx[i]=-1,X[tmp[i]]=t;

		else Xx[i]=t,X[tmp[i]]=t++;

	}

	for(int i=1;i<=n;i++){

		for(int j=Y[tmp[i]];j<=n;j +=j &(-j))T[j]++;

		if(hsh[i]&& Xx[i]!=-1)memcpy(Tx[hsh[i]],T,sizeof T);

	}

	for(;;){

		int a1=A[1]+A[2]+A[3],a2=a1+A[4]+A[5]+A[6],a3=A[1]+A[4]+A[7],a4=a3+A[2]+A[5]+A[8];

		if(check(a3,a1,A[1])&& check(a4,a1,A[1]+A[2])&& check(a3,a2,A[1]+A[4])&& check(a4,a2,A[1]+A[2]+A[4]+A[5]))

			return printf("%.12lf %.12lf\n%.12lf %.12lf\n",SymX[Xx[a3]]+0.5,SymX[Xx[a4]]+0.5,SymY[Yx[a1]]+0.5,SymY[Yx[a2]]+0.5),0;

		if(!next_permutation(A+1,A+10))return puts("-1"),0;

	}

}