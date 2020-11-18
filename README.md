# agentserver
S3 transfer(download) in multi-nodes. It's a fast way to download file from S3 to local.

## Single mode 
### How to use
```
Usage:
  agentserver.exe [OPTIONS]

Application Options:
      /cmd:          Command [agent, master, execute] (default: execute)
      /dir:          work directory
      /s3:           s3 url
  /H, /host:         master host or ip
  /p, /port:         listen port
      /accesskey:    AWS access key id.
      /secretkey:    AWS sercret key
      /region:       AWS region (default: ap-northeast-1)
      /token:        aws credentials token
      /proxy:        Proxy server
      /retry:        Retry count (default: 3)
      /jobid:        Job ID when cmd=worker
      /log:          Log level (trace,debug,info,warn,error,fatal)
  /d, /distributed   work with execute.
      /dryrun        debug mode
  /e, /env:          Stage [deve]
  /a, /agent:        agent server [host:port]
  /m, /master:       master server [host:port]
  /s, /storage:      storage [host:port]
      /wait          wait for download finish
  /v, /version       Show version

Help Options:
  /?                 Show this help message
  /h, /help          Show this help message
```

**download s3://dittos/days/09/21 to c:\tmp\data.**<br>
agentserver.exe --dir c:\tmp\data --s3 s3://dittos/days/09/21 <br>


## Multi-Nodes mode
![mode](https://github.com/cereskou/agentserver/blob/main/images/layout.png)

### Execute in master mode
```
agentserver.exe --cmd master --log trace

time="2020-11-18T18:41:46+09:00" level=info msg="Run Master"
time="2020-11-18T18:41:46+09:00" level=debug msg="http server started on 0.0.0.0:9090"
time="2020-11-18T18:41:50+09:00" level=trace msg="Received {\"mode\":\"agent\",\"ip\":\"192.168.50.113\",\"port\":9091,\"timestam\":\"2020-11-18T18:41:50.5056251+09:00\",\"jobid\":0,\"status\":0} from 192.168.50.113:56077\n"
time="2020-11-18T18:41:51+09:00" level=trace msg="Check alive at 2020-11-18 18:41:51.95758 +0900 JST m=+5.244475001"
time="2020-11-18T18:41:51+09:00" level=trace msg="Check 192.168.50.113:9091 ..."
time="2020-11-18T18:41:51+09:00" level=trace msg="Received {\"mode\":\"master\",\"ip\":\"192.168.50.113\",\"port\":9090,\"timestam\":\"2020-11-18T18:41:51.961565+09:00\",\"jobid\":0,\"status\":0} from 192.168.50.113:56078\n"
time="2020-11-18T18:41:51+09:00" level=trace msg="192.168.50.113:9091 is alive"
```

### Execute in agent mode
Multiple agents can be started on others machine.And the agent will registerate to the master by using broadcast.

```
agentserver.exe --cmd agent --log trace

time="2020-11-18T18:41:35+09:00" level=info msg="Run Agent"
time="2020-11-18T18:41:35+09:00" level=debug msg="http server started on 0.0.0.0:9091"
{"time":"2020-11-18T18:41:35.7345338+09:00","id":"","remote_ip":"192.168.50.113","host":"192.168.50.113:9091","method":"GET","uri":"/health","user_agent":"s3transfer distributed","status":200,"error":"","latency":0,"latency_human":"0s","bytes_in":0,"bytes_out":11}
{"time":"2020-11-18T18:41:39.7338912+09:00","id":"","remote_ip":"192.168.50.113","host":"192.168.50.113:9091","method":"GET","uri":"/health","user_agent":"s3transfer distributed","status":200,"error":"","latency":998800,"latency_human":"998.8µs","bytes_in":0,"bytes_out":11}
```

### Download
--distributed will start the download process on agent node.<br>
the data will saved to master pc's c:\tmp\data if you are not specified the target by using --storage. <br>
--storage can be a agent or master.

```
agentserver.exe --dir c:\tmp\data --s3 s3://dittos/days/09/21 --distributed --master 192.168.50.113
```

