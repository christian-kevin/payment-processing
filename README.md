# Payment Processing

###Assumption
- DB Sharding is out of scope
- Redis cluster, master slave, sentinel mode, RDB, AOF are out of scope.
- Card data encryption is out of scope.
- Auth service can be separate on different microservice.
- Can use multiple instance of payment-processing service, using loadbalancer in front
- Rate limit is implemented on service, can use nginx/kong/istio ingress + envoy,etc
- Rate limit using simple rate limit for current QPS. can use more complex rate limit algorithm like leaky bucket, sliding window, etc for improvement.
- On user table, plain password is used for simplicity
- Auth method skipped, should be cookies on header then server compared with cookies stored on redis / db
- CVV must not stored in DB (PCI DSS) 
- log and cache library using abstraction that I've implemented before.
- All table is on single database for simplicity. If using multiple database with multiple service will add complexity 
  such as we can't use transaction, consistency need to be enforced using saga / double phased commit.
- Delete card is soft delete.
- I don't know about update card, I think card should be immutable.
- Usually for routes, there is nginx replace for replacing url path. For example `https://spenmo.com/order/v1/status/xxx` will redirected
  to order service and path replaced to `https://public/v1/status/xxxxx`
- Card number generation use Luhn algorithm with IIN of mastercard
- Expiry date will took ID current local date time + 2 years
- When creating balance will automatically add 10 mil IDR credit
- For now currency will be IDR, but db support other currency based on `country` field in `wallet` table
- Every card limit change, limit will reset

