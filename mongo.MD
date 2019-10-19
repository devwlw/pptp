vim /etc/yum.repos.d/mongodb-org-4.2.repo

```
[mongodb-org-4.2]
name=MongoDB Repository
baseurl=https://repo.mongodb.org/yum/redhat/$releasever/mongodb-org/4.2/x86_64/
gpgcheck=1
enabled=1
gpgkey=https://www.mongodb.org/static/pgp/server-4.2.asc
```

sudo yum install -y mongodb-org
sudo yum install -y mongodb-org-4.2.0 mongodb-org-server-4.2.0 mongodb-org-shell-4.2.0 mongodb-org-mongos-4.2.0 mongodb-org-tools-4.2.0

vim /etc/mongod.conf
change bindIp  to 0.0.0.0


mongod --config /etc/mongod.conf
```
use admin
db.createUser(
  {
    user: "root",
    pwd:  "abcd123456.",
    roles: [ { role: "root", db: "admin" }]
  }
)

```

```
use mail
db.createUser(
  {
    user: "mailUser",
    pwd:  "abcd1234.",
    roles: [ { role: "readWrite", db: "mail" },{role:"dbAdmin",db:"mail"}]
  }
)
```

db.adminCommand( { shutdown: 1 } )

mongod --auth --config /etc/mongod.conf

mongo -u "mailUser" -p "abcd1234." --authenticationDatabase "mail"