# Simple GRPC/REST API as a Knative Service.

This repository contains a simple GRPC/REST API, which can be run as a Knative
service.

## Local

To run the service locally:

```shell
go run ./cmd/server
```

To hit the server locally:

```shell
# Create
curl -X POST -d '{"field1":"foo","field2":"bar"}' http://localhost:8080/v1/things

# List
curl http://localhost:8080/v1/things
```

## Knative

To install Knative, my goto is [`mink`](https://github.com/mattmoor/mink) (replace `mattmoor.dev` with your own domain):

```shell
mink install --domain mattmoor.dev
```

> _Note:_ `mink` prints the DNS record you need to configure for the specified domain when it is done.


To deploy this to Knative Serving:

```shell
ko apply -Bf config/
```


Now you should be able to curl things:

```shell
curl https://foo.default.mattmoor.dev/v1/things
```

Or load test things, my goto is [`vegeta`](https://github.com/tsenart/vegeta) (edit `attack.log` with your domain):

```shell
vegeta -cpus=1 attack -duration=4m -rate=1000/1s -targets=attack.log | vegeta report -type='hist[0,10ms,100ms,1s,10s]'
```

This should generate a nice little report with latency breakdown (it is most interesting if you let things scale to zero first!):

```
Bucket           #       %       Histogram
[0s,     10ms]   163     0.07%   
[10ms,   100ms]  233302  97.21%  ########################################################################
[100ms,  1s]     4487    1.87%   #
[1s,     10s]    2048    0.85%   
[10s,    +Inf]   0       0.00%   
```
