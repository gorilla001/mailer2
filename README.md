## install docker engine
[Guide](https://docs.docker.com/engine/installation/linux/)

## run
```bash
make deploy
```

## CLI

### server
```bash
mailer show server
mailer load server --host=smtp.126.com --port=25  --user=xyz@126.com --password=xyz
mailer rm   server --id=58288bd168ec1b6c69000001
```


### recipient
```bash
mailer show recipient
mailer load recipient  --name=rec1 --emails=user1@sohu.com,user2@126.com,user3@163.com
mailer load recipient  --name=rec2 --emails=file:///tmp/emails.txt
mailer rm   recipient  --id=58288bf768ec1b6c97000001
```

### mail
```bash
mailer show mail
mailer load mail --from-name="推广中心" --subject=xxx --body="yyy"
mailer load mail --from-name="推广中心" --subject=xxx --body=file:///tmp/mailbody.html
mailer rm   mail --id=58288c1568ec1b6cab000001
```

### task
```bash
mailer task show
mailer task create --recipient=5828881368ec1b650b000001 --servers=5828849468ec1b62ef000001,5828879568ec1b6499000001 --mails=582889b568ec1b6b58000001
mailer task run --id=582b24a168ec1b13a9000001
mailer task rm  --id=5829b1235cfa91fdfe68bf84
```

## HTTP API
