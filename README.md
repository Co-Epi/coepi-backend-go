# CoEpi Go Server

The client-server flow is as follows:`
1. Users record symptoms in their CoEpi app, resulting in a POST to `/exposureandsymptoms` endpoint with their `Symptoms` and UUIDs/Dates of `Contacts`
2. All CoEpi apps poll every N mins with POST to `/exposurecheck` to see if the server has received any symptoms of the above matching those that the device has seen.

## CoEpi API Endpoints

Test Endpoint: https://coepi.wolk.com:8081

### POST `/exposureandsymptoms`
Body: JSON Object - `ExposureAndSymptoms`
```
{"Symptoms":"SlNPTkJMT0I6c2V2ZXJlIGZldmVyLGNvdWdoaW5n","Contacts":[{"UUID":"ax","Date":"2020-03-04"},{"UUID":"by","Date":"2020-03-15"},{"UUID":"cz","Date":"2020-03-20"}]}
```

Response:
```
200 OK
400 Error
```

Behavior:
* Server records ContactIDs (BLE Proximity Event UUIDs) in `contacts` table with { Date, SymptomHash }
* Server records Symptoms in `symptoms` table keyed by SymptomHash

### POST `/exposurecheck`
Body: JSON Object - `ExposureCheck`
Sample:
```
{"Contacts":[{"UUID":"by","Date":"2020-03-04"}]}
```
Behavior:
* Server gets UUIDs (BLE Proximity Event UUIDs) and checks `contacts` table for potential { Date, SymptomHash } combinations
* With hits, Server fetches Symptoms from `symptoms` table keyed by SymptomHash and returns array of byte blobs


### BigTable Setup

Use `cbt` https://cloud.google.com/bigtable/docs/quickstart-cbt

1. After setting up your BigTable `co-epi` instance and updating the project/instance strings in server.go, create 2 tables and their families with `cbt`
```
cbt createtable contacts
cbt createtable symptoms
cbt createfamily symptoms case
cbt createfamily contacts symptoms
```
Check with:
```
# cbt ls
contacts
symptoms
# cbt ls symptoms
Family Name	GC Policy
case		<never>
# cbt ls contacts
Family Name	GC Policy
symptoms	<never>
```

2. Check that you can write some data with `go test -run TestBackendSimple`
```
# go test -run TestBackendSimple
processExposureCheck(check1) SUCCESS: [JSONBLOB:severe fever,coughing]
processExposureCheck(check0) SUCCESS: []
PASS
ok	github.com/wolkdb/coepi-backend-go/server	0.412s
```

`TestBackendSimple` goes through a Backend `exposureandsymptoms` and `exposure check`

## Build + Run

```
# make coepi
go build -o bin/coepi
Done building coepi.  Run "/bin/coepi" to launch coepi.
# bin/coepi
...
```

## Test

After getting your SSL Certs in the right spot with a DNS entry that matches and running `bin/coepi`, you can run this test:
```
# go test -run TestCoepi
ExposureAndSymptoms Sample: {"Symptoms":"SlNPTkJMT0I6c2V2ZXJlIGZldmVyLGNvdWdoaW5n","Contacts":[{"UUID":"ax","Date":"2020-03-04"},{"UUID":"by","Date":"2020-03-15"},{"UUID":"cz","Date":"2020-03-20"}]}
exposureandsymptoms[OK]exposurecheck(check1) SUCCESS: [JSONBLOB:severe fever,coughing]
exposurecheck(check0) SUCCESS: []
PASS
ok	github.com/wolkdb/coepi-backend-go	0.589s
```

`TestCoepiSimple` does the same thing as `TestBackendSimple` except going through the HTTP Server.

### How it works (at a glance)

In the `contacts` table, there is a map between each UUID and a `symptomHash`

```
# cbt read contacts
2020/03/21 02:36:20 -creds flag unset, will use gcloud credential
----------------------------------------
ax
  symptoms:2020-03-04                      @ 2020/03/21-02:34:55.307000
    "b93a90a843ed293522aa803781298dac436040fa231a189e52c6994a5d591f09"
