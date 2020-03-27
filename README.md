# CEN API Go Server

This is the Go implementation of the [CEN API](https://docs.google.com/document/d/1f65V3PI214-uYfZLUZtm55kdVwoazIMqGJrxcYNI4eg/edit) for _Privacy-Preserving Distributed Contact Tracing_.

Tech Lead: Sourabh Niyogi (`sourabh@wolk.com`)

The flow is as follows:`
1. Apps use secret 128-bit AES keys (CENKeys) to broadcast CENs using BLE

2. CEN Apps records neighbor CENs in close proximity to the user

3. Users submit symptom / infection reports and reveal the secret CENKey in their application to a CEN API endpoint, resulting in a POST to `/cenreport` endpoint with `Report`  

4. Apps poll every N mins to `/cenkeys` for CEN Keys posts and matches them to CENs the user has observed.

5. Upon match the application retrieves the report

See also:
* Design documents:
* CEN API Documentation: (to be updated)
* CEN API Endpoint: (under construction) https://cenapi.wolk.com:8080   
* Rust Implementation:

### POST `/cenreport`
Request Body: JSON Object - `CENReport`
```
{"reportID":"1","report":"c2V2ZXJlIGZldmVyLGNvdWdoaW5nLGhhcmQgdG8gYnJlYXRoZQ==","cenKeys":"2cb87ba2f39a3119e4096cc6e04e68a8,4bb5c242916d923fe3565bd5c6b09dd3"}
```

### GET `/cenkeys/<timestamp>`

Returns CENKeys from `<timestamp>` to `<timestamp-3600>` (Under consideration)

Sample:
```
["67f48de38f35c231e34e533649ecbfeb","b85c4b373adde4c66651ba63aef40f48"]
```

### POST `/cenreport/<cenkey>`

Returns reports associated with CENKey.

Sample:
```
[{"reportID":"6d68785abd6f99ea646fd7c775a14a36f7cec7635a76d3b5554908c6c3a09af5","report":"c2V2ZXJlIGZldmVyLGNvdWdoaW5nLGhhcmQgdG8gYnJlYXRoZQ==","reportTimeStamp":1585326048}]
```


### MySQL Setup

1. After setting up your Mysql instance and updating the Default Connection strings in server.go, create tables with `cen.sql`
```

2. Check that you can (a) POST a CEN Report; (b) GET CEN Keys; (c) GET CEN Reports  with `go test -run TestBackendSimple`
```
[root@d5 backend]# go test -run TestBackendSimple
CENReportJSON Sample: {"reportID":"1","report":"c2V2ZXJlIGZldmVyLGNvdWdoaW5nLGhhcmQgdG8gYnJlYXRoZQ==","cenKeys":"79c256d12f5cad1af81e3658c0416e4e,4c29aa2b332199c82c825887ba179e27"}
CENKey 0: 79c256d12f5cad1af81e3658c0416e4e
CENKey 1: 4c29aa2b332199c82c825887ba179e27
ProcessGetCENKeys 18 records
Recent Keys: [67f48de38f35c231e34e533649ecbfeb]
Recent Keys: [b85c4b373adde4c66651ba63aef40f48]
Recent Keys: [046316753d2e5684d542a0a53944c6f7]
Recent Keys: [41371323cc938a0e3c55b0694bfd23f5]
Recent Keys: [8a54afbc79912d56d30d357ca6d86f06]
Recent Keys: [a8c15b2e81c7aa185618698b9bc61eb8]
Recent Keys: [001cb7c122064b79ebd2e878889c17c0]
Recent Keys: [fc87cbdf6e3638cb9ddffbe77cddff07]
Recent Keys: [5c52f9a5e539d69b179c5badb4aed724]
Recent Keys: [6413987722a3590faf38ec84461e5d2f]
Recent Keys: [0c0b9ddff6ab2d22519a793d941bdeb5]
Recent Keys: [befa6a1c065c9162890edd614ec33fcd]
Recent Keys: [4cfc2187fa3c6e5c9cafaeb03d09298c]
Recent Keys: [a1c7d918f2f7d464b493a6ed3e200bfd]
Recent Keys: [2cb87ba2f39a3119e4096cc6e04e68a8]
Recent Keys: [4bb5c242916d923fe3565bd5c6b09dd3]
Recent Keys: [4c29aa2b332199c82c825887ba179e27]
Recent Keys: [79c256d12f5cad1af81e3658c0416e4e]
ProcessGetCENReport SUCCESS (67f48de38f35c231e34e533649ecbfeb): [severe fever,coughing,hard to breathe]
ProcessGetCENReport SUCCESS (b85c4b373adde4c66651ba63aef40f48): [severe fever,coughing,hard to breathe]
PASS
ok	github.com/Co-Epi/coepi-backend-go/backend	0.032s
```

## Build + Run
```
$ make cen
go build -o bin/cen
Done building cen.  Run "bin/cen" to launch CEN Server.
```

## Test

After getting your SSL Certs in the right spot with a DNS entry that matches and running `bin/cen`, you can run this test:
```
[root@d5 coepi-backend-go]# go test -run TestCEN
EndpointCENReport[OK]
EndpointCENKeys: ["67f48de38f35c231e34e533649ecbfeb","b85c4b373adde4c66651ba63aef40f48","046316753d2e5684d542a0a53944c6f7","41371323cc938a0e3c55b0694bfd23f5","8a54afbc79912d56d30d357ca6d86f06","a8c15b2e81c7aa185618698b9bc61eb8","001cb7c122064b79ebd2e878889c17c0","fc87cbdf6e3638cb9ddffbe77cddff07","5c52f9a5e539d69b179c5badb4aed724","6413987722a3590faf38ec84461e5d2f","0c0b9ddff6ab2d22519a793d941bdeb5","befa6a1c065c9162890edd614ec33fcd","4cfc2187fa3c6e5c9cafaeb03d09298c","a1c7d918f2f7d464b493a6ed3e200bfd","2cb87ba2f39a3119e4096cc6e04e68a8","4bb5c242916d923fe3565bd5c6b09dd3","4c29aa2b332199c82c825887ba179e27","79c256d12f5cad1af81e3658c0416e4e","0e987952d86e19197a2be845d4423ae7","e98a09d7c3e9f7a9e43906f2eb8635d7"]
EndpointCENKeys SUCCESS: [severe fever,coughing,hard to breathe]
EndpointCENKeys SUCCESS: [severe fever,coughing,hard to breathe]
PASS
ok	github.com/Co-Epi/coepi-backend-go	0.124s
```

which does the same things as the above backend test except going through the HTTP Server.

### How it works (at a glance)

Right now the schema is pretty simple.  

1. `CENKeys` holds the mapping between revealed keys (cenKey) and reports (reportID), with `reportID` being the join key between the 2 tables.

```
mysql> select * from CENKeys;
+----------------------------------+------------------------------------------------------------------+------------+
| cenKey                           | reportID                                                         | reportTS   |
+----------------------------------+------------------------------------------------------------------+------------+
| 67f48de38f35c231e34e533649ecbfeb | e6f6cc772d11270cbac9b6e6cb24c451dc54b54172547f9ce327e9b5c5961609 | 1585326035 |
| b85c4b373adde4c66651ba63aef40f48 | e6f6cc772d11270cbac9b6e6cb24c451dc54b54172547f9ce327e9b5c5961609 | 1585326035 |
| 046316753d2e5684d542a0a53944c6f7 | 601e2a74dd04f3ec8b29012cefb92975eb31b60015584c803a4c27746537c23e | 1585326038 |
| 41371323cc938a0e3c55b0694bfd23f5 | 601e2a74dd04f3ec8b29012cefb92975eb31b60015584c803a4c27746537c23e | 1585326038 |
....
| 2cb87ba2f39a3119e4096cc6e04e68a8 | ef932667e953f400ee0ea3b6bf70e463f6adf57f124673db71b34d9ca93b91e4 | 1585326239 |
| 4bb5c242916d923fe3565bd5c6b09dd3 | ef932667e953f400ee0ea3b6bf70e463f6adf57f124673db71b34d9ca93b91e4 | 1585326239 |
| 4c29aa2b332199c82c825887ba179e27 | 68883a2b8ea382397114376c9c099fdde6a18e4f9c9b35f1177a9e0010a7818f | 1585326764 |
| 79c256d12f5cad1af81e3658c0416e4e | 68883a2b8ea382397114376c9c099fdde6a18e4f9c9b35f1177a9e0010a7818f | 1585326764 |
| 0e987952d86e19197a2be845d4423ae7 | 80d96e77ba86ff1766e8b0e6086f3b9266e9edc5e5a87cc87cc222331ecc6de8 | 1585326891 |
| e98a09d7c3e9f7a9e43906f2eb8635d7 | 80d96e77ba86ff1766e8b0e6086f3b9266e9edc5e5a87cc87cc222331ecc6de8 | 1585326891 |
+----------------------------------+------------------------------------------------------------------+------------+
```

2. `CENReport` holds the reports

```
mysql> select * from CENReport;
+------------------------------------------------------------------+---------------------------------------+----------------+------------+
| reportID                                                         | report                                | reportMimeType | reportTS   |
+------------------------------------------------------------------+---------------------------------------+----------------+------------+
| 0075b901445a732a5d6f299cbf9d14692b12f2724b9443d7b3740a62c7a81ab6 | severe fever,coughing,hard to breathe |                | 1585326226 |
| 374ecfb7b0c437a3054e1b75f546faa8c07f81809a9a77ce5752bed1b1218ab7 | severe fever,coughing,hard to breathe |                | 1585325928 |
| 601e2a74dd04f3ec8b29012cefb92975eb31b60015584c803a4c27746537c23e | severe fever,coughing,hard to breathe |                | 1585326038 |
| 68883a2b8ea382397114376c9c099fdde6a18e4f9c9b35f1177a9e0010a7818f | severe fever,coughing,hard to breathe |                | 1585326764 |
| 692ad09bf3adb0467f7b19ba4f28076a7de62b9235eb7495da0b57036162a586 | severe fever,coughing,hard to breathe |                | 1585326043 |
| 6d68785abd6f99ea646fd7c775a14a36f7cec7635a76d3b5554908c6c3a09af5 | severe fever,coughing,hard to breathe |                | 1585326048 |
| 7f9dcb9b89870384c91f19c7e07bce1b0bc6b63f9f331f64287d81367e7a82e8 | severe fever,coughing,hard to breathe |                | 1585326041 |
| 80d96e77ba86ff1766e8b0e6086f3b9266e9edc5e5a87cc87cc222331ecc6de8 | severe fever,coughing,hard to breathe |                | 1585326891 |
| c41c3715d23127a75d93e0275b476badace16193d5cf07eac7b256db18c37a00 | severe fever,coughing,hard to breathe |                | 1585326046 |
| e6f6cc772d11270cbac9b6e6cb24c451dc54b54172547f9ce327e9b5c5961609 | severe fever,coughing,hard to breathe |                | 1585326035 |
| ef932667e953f400ee0ea3b6bf70e463f6adf57f124673db71b34d9ca93b91e4 | severe fever,coughing,hard to breathe |                | 1585326239 |
+------------------------------------------------------------------+---------------------------------------+----------------+------------+
```

Implementation:
 - The `/cenreport` POST writes the raw CENReport to the `cenkeys` table (1-3 rows, with 1 row per week) and `CENReport` table (1 new row).
 - The `/cenkeys` GET endpoint reads just the `CENKeys` table, requiring the index.  
 - The `/cenreport` GET endpoint reads the join between the 2 tables with the CENKey.

It is expected that CDNs can replace this, with mobile applications 

Importantly, no PII data is held in any table.

## Contributions to Go Server

Great!  Join the CoEpi Slack channel and say hi to @sourabh, or email sourabh@wolk.com.

* If you are an experienced Go person (channels, mutexes / concurrency) please consider extending this server to support infection reports.

* If you are an experienced iOS or Android person with a strong interest in Bluetooth Low Energy and want to do client-server bridge work, check out [app-ios](https://github.com/Co-Epi/app-ios) or [app-android](https://github.com/Co-Epi/app-android).
