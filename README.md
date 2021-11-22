# Payment Processing

###Assumption
- DB Sharding is out of scope
- Redis cluster, master slave, sentinel mode, RDB, AOF are out of scope.
- Auth service can be separate on different microservice.
- Can use multiple instance of payment-processing service, using loadbalancer in front
- Rate limit is implemented on service now, can use nginx/kong/istio ingress + envoy,etc for rate limiter if using multiple instances.
- Rate limit using simple rate limit for current QPS. can use more complex rate limit algorithm like leaky bucket, sliding window, etc for improvement.

###Directory Structure
I am using directory structure guideline from [golang standard](https://github.com/golang-standards/project-layout)
- `cmd` as runmode (webservice, cron-mode, consumer-mode, etc)
- `config` for application config file
- `internal` for private application code, consist of:
    1. `app`, this is for api controller and handler (manager in java)
    2. `pkg`, this is for application shared component and DAO, request response struct also included here in `dto` dir
- `migration` for history of migration / schema alter. note that all schema alter copied to `init.sql` used by `docker-compose`
- `pkg` for shared  / common library
- `vendor` for dependencies
- `apis` for api documentation
