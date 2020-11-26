# Logstash hook for logrus <img src="http://i.imgur.com/hTeVwmJ.png" width="40" height="40" alt=":walrus:" class="emoji" title=":walrus:" />
This repository is a fork of https://github.com/bshuster-repo/logrus-logstash-hook, with opinionated changes related to logrus/logstash.

Basically, you might be interested in this fork if you're likely to use logrus with the Report Caller, and if you want a ns precision on the logs sent to logstash.

All the credit of the work is given to Boaz Shuster.

Use this hook to send the logs to [Logstash](https://www.elastic.co/products/logstash).

# Usage

```go
package main

import (
        "github.com/bshuster-repo/logrus-logstash-hook"
        "github.com/sirupsen/logrus"
        "net"
)

func main() {
        log := logrus.New()
        conn, err := net.Dial("tcp", "logstash.mycompany.net:8911")
        if err != nil {
                log.Fatal(err)
        }
        hook := logrustash.New(conn, logrustash.DefaultFormatter(logrus.Fields{"type": "myappName"}))

        log.Hooks.Add(hook)
        ctx := log.WithFields(logrus.Fields{
                "method": "main",
        })
        ctx.Info("Hello World!")
}

```

This is how it will look like:

```ruby
{
    "@timestamp" => "2016-02-29T16:57:23.456Z",
      "@version" => "1",
         "level" => "info",
       "message" => "Hello World!",
        "method" => "main",
          "host" => "172.17.0.1",
          "port" => 45199,
          "type" => "myappName"
}
```

This is how it will look like with report caller enabled:

```ruby
{
    "@timestamp" => "2016-02-29T16:57:23.456Z",
      "@version" => "1",
         "level" => "info",
       "message" => "Hello World!",
        "method" => "main",
          "host" => "172.17.0.1",
          "port" => 45199,
          "type" => "myappName"
          "file" => "github.com/user/project/main.go:10"
          "func" => "github.com/user/project/moduleMain.Func
}
```

# Maintainers

Name         | Github    | Twitter    |
------------ | --------- | ---------- |
Boaz Shuster | ripcurld0 | @ripcurld0 |

# License

Boaz Shuster/MIT.
