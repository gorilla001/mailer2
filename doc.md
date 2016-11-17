##build
```bash
make product
```
## run
```bash
make deploy
```
broswer visit `http://localhost/api`

## CLI

### load
#### server
```bash
mailer load server --host=smtp.126.com --port=25  --user=eyou_uetest@126.com --password=test123
mailer load server --host=smtp.163.com --port=25  --user=eyou_uetest@163.com --password=eyou-uetest
```


#### recipient
```bash
mailer load recipient  --name=rec1 --emails=root_bbk@sohu.com,root_bbk@126.com,root_bbk@163.com
```

#### mail
```bash
mailer load mail --from-name="任超奇" --subject="xxx" --body="yyy"
```

### show
```bash
mailer show server
mailer show recipient
mailer show mail
```

### remove
```bash
mailer rm server --id=58288bd168ec1b6c69000001
mailer rm recipient --id=58288bf768ec1b6c97000001
mailer rm mail --id=58288c1568ec1b6cab000001
```

### task
#### create
```bash
mailer task create --recipient=5828881368ec1b650b000001 --servers=5828849468ec1b62ef000001,5828879568ec1b6499000001 --mails=582889b568ec1b6b58000001
```

#### show
```bash
mailer task show
```

#### rm 
```bash
mailer task rm --id=5829b1235cfa91fdfe68bf84
```

#### run
```bash
mailer task run --id=582b24a168ec1b13a9000001
curl -XPATCH http://127.0.0.1/api/task/run?id=582b24a168ec1b13a9000001
```