----------------------------------------
by
  symptoms:2020-03-15                      @ 2020/03/21-02:34:55.389000
    "b93a90a843ed293522aa803781298dac436040fa231a189e52c6994a5d591f09"
----------------------------------------
cz
  symptoms:2020-03-20                      @ 2020/03/21-02:34:55.443000
    "b93a90a843ed293522aa803781298dac436040fa231a189e52c6994a5d591f09"
```

In the `symptoms` table, there is map between `symptomHash`es held in `contact` and a blob of bytes.
```
# cbt read symptoms
----------------------------------------
b93a90a843ed293522aa803781298dac436040fa231a189e52c6994a5d591f09
  case:symptoms                            @ 2020/03/21-02:03:55.604000
    "JSONBLOB:severe fever,coughing"
```

## Sourabh's Brass Tacks Questions

0. I'm going to use the word "contact" in the BLE sense rather than the address book sense, but I saw Jack using the word "interaction" and not use UUID, maybe can we get our word usages straightened out?

1. The server function in `Simplified API` is to do blind contact matching and reveal `Symptoms` only on matched contact, where the server can consider `Symptoms` to be just a blob of bytes... **except what**?
Guess: The server could usefully parse out { `QuizTime`, `DiseaseID` } to filter "old" contacts in the matching based on (a)+(b).  If nothing is disease specific, then we can treat all diseases with logic like "discard all interactions greater than 14 days" and there is just a free parameter "14" on the server.  Then `Symptoms` is just `[]byte` in `ExposureAndSymptoms` coming in `[][]byte` in `ExposureCheck` response.

2. I still don't understand what the "first 64 bits/full 128 bits" is all about yet.  
If the goal is _information hiding_, I would think that when A and B don't want to reveal that A+B interacted to the server, since they both know id(A) and id(B) (where id(A) is the non-persistent BLE identifier that B's device can see), they can both compute:
```
(1) H(id(A)||id(B))     OR
(2) H(id(B)||id(A))
```
and choose the smaller one in `H` space or `id` space (pick one).  When A posts N observations of `ExposureAndSymptoms` to a CoEpi server, the server knows nothing about id(A), id(B), or that A and B interactions, with H being SHA256 which is easy to do in mobile.

 If the goal is to reduce bandwidth (because 32 bytes is overkill), you can just take the first 8 bytes of H(...), and do a special purpose wire format like:
```
CONCAT(H(...)(16 chars)|geohash (5-6 chars)|daysago(2 chars))
```
mashing every pair of interactions into like 24 chars, that the server splits every 24 bytes without any of JSON's overhead.  Then 10K interactions over a 30 day period will be 240K, and with redundancy in the `geohash` and `daysago` you probably will see greater than 50% compression and we get straight into testing gzip compression.  Then the `/exposurecheck` API can just return matched an array of matched `Symptoms` directly with the client supplied H(...).  I don't think there needs to a another endpoint, and if we did, the return value would be an Array of `SymptomHash`.

3. What is the role of the (optional) `GeoHash` and `Datestamp`?  Is the `GeoHash` or `Datestamp` supposed to be used for contact matching at all?   They aren't doing any work yet (except to filter out old interactions, maybe).  So, to make them do some work and store them in the backend, we should have a `/symptoms` API endpoint to filter on these attributes.  We can organize the backend tables around that and make have someone who knows Google Maps showing symptoms in a geohash and how exposurecheck works visually in simulations at least.  

4. Based on the above I think we should nail down sample JSON POSTs and get parity between `coepi-backend-go` and `coepi-backend-rust` (and others) immediately.  We can use this to name our internal tables/attributes identically
across `coepi-backend-*`, but really I want to get all the capitalization/underscore and names mapped to the concepts done today.

5. After we sort out the above, it would be good to talk through the Cache expires header idea and how it relates to the "pull X times a day" design pattern exactly.   
