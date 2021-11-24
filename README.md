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
- Can send metrics to prometheus for every transaction (ex: `create card` / `create wallet`), showed in grafana
- Log already standardized (`context_id`, `message`, `stack_trace`, etc), can add ingress / egress log if using microservice / 3rd party service. 
  Log then scraped from stdout by fluentd using regex to elasticsearch.
- Log in elasticsearch can put on graylog / kibana / even add in grafana for monitoring.
- From kibana / graylog can add alerting when reaches certain threshold, sent to telegram / slack channel /etc
