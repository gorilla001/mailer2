## install docker engine
[Guide](https://docs.docker.com/engine/installation/linux/)

## run
```bash
git clone https://github.com/tinymailer/mailer
cd mailer/
make deploy
```
this will start up an container `mailer` which will occupy OS port `0.0.0.0:80` and `127.0.0.1:27017`

## CLI Usage 

switch into the mailer container
```bash
docker exec -it mailer /bin/sh
```

### smtp server
```bash
mailer show server
mailer show server --id=58288bd168ec1b6c69000001
mailer load server --host=smtp.126.com --port=25  --user=xyz@126.com --password=xyz
mailer rm   server --id=58288bd168ec1b6c69000001
```


### recipient list
```bash
mailer show recipient
mailer show recipient  --id=58288bf768ec1b6c97000001
mailer load recipient  --name=rec1 --emails=user1@sohu.com,user2@126.com,user3@163.com
mailer load recipient  --name=rec2 --emails=file:///tmp/emails.txt
mailer rm   recipient  --id=58288bf768ec1b6c97000001
```

### mail content
```bash
mailer show mail
mailer show mail --id=58288c1568ec1b6cab000001
mailer load mail --from-name="推广中心" --subject=xxx --body="yyy"
mailer load mail --from-name="推广中心" --subject=xxx --body=file:///tmp/mailbody.html
mailer rm   mail --id=58288c1568ec1b6cab000001
```

### task management
```bash
mailer task show
mailer task show   --id=5829b1235cfa91fdfe68bf84
mailer task create --recipient=5828881368ec1b650b000001 --servers=5828849468ec1b62ef000001,5828879568ec1b6499000001 --mails=582889b568ec1b6b58000001
mailer task run    --id=582b24a168ec1b13a9000001
mailer task run    --id=582b24a168ec1b13a9000001 --max-concurrency=10 --max-retry=1 --wait-delay=10
mailer task rm     --id=5829b1235cfa91fdfe68bf84
```

## HTTP API

### smtp server

#### list all
`GET` `/api/server`

```json
[
  {
    id: "5828849468ec1b62ef000001",
    host: "smtp.126.com",
    port: "25",
    auth_user: "uetest@126.com",
    auth_pass: "test123"
  },
  {
    id: "5828879568ec1b6499000001",
    host: "smtp.163.com",
    port: "25",
    auth_user: "uetest@163.com",
    auth_pass: "uetest"
  },
  {
    id: "5843905168ec1b289f000001",
    host: "smtp.xxx.com",
    port: "25",
    auth_user: "xxx@xxx.com",
    auth_pass: "xxxxxxx"
  }
]
```

#### get specified one
`GET` `/api/server?id=5828879568ec1b6499000001`

#### remove specified
`DELETE` `/api/server?id=5828879568ec1b6499000001`


### recipient list

#### list all
`GET` `/api/recipient`

```json
[
  {
    id: "5828881368ec1b650b000001",
    name: "rec1",
    emails: [
      "xx1@sohu.com",
      "xx2@126.com",
      "xx3@163.com"
    ]
  },
  {
    id: "583813ca68ec1b3e29000001",
    name: "sohu",
    emails: [
      "xx1@sohu.com",
      "xx2@sohu.com",
      "xx3@sohu.com",
      "xx4@sohu.com",
    ]
  }
]
```

#### get specified one
`GET` `/api/recipient?id=5828881368ec1b650b000001`

#### remove specified
`DELETE` `/api/recipient?id=5828881368ec1b650b000001`

### mail content

####  list all
`GET` `/api/mail`

```json
[
  {
    id: "582889b568ec1b6b58000001",
    from_name: "任超奇",
    subject: "部署工作已经完成",
    body: "这里是部署mailer的文档, 包括docker容器方式的部署和启动文档, 请查收"
  },
  {
    id: "5844224f15ec240067000001",
    from_name: "推广中心",
    subject: "xxx",
    body: "yyy"
  },
  {
    id: "5844227915ec24007d000001",
    from_name: "推广中心",
    subject: "xxx",
    body: "<html> <h1>title</h1> <bold>aaaaaaaa</bold> </html>"
  }
]
```

#### get specified one
`GET` `/api/mail?id=582889b568ec1b6b58000001`

#### remove specified
`DELETE` `/api/mail?id=582889b568ec1b6b58000001`


### task management

#### list all
`GET` `/api/task` 

```json
[
  {
    id: "5843908f68ec1b2941000001",
    recipient: "5828881368ec1b650b000001",
    servers: [
      "5828849468ec1b62ef000001",
      "5843905168ec1b289f000001"
    ],
    mails: [
      "5837abad68ec1b5b37000001"
    ],
    status: "created",
    recipient_name: "rec1",
    server_names: [
      "smtp.126.com:25-eyou_uetest@126.com",
      "smtp.xxx.com:25-xxx@xxx.com"
    ]
  },
  {
    id: "584422d415ec240099000001",
    recipient: "583813ca68ec1b3e29000001",
    servers: [
      "5828849468ec1b62ef000001",
      "5828879568ec1b6499000001",
      "5843905168ec1b289f000001"
    ],
    mails: [
      "5844227915ec24007d000001"
    ],
    status: "created",
    recipient_name: "sohu",
    server_names: [
      "smtp.126.com:25-eyou_uetest@126.com",
      "smtp.163.com:25-eyou_uetest@163.com",
      "smtp.xxx.com:25-xxx@xxx.com"
    ]
  }
]
```

#### get specified one
`GET` `/api/task?id=5843908f68ec1b2941000001`

#### remove specified
`DELETE` `/api/task?id=5843908f68ec1b2941000001`

#### run a task
`PATCH` `/api/task/run?id=5843908f68ec1b2941000001`

Request Header:
```liquid
Content-Type: application/json
```

Request Payload:
```json
{
  "max_concurrency": 5,
  "max_retry": 1,
  "wait_delay": 10 
}
```

Response Stream:
```json
{"details":{"finished":false,"mail_id":"5844227915ec24007d000001","mail_to":"xxx1@sohu.com","server_ranks":{"5828849468ec1b62ef000001":{"negative":0,"positive":0,"rank":1},"5828879568ec1b6499000001":{"negative":0,"positive":0,"rank":1},"5843905168ec1b289f000001":{"negative":0,"positive":0,"rank":1}},"smtp_server":"5828849468ec1b62ef000001","start_at":"2016-12-04T15:20:36.817791018Z"}}
{"details":{"finished":true,"mail_id":"5844227915ec24007d000001","mail_to":"xxx1@sohu.com","server_ranks":{"5828849468ec1b62ef000001":{"negative":0,"positive":0,"rank":1},"5828879568ec1b6499000001":{"negative":0,"positive":0,"rank":1},"5843905168ec1b289f000001":{"negative":0,"positive":0,"rank":1}},"smtp_server":"5828849468ec1b62ef000001","start_at":"2016-12-04T15:20:36.817791018Z","time_elapsed":"1.784213816s"},"succ":1}
{"details":{"finished":false,"mail_id":"5844227915ec24007d000001","mail_to":"xxx2@sohu.com","server_ranks":{"5828849468ec1b62ef000001":{"negative":0,"positive":1,"rank":2},"5828879568ec1b6499000001":{"negative":0,"positive":0,"rank":1},"5843905168ec1b289f000001":{"negative":0,"positive":0,"rank":1}},"smtp_server":"5828849468ec1b62ef000001","start_at":"2016-12-04T15:20:39.902752008Z"}}
....
....
....
{"details":{"finished":true,"mail_id":"5844227915ec24007d000001","mail_to":"xxx9@sohu.com","server_ranks":{"5828849468ec1b62ef000001":{"negative":0,"positive":15,"rank":16},"5828879568ec1b6499000001":{"negative":0,"positive":0,"rank":1},"5843905168ec1b289f000001":{"negative":2,"positive":0,"rank":1}},"smtp_server":"5828849468ec1b62ef000001","start_at":"2016-12-04T15:21:32.224996714Z","time_elapsed":"290.338177ms"},"succ":16,"fail":2}
{"details":{"finished":false,"mail_id":"5844227915ec24007d000001","mail_to":"xxx10@sohu.com","server_ranks":{"5828849468ec1b62ef000001":{"negative":0,"positive":16,"rank":17},"5828879568ec1b6499000001":{"negative":0,"positive":0,"rank":1},"5843905168ec1b289f000001":{"negative":2,"positive":0,"rank":1}},"smtp_server":"5828849468ec1b62ef000001","start_at":"2016-12-04T15:21:33.816056949Z"}}
{"details":{"finished":true,"mail_id":"5844227915ec24007d000001","mail_to":"xxx10@sohu.com","server_ranks":{"5828849468ec1b62ef000001":{"negative":0,"positive":16,"rank":17},"5828879568ec1b6499000001":{"negative":0,"positive":0,"rank":1},"5843905168ec1b289f000001":{"negative":2,"positive":0,"rank":1}},"smtp_server":"5828849468ec1b62ef000001","start_at":"2016-12-04T15:21:33.816056949Z","time_elapsed":"1.097625422s"},"succ":17,"fail":2}
{"finish":true,"succ":17,"fail":2,"elapsed":"59.701636061s"}
```
