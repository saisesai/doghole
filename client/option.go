package main

import "flag"

var optionName = flag.String("name", "default", "service name")
var optionLocalAddr = flag.String("local-addr", "localhost:3001", "client listen address")
var optionServerAddr = flag.String("server-addr", "localhost:3000", "server listen address")
var optionServiceAddr = flag.String("service-addr", "localhost:80", "service listen addr")
