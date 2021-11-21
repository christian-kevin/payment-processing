# Payment Processing

###Assumption
- DB Sharding is out of scope
- Redis cluster, master slave, sentinel mode, RDB, AOF are out of scope.
- Auth service can be separate on different microservice.
- Can use multiple instance of payment-processing service, using loadbalancer in front
- Rate limit is implemented on service now, can use nginx/kong/istio ingress + envoy,etc for rate limiter if using multiple instances.
- Rate limit using simple rate limit for current QPS. can use more complex rate limit algorithm like leaky bucket, sliding window, etc for improvement.