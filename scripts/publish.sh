#!/bin/bash
cat ./example.json | nats pub ORDERS.test
