#include<cstdio>

#include<cstring>

#include<iostream>



#define nn 1000010

using namespace std;



long long ans,f[nn][20];

int bin[20],w[20],n,m,sz[nn];



int main()

{

	scanf("%d%d",&n,&m);

	for (int i=0;i<=n;i++) bin[i]=1<<i;

	for (int u,v,i=1;i<=m;i++)

	{

		scanf("%d%d",&u,&v);

		w[u-1]|=bin[v-1],w[v-1]|=bin[u-1];

	}

	

	for (int pos,i=1;i<=bin[n];i++)

	{

		for (pos=0;pos<n;pos++) if (i&bin[pos]) break;

		if (i==bin[pos]) {f[i][pos]=sz[i]=1;continue;}

		

		sz[i]=sz[i^bin[pos]]+1;

		for (int j=pos+1;j<n;j++) if (i&bin[j])

		{

			for (int k=0;k<n;k++) if (w[j]&bin[k])

			f[i][j]+=f[i^bin[j]][k];

			

			if (w[j]&bin[pos]&&sz[i]>2) ans+=f[i][j];

		}

	}

	

	printf("%I64d\n",ans/2ll);

}