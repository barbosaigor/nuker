# WIP: Nuker

**WIP (Work in Progress)**  
_This is not the final product_    

Nuker is a CLI tool for load testing, with a powerful configuration file (but easy) for planning your tests.   
It's a suitable alternative for [JMeter](https://jmeter.apache.org/) like tools.  

## Install  
```sh
$ go get -u github.com/barbosaigor/nuker/...  
```

**How to use**  
```sh
$ nuker plan.yaml  
```

## (Goals) Features  
* High throughput   
* Easy to write configuration file with an expressive config plan   
* Observability - Metrics  
* Detailed request logs (nuker will create a log file by default)  
* Easy to connect with another node, replicate plan for each node or distribute work (should be configurable which algorithm)  
