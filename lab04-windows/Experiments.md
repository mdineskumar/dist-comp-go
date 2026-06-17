```bash
docker exec lab04-client /lab04/client/client_bin rpc put hello world
docker exec lab04-client /lab04/client/client_bin grpc get hello
docker exec lab04-client /lab04/client/client_bin rest delete hello
docker exec lab04-client /lab04/client/client_bin rest list

# Test same operation with ALL 3 protocols at once
docker exec lab04-client /lab04/client/client_bin all put city London
docker exec lab04-client /lab04/client/client_bin all get city

# Benchmark — compare performance of all 3
docker exec lab04-client /lab04/client/client_bin bench 100

# Test REST directly with curl (from your LAPTOP terminal, not container)
curl -X PUT http://localhost:8080/keys/city \
 -H 'Content-Type: application/json' \
 -d '{"value":"London"}'
curl http://localhost:8080/keys/city
curl http://localhost:8080/keys
curl -X DELETE http://localhost:8080/keys/city


# RPC
# Test
docker exec lab04-client /lab04/client/client_bin rpc put city London
# Expected: [net/rpc] Stored key="city" value="London"
docker exec lab04-client /lab04/client/client_bin rpc get city
# Expected: [net/rpc] key="city" value="London"
docker exec lab04-client /lab04/client/client_bin rpc list
# Expected: [net/rpc] 1 keys: [city]

#gRPC tasks
docker exec lab04-client /lab04/client/client_bin grpc put city London
# Expected: [gRPC] Stored key="city" value="London"
docker exec lab04-client /lab04/client/client_bin grpc get city
# Expected: [gRPC] key="city" value="London"
# Cross-protocol test — put via net/rpc, get via gRPC
docker exec lab04-client /lab04/client/client_bin rpc put country UK
docker exec lab04-client /lab04/client/client_bin grpc get country
# Expected: key="country" value="UK"
# Both reach the same underlying store!


# REST 
docker exec lab04-client /lab04/client/client_bin rest put city London
docker exec lab04-client /lab04/client/client_bin rest get city
docker exec lab04-client /lab04/client/client_bin rest list

# BENCH MARK
docker exec lab04-client /lab04/client/client_bin bench 100
# Expected output:
# Protocol Total time Avg/op Ops/sec
# net/rpc 45ms 225us 4444
# gRPC 38ms 190us 5263
# REST 210ms 1050us 952
# (your numbers will differ — what matters is the relative order)

# Put via one protocol, get via another
docker exec lab04-client /lab04/client/client_bin rpc put benchmark test123
docker exec lab04-client /lab04/client/client_bin grpc get benchmark
docker exec lab04-client /lab04/client/client_bin rest get benchmark
# All three should return: test123
```


PS C:\Users\DM2719\Desktop\dist-comp-go\lab04-windows> docker exec lab04-client /lab04/client/client_bin bench 100

══════════════════════════════════════════════
 Lab 04 — Benchmark: 100 operations per protocol
══════════════════════════════════════════════

[net/rpc] Connected to lab04-server:7000
[gRPC] Connected to lab04-server:7001
[REST] PUT key="bench-key-0" value="value"
[REST] GET key=bench-key-0
[REST] PUT key="bench-key-1" value="value"
[REST] GET key=bench-key-1
[REST] PUT key="bench-key-2" value="value"
[REST] GET key=bench-key-2
[REST] PUT key="bench-key-3" value="value"
[REST] GET key=bench-key-3
[REST] PUT key="bench-key-4" value="value"
[REST] GET key=bench-key-4
[REST] PUT key="bench-key-5" value="value"
[REST] GET key=bench-key-5
[REST] PUT key="bench-key-6" value="value"
[REST] GET key=bench-key-6
[REST] PUT key="bench-key-7" value="value"
[REST] GET key=bench-key-7
[REST] PUT key="bench-key-8" value="value"
[REST] GET key=bench-key-8
[REST] PUT key="bench-key-9" value="value"
[REST] GET key=bench-key-9
[REST] PUT key="bench-key-10" value="value"
[REST] GET key=bench-key-10
[REST] PUT key="bench-key-11" value="value"
[REST] GET key=bench-key-11
[REST] PUT key="bench-key-12" value="value"
[REST] GET key=bench-key-12
[REST] PUT key="bench-key-13" value="value"
[REST] GET key=bench-key-13
[REST] PUT key="bench-key-14" value="value"
[REST] GET key=bench-key-14
[REST] PUT key="bench-key-15" value="value"
[REST] GET key=bench-key-15
[REST] PUT key="bench-key-16" value="value"
[REST] GET key=bench-key-16
[REST] PUT key="bench-key-17" value="value"
[REST] GET key=bench-key-17
[REST] PUT key="bench-key-18" value="value"
[REST] GET key=bench-key-18
[REST] PUT key="bench-key-19" value="value"
[REST] GET key=bench-key-19
[REST] PUT key="bench-key-20" value="value"
[REST] GET key=bench-key-20
[REST] PUT key="bench-key-21" value="value"
[REST] GET key=bench-key-21
[REST] PUT key="bench-key-22" value="value"
[REST] GET key=bench-key-22
[REST] PUT key="bench-key-23" value="value"
[REST] GET key=bench-key-23
[REST] PUT key="bench-key-24" value="value"
[REST] GET key=bench-key-24
[REST] PUT key="bench-key-25" value="value"
[REST] GET key=bench-key-25
[REST] PUT key="bench-key-26" value="value"
[REST] GET key=bench-key-26
[REST] PUT key="bench-key-27" value="value"
[REST] GET key=bench-key-27
[REST] PUT key="bench-key-28" value="value"
[REST] GET key=bench-key-28
[REST] PUT key="bench-key-29" value="value"
[REST] GET key=bench-key-29
[REST] PUT key="bench-key-30" value="value"
[REST] GET key=bench-key-30
[REST] PUT key="bench-key-31" value="value"
[REST] GET key=bench-key-31
[REST] PUT key="bench-key-32" value="value"
[REST] GET key=bench-key-32
[REST] PUT key="bench-key-33" value="value"
[REST] GET key=bench-key-33
[REST] PUT key="bench-key-34" value="value"
[REST] GET key=bench-key-34
[REST] PUT key="bench-key-35" value="value"
[REST] GET key=bench-key-35
[REST] PUT key="bench-key-36" value="value"
[REST] GET key=bench-key-36
[REST] PUT key="bench-key-37" value="value"
[REST] GET key=bench-key-37
[REST] PUT key="bench-key-38" value="value"
[REST] GET key=bench-key-38
[REST] PUT key="bench-key-39" value="value"
[REST] GET key=bench-key-39
[REST] PUT key="bench-key-40" value="value"
[REST] GET key=bench-key-40
[REST] PUT key="bench-key-41" value="value"
[REST] GET key=bench-key-41
[REST] PUT key="bench-key-42" value="value"
[REST] GET key=bench-key-42
[REST] PUT key="bench-key-43" value="value"
[REST] GET key=bench-key-43
[REST] PUT key="bench-key-44" value="value"
[REST] GET key=bench-key-44
[REST] PUT key="bench-key-45" value="value"
[REST] GET key=bench-key-45
[REST] PUT key="bench-key-46" value="value"
[REST] GET key=bench-key-46
[REST] PUT key="bench-key-47" value="value"
[REST] GET key=bench-key-47
[REST] PUT key="bench-key-48" value="value"
[REST] GET key=bench-key-48
[REST] PUT key="bench-key-49" value="value"
[REST] GET key=bench-key-49
[REST] PUT key="bench-key-50" value="value"
[REST] GET key=bench-key-50
[REST] PUT key="bench-key-51" value="value"
[REST] GET key=bench-key-51
[REST] PUT key="bench-key-52" value="value"
[REST] GET key=bench-key-52
[REST] PUT key="bench-key-53" value="value"
[REST] GET key=bench-key-53
[REST] PUT key="bench-key-54" value="value"
[REST] GET key=bench-key-54
[REST] PUT key="bench-key-55" value="value"
[REST] GET key=bench-key-55
[REST] PUT key="bench-key-56" value="value"
[REST] GET key=bench-key-56
[REST] PUT key="bench-key-57" value="value"
[REST] GET key=bench-key-57
[REST] PUT key="bench-key-58" value="value"
[REST] GET key=bench-key-58
[REST] PUT key="bench-key-59" value="value"
[REST] GET key=bench-key-59
[REST] PUT key="bench-key-60" value="value"
[REST] GET key=bench-key-60
[REST] PUT key="bench-key-61" value="value"
[REST] GET key=bench-key-61
[REST] PUT key="bench-key-62" value="value"
[REST] GET key=bench-key-62
[REST] PUT key="bench-key-63" value="value"
[REST] GET key=bench-key-63
[REST] PUT key="bench-key-64" value="value"
[REST] GET key=bench-key-64
[REST] PUT key="bench-key-65" value="value"
[REST] GET key=bench-key-65
[REST] PUT key="bench-key-66" value="value"
[REST] GET key=bench-key-66
[REST] PUT key="bench-key-67" value="value"
[REST] GET key=bench-key-67
[REST] PUT key="bench-key-68" value="value"
[REST] GET key=bench-key-68
[REST] PUT key="bench-key-69" value="value"
[REST] GET key=bench-key-69
[REST] PUT key="bench-key-70" value="value"
[REST] GET key=bench-key-70
[REST] PUT key="bench-key-71" value="value"
[REST] GET key=bench-key-71
[REST] PUT key="bench-key-72" value="value"
[REST] GET key=bench-key-72
[REST] PUT key="bench-key-73" value="value"
[REST] GET key=bench-key-73
[REST] PUT key="bench-key-74" value="value"
[REST] GET key=bench-key-74
[REST] PUT key="bench-key-75" value="value"
[REST] GET key=bench-key-75
[REST] PUT key="bench-key-76" value="value"
[REST] GET key=bench-key-76
[REST] PUT key="bench-key-77" value="value"
[REST] GET key=bench-key-77
[REST] PUT key="bench-key-78" value="value"
[REST] GET key=bench-key-78
[REST] PUT key="bench-key-79" value="value"
[REST] GET key=bench-key-79
[REST] PUT key="bench-key-80" value="value"
[REST] GET key=bench-key-80
[REST] PUT key="bench-key-81" value="value"
[REST] GET key=bench-key-81
[REST] PUT key="bench-key-82" value="value"
[REST] GET key=bench-key-82
[REST] PUT key="bench-key-83" value="value"
[REST] GET key=bench-key-83
[REST] PUT key="bench-key-84" value="value"
[REST] GET key=bench-key-84
[REST] PUT key="bench-key-85" value="value"
[REST] GET key=bench-key-85
[REST] PUT key="bench-key-86" value="value"
[REST] GET key=bench-key-86
[REST] PUT key="bench-key-87" value="value"
[REST] GET key=bench-key-87
[REST] PUT key="bench-key-88" value="value"
[REST] GET key=bench-key-88
[REST] PUT key="bench-key-89" value="value"
[REST] GET key=bench-key-89
[REST] PUT key="bench-key-90" value="value"
[REST] GET key=bench-key-90
[REST] PUT key="bench-key-91" value="value"
[REST] GET key=bench-key-91
[REST] PUT key="bench-key-92" value="value"
[REST] GET key=bench-key-92
[REST] PUT key="bench-key-93" value="value"
[REST] GET key=bench-key-93
[REST] PUT key="bench-key-94" value="value"
[REST] GET key=bench-key-94
[REST] PUT key="bench-key-95" value="value"
[REST] GET key=bench-key-95
[REST] PUT key="bench-key-96" value="value"
[REST] GET key=bench-key-96
[REST] PUT key="bench-key-97" value="value"
[REST] GET key=bench-key-97
[REST] PUT key="bench-key-98" value="value"
[REST] GET key=bench-key-98
[REST] PUT key="bench-key-99" value="value"
[REST] GET key=bench-key-99

── Results (100 put+get pairs) ─────────────────
  Protocol      Total time    Avg/op        Ops/sec
  ────────────────────────────────────────────────────
  net/rpc       57ms          286µs         3499
  gRPC          64ms          318µs         3145
  REST          158ms         792µs         1262

── Observations ─────────────────────────────
  Which protocol was fastest? Why?
  Can you call the REST server with curl? (Try it!)
  Can you call net/rpc with curl? Why not?


  DM2719@225L11-DM2719 MINGW64 ~
$ curl -X PUT http://localhost:8080/keys/city \
 -H 'Content-Type: application/json' \
 -d '{"value":"London"}'
{"success":true}

DM2719@225L11-DM2719 MINGW64 ~
$ curl http://localhost:8080/keys/city
{"found":true,"value":"London"}

DM2719@225L11-DM2719 MINGW64 ~
$ curl http://localhost:8080/keys/cityr
{"found":false}

DM2719@225L11-DM2719 MINGW64 ~
$cucurl http://localhost:8080/keys
bash: $'\302\226curl': command not found

DM2719@225L11-DM2719 MINGW64 ~
$ curl http://localhost:8080/keys
{"count":102,"keys":["bench-key-35","bench-key-66","bench-key-83","bench-key-30","bench-key-38","bench-key-62","bench-key-65","bench-key-50","bench-key-67","bench-key-16","bench-key-72","bench-key-78","bench-key-88","bench-key-13","bench-key-31","bench-key-46","bench-key-59","bench-key-41","bench-key-8","bench-key-18","bench-key-23","bench-key-26","bench-key-33","bench-key-7","test","bench-key-9","bench-key-14","bench-key-80","bench-key-95","bench-key-55","bench-key-4","bench-key-51","bench-key-90","bench-key-85","bench-key-99","bench-key-40","bench-key-71","bench-key-19","bench-key-15","bench-key-20","bench-key-73","bench-key-97","bench-key-12","bench-key-32","bench-key-86","bench-key-3","bench-key-47","bench-key-48","bench-key-22","bench-key-25","bench-key-81","bench-key-87","bench-key-91","bench-key-0","bench-key-21","bench-key-57","bench-key-77","bench-key-89","bench-key-64","bench-key-84","bench-key-98","bench-key-5","bench-key-6","bench-key-52","bench-key-53","bench-key-93","bench-key-1","bench-key-37","bench-key-39","bench-key-54","bench-key-70","bench-key-2","bench-key-10","bench-key-24","bench-key-44","bench-key-75","bench-key-17","bench-key-96","bench-key-43","bench-key-68","bench-key-34","bench-key-36","bench-key-74","bench-key-82","bench-key-94","city","bench-key-58","bench-key-60","bench-key-63","bench-key-69","bench-key-79","bench-key-92","bench-key-11","bench-key-27","bench-key-28","bench-key-42","bench-key-61","bench-key-29","bench-key-45","bench-key-49","bench-key-56","bench-key-76"]}

DM2719@225L11-DM2719 MINGW64 ~
$ curl -X DELETE http://localhost:8080/keys/city
{"deleted":true}