###Directory Structure
I am using directory structure guideline from [golang standard](https://github.com/golang-standards/project-layout)
- `cmd` as runmode (webservice, cron-mode, consumer-mode, etc)
- `config` for application config file
- `internal` for private application code, consist of:
    1. `app` this is for api controller and handler (manager in java)
    2. `pkg` this is for application shared component and DAO, request response struct also included here in `dto` dir
- `migration` for history of migration / schema alter. note that all schema alter copied to `init.sql` used by `docker-compose`
- `pkg` for shared  / common library
- `vendor` for dependencies
- `apis` for api documentation

##To Run
1. install `docker-compose`
2. run `docker-compose up --build`
3. api can accessed on `http://localhost:8080/public/{api}`

##Observability Plan
- Can send metrics to prometheus / thanos as storage for every transaction (ex: `create card` / `create wallet`), showed in grafana
- Log already standardized (`context_id`, `message`, `stack_trace`, etc), can add ingress / egress log if using microservice / 3rd party service. 
  Log then scraped from stdout by fluentd using regex to elasticsearch.
- Log in elasticsearch can put on graylog / kibana / even add in grafana for monitoring.
- From kibana / graylog can add alerting when reaches certain threshold, sent to telegram / slack channel /etc

##Note
- GetCardsFromWallet and CreateCard rate limited to 2 req / minute. Can change route / disable in config + router.

##API Docs
Domain: `https://localhost:8080/public`

*CreateWallet*
```
/POST /v1/wallet
Header: 
X-User-ID: 1 ##dummy user id injected in init.sql
X-AU: xxxx ##any value (as cookie auth token)
X-Tenant: id  ## for now wallet only support id / IDR

200 
{
    code: 200, #number
    err_msg: "", #string
    data: {
            "wallet_id": 1, #number
            "balance": "10000000", #string
            "country": "id" #string
    }
}

400
{
    "code": 400,
    "err_msg": "wallet already exist",
    "data": null
}

401
{
    "code": 401,
    "err_msg": "unauthorized",
    "data": null
}
```
*GetWallet*
```
/GET /v1/wallet
Header: 
X-User-ID: 1 ##dummy user id injected in init.sql
X-AU: xxxx ##any value (as cookie auth token)
X-Tenant: id  ## for now wallet only support id / IDR

200 
{
    code: 200, #string
    err_msg: "",
    data: {
            "wallet_id": 1, #number
            "balance": "10000000", #string
            "country": "id" #string
    }
}

500
{
    "code": 500,
    "err_msg": "server error",
    "data": null
}

401
{
    "code": 401,
    "err_msg": "unauthorized",
    "data": null
}
```

*CreateCard*
```
/POST /v1/card
Header: 
X-User-ID: 1 ##dummy user id injected in init.sql
X-AU: xxxx ##any value (as cookie auth token)
X-Tenant: id  ## for now wallet only support id / IDR
Body:
{
    "name": "XXXX" ## name on card (string)
}

200
{
    "code": 200, #number
    "err_msg": "",
    "data": {
        "card_id": 3, #number
        "card_number": "5100001399840608", #string
        "expiry_date": "11/23", #string
        "name": "Kevin" #string
    }
}

500
{
    "code": 500,
    "err_msg": "server error",
    "data": null
}

404
{
    "code": 404,
    "err_msg": "wallet not found",
    "data": null
}

401
{
    "code": 401,
    "err_msg": "unauthorized",
    "data": null
}

429
{
    "code": 429,
    "err_msg": "rate limit exceeded",
    "data": null
}
```

*GetCards*
```
/GET /v1/card/multiple
Header: 
X-User-ID: 1 ##dummy user id injected in init.sql
X-AU: xxxx ##any value (as cookie auth token)
X-Tenant: id  ## for now wallet only support id / IDR

200
{
    "code": 200, #number
    "err_msg": "",
    "data": {
        "cards": [
                    {
                        "card_id": 1, #number
                        "card_number": "5100004984980821", #string
                        "expiry_date": "11/23", #string
                        "name": "Kevin", #string
                        "created_at": 1637735190586 #number
                    },
                    {
                        "card_id": 2,
                        "card_number": "5100006271318482",
                        "expiry_date": "11/23",
                        "name": "Kevin",
                        "created_at": 1637735191138
                    },
                    {
                        "card_id": 3,
                        "card_number": "5100001399840608",
                        "expiry_date": "11/23",
                        "name": "Kevin",
                        "created_at": 1637737509700
                    },
                    {
                        "card_id": 4,
                        "card_number": "5100001119020820",
                        "expiry_date": "11/23",
                        "name": "Kevin",
                        "created_at": 1637737574897
                    },
                    {
                        "card_id": 5,
                        "card_number": "5100006749413196",
                        "expiry_date": "11/23",
                        "name": "Kevin",
                        "created_at": 1637737575582
                    }
                ]
    }
}

500
{
    "code": 500,
    "err_msg": "server error",
    "data": null
}

401
{
    "code": 401,
    "err_msg": "unauthorized",
    "data": null
}

429
{
    "code": 429,
    "err_msg": "rate limit exceeded",
    "data": null
}
```
*DeleteCard*
```
/DELETE /v1/card
Header: 
X-User-ID: 1 ##dummy user id injected in init.sql
X-AU: xxxx ##any value (as cookie auth token)
X-Tenant: id  ## for now wallet only support id / IDR
Body:
{
    "card_id": 1 ## card id got from get multiple card (number)
}

200
{
    "code": 200,
    "err_msg": "",
    "data": null
}

500
{
    "code": 500,
    "err_msg": "server error",
    "data": null
}

404
{
    "code": 404,
    "err_msg": "card not found",
    "data": null
}

401
{
    "code": 401,
    "err_msg": "unauthorized",
    "data": null
}
```

*CreateTransaction*
```
/GET /v1/card/public/transaction
Body
{
    "card_number": "5100004984980821", #string, from get cards / create card
    "expiry_date": "11/23", #string
    "amount": "1000000", #string
    "country": "id" #string
}

200
{
    "code": 200,
    "err_msg": "",
    "data": null
}


500
{
    "code": 500,
    "err_msg": "server error",
    "data": null
}

404
{
    "code": 404,
    "err_msg": "card not found",
    "data": null
}
