#include <cstdio>

#include <cstring>



using namespace std;



const char fuck[100][3]={"H","He","Li","Be","B","C","N","O","F","Ne","Na","Mg","Al","Si","P","S","Cl","Ar","K","Ca","Sc","Ti","V","Cr","Mn","Fe","Co","Ni","Cu","Zn","Ga","Ge","As","Se","Br","Kr","Rb","Sr","Y","Zr","Nb","Mo","Tc","Ru","Rh","Pd","Ag","Cd","In","Sn","Sb","Te","I","Xe","Cs","Ba","La","Ce","Pr","Nd","Pm","Sm","Eu","Gd","Tb","Dy","Ho","Er","Tm","Yb","Lu","Hf","Ta","W","Re","Os","Ir","Pt","Au","Hg","Tl","Pb","Bi","Po","At","Rn","Fr","Ra","Ac","Th","Pa","U","Np","Pu","Am","Cm","Bk","Cf","Es","Fm"};



int n,m;

int seq[20],tar[20];

char tmp[10];

int link[140000];

struct data{

	int id,val;

	data()

	{

		id=0,val=0;

	}

	data(int _id,int _val)

	{

		id=_id,val=_val;

	}

	bool operator < (const data &n1) const

	{

		return id==n1.id ? val<n1.val : id<n1.id;

	}

	bool operator == (const data &n1) const

	{

		return id==n1.id&&val==n1.val;

	}

	data operator + (const int &p) const

	{

		if (tar[id]-val<p)

			return data(-1,0);

		if (tar[id]-val==p)

			return data(id+1,0);

		if (tar[id]-val>p)

			return data(id,val+p);

	}

}f[140000];



int trans()

{

	scanf("%s",tmp);

	for (int i=0;i<100;i++)

		if (strcmp(tmp,fuck[i])==0)

			return i+1;

	return -1;

}

int main()

{

	scanf("%d%d",&n,&m);

	for (int i=0;i<n;i++)

		seq[i]=trans();

	for (int i=0;i<m;i++)

		tar[i]=trans();

	memset(link,0xff,sizeof(link));

	int lim=(1<<n)-1;

	for (int i=0;i<=lim;i++)

		for (int j=0;j<n;j++)

			if (!(i&(1<<j)))

				if (f[i|(1<<j)]<f[i]+seq[j])

				{

					f[i|(1<<j)]=f[i]+seq[j];

					link[i|(1<<j)]=i;

				}

	if (f[lim]==data(m,0))

	{

		printf("YES\n");

		int now=lim,tmp=0,v=0;

		while (now)

		{

			tmp=now-link[now];

			for (int i=0;i<n;i++)

				if (tmp&(1<<i))

					printf(v++ ? "+%s" : "%s",fuck[seq[i]-1]);

			if (!f[link[now]].val)

			{

				v=0;

				printf("->%s\n",fuck[tar[f[link[now]].id]-1]);

			}

			now=link[now];

		}

	}

	else

		printf("NO\n");

	return 0;

}