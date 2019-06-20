# Messaging Broker based on NATS

## Motivation

We want to build a special messenger for distributed system.

* Support consistent hash for queue message dispatching. This is strongly needed by our stateful service.
* Support some kind of retry and some kind of message storage. Use ack/retry to make a kind of reliable message routing.
* Support service discovery. There are no stable service discovery on k8s. The kub-master crushed for a day sometimes. ZK k8s deployment often broken, nats streaming server not stable also. Looks like RAFT/Gossip works not good at k8s to synchronize data (Don't support consistent file storage at our k8s). We want to build our own service discovery based on database as we have a stable database to use.

### Stateful micro-service

There are a lot of case we need stateful micro-service. Cache should be distribute to different PODs for memory, storage limitation or state consistent requirement. The traditional solution is service discovery + consistent hashing. There is another solution is use message queue and routing messages to subscribers by consistent hash. Now only RabbitMQ support consistent hash by a plugin. For heavy traffic would based on service discovery.

## Infrastructure

The 