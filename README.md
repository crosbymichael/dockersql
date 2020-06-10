## dockersql

Query your docker with SQL.  Why?  Why not.


### Build From Source Code

You can build from the source code as below:

```bash
git clone https://github.com/crosbymichael/dockersql.git

cd dockersql/cmd/dockersql

sudo go build -o /usr/local/bin/dockersql .

dockersql
```


#### Connect to your docker
```bash
dockersql --docker http://docker:2375
```

#### Get the count of all images
```bash
dockersql
> select count(*) from images;
count(*)
25
>
```

#### Select all images larger than 1G
```bash
dockersql
> select * from images where virtual_size > 1024*1024*1024
id                                                                 parent_id                                                          size       virtual_size   tag
8803adbf8f9c90dd66c146651401aab40a1c12cdbb6fb5bd3934c4ad6f05389d   e16a6ea01ac787c7bebdd73c92037a4fa9ca2bd61979ba00d1831d535aefdaf2   55810869   1196777592     docker:latest
f2505178b040558983f5a54ec32119d1407b68aa936be4942d5a41a199befaf2   e16a6ea01ac787c7bebdd73c92037a4fa9ca2bd61979ba00d1831d535aefdaf2   59775724   1200742447     release:latest
>
```

#### Join containers and images
```bash
dockersql
> SELECT c.id, c.name, c.image FROM containers AS c JOIN images AS i ON c.image=i.tag
id                                                                 name              image
3481862d9bb6e0fd9ad4d951faa7ebf9a477a0194975fb1698a1d3e9691c7935   /redis            crosbymichael/redis:latest
f733305c31fd3e7ea7269a96ff1747f977ee5a19decf8b4778a481a010434a0c   /local-registry   registry:latest
>
```

### License - MIT
