# CoEpi Go Server

Tech Lead: Sourabh Niyogi (`sourabh@wolk.com`)

The client-server flow is as follows:`
1. Users record symptoms in their CoEpi app, resulting in a POST to `/exposureandsymptoms` endpoint with their `Symptoms` and UUIDs/Dates of `Contacts`
2. All CoEpi apps poll every N mins with POST to `/exposurecheck` to see if the server has received any symptoms of the above matching those that the device has seen.

* CoEpi API Documentation: https://documenter.getpostman.com/view/10811660/SzS8sQQY?version=latest
* Endpoint: (under construction) https://coepi.wolk.com:8080   
* CoEpi ExposureCheck Complexity Problem: https://docs.google.com/spreadsheets/d/1WnlshNGkPJOajQXCNCbLntfDrAsUCbS5ASEJ3D2nmGU/edit#gid=0

### POST `/exposureandsymptoms`
Request Body: JSON Object - `ExposureAndSymptoms`
```
{
	"symptoms": "SlNPTkJMT0I6c2V2ZXJlIGZldmVyLGNvdWdoaW5n",
	"contacts": [{
		"uuidHash": "ax",
		"dateStamp": "2020-03-04"
	}, {
		"uuidHash": "by",
		"dateStamp": "2020-03-15"
	}, {
		"uuidHash": "cz",
		"dateStamp": "2020-03-20"
	}]
}
```

Implementation Behavior:
* Server records ContactIDs (Hashes of pairs of BLE Proximity Event UUIDs) in `contacts` table with { Date, SymptomHash }
* Server records Symptoms in `symptoms` table keyed by SymptomHash

### POST `/exposurecheck`
Request Body: JSON Object - `ExposureCheck`
Sample:
```
{
	"contacts": [{
		"uuidHash": "by",
		"dateStamp": "2020-03-04"
	}]
}
```

Implementation Behavior:
* Server gets ContactIDs (Hashes of pairs of BLE Proximity Event UUIDs) and checks `contacts` table for potential { Date, SymptomHash } combinations
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

2. Check that you can write some contacts with `go test -run TestBackendSimple`
```
[root@d5 server]# go test -run TestBackendSimple
processExposureCheck(check1) SUCCESS: [JSONBLOB:severe fever,coughing]
processExposureCheck(check0) SUCCESS: []
PASS
ok	github.com/wolkdb/coepi-backend-go/server	0.412s
```

## Build + Run
```
$ make coepi
go build -o bin/coepi
Done building coepi.  Run "/bin/coepi" to launch coepi.
```

## Test

After getting your SSL Certs in the right spot with a DNS entry that matches and running `bin/coepi`, you can run this test:
```
# go test -run TestCoepi
exposureandsymptoms[OK]
exposurecheck(check1) SUCCESS: [JSONBLOB:severe fever,coughing]
exposurecheck(check0) SUCCESS: []
PASS
ok	github.com/wolkdb/coepi-backend-go	0.589s
```

which does the same things as the above backend test except going through the HTTP Server.

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

In the `symptoms` table, there is map between symptomHashes and a blob of bytes.
```
# cbt read symptoms
----------------------------------------
b93a90a843ed293522aa803781298dac436040fa231a189e52c6994a5d591f09
  case:symptoms                            @ 2020/03/21-02:03:55.604000
    "JSONBLOB:severe fever,coughing"
```


## Sourabh's Brass Tacks Questions

0. I'm going to use the word "contact" in the BLE sense rather than the address book sense, but I saw Jack using the word "interaction" and not use UUID, maybe can we get our word usages straightened out?

1. The server function in `Simplified API` is to do blind contact matching and reveal `Symptoms` only on matched contact, where the server can consider `Symptoms` to be just a blob of bytes... **except what**?  The server could usefully parse out { `QuizTime`, `DiseaseID` } to filter "old" contacts in the matching based on (a)+(b).  If nothing is disease specific, then we can treat all diseases with logic like "discard all interactions greater than 14 days" and there is just a free parameter "14" on the server.  Then `Symptoms` is just `[]byte` in `ExposureAndSymptoms` coming in `[][]byte` in `ExposureCheck` response.
_Scott: I don’t know that we have to worry too much at this stage about the exact timing of symptoms. I agree that we should age out old reports, and 14 days is a reasonable time for that.  We could implement that by only having the client ask for the last 14 days on each pull. Whether we keep old data on the server then becomes just a question of whether we have permission to store it for public health research, etc.  the “real” BLE MAC identifiers of mobile devices are hidden from apps by iOS, so the BLE protocol as it exists today generates random numbers and shares them with BLE contacts. I think your H thing might be essentially the same thing as we’re doing: if there are nuanced differences I don’t understand your proposal well enough to understand what they are_


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
_Scott: The point of the “first 64 bits” thing is to prevent some malicious client from downloading symptom information from everyone who’s reported symptoms. if you only get a partial UUID, and then have to provide the whole thing if it’s a match, you have an infinitesimal chance of being able to guess and download something that wasn’t really a match._
_Hamish:  Line 5 - the data exchanged is a fully random uuid, not a hash of each phone’s constant.  That would let clients know that you have seen the same person again which we don’t want_
_Hamish: Line 26, with Scott’s suggestion from mid last week the clients would be sending the partial contact event UUID in the initial exposure request.  This could be 4 bytes - enough to prevent the client from fishing for the whole list, and enough specificity to ensure the server is sending little information that is not required_

3. What is the role of the (optional) `GeoHash` and `Datestamp`?  Is the `GeoHash` or `Datestamp` supposed to be used for contact matching at all?   They aren't doing any work yet (except to filter out old interactions, maybe).  So, to make them do some work and store them in the backend, we should have a `/symptoms` API endpoint to filter on these attributes.  We can organize the backend tables around that and make have someone who knows Google Maps showing symptoms in a geohash and how exposurecheck works visually in simulations at least.  

_Scott: correct, the geohash is only useful for filtering/sharding/fragmenting the data in such a way that it’s easily cachable in manageable chunks, and the clients don’t have to download the array of BLE proximity identifiers for the entire world.  The datestamp was intended to be a way to implement the 14-day lookback by managing what’s requested by clients, while making each day’s data cachable on a CDN or something. if you want to instead implement that all on the server side, we could_

4. TODO: Cache expires header idea and how it relates to the "pull X times a day" design pattern


## Interested in Working on CoEpi Go Server?

Great!  Join the CoEpi Slack channel and say hi to @sourabh, or email sourabh@wolk.com --
* If you are an experienced Go person (channels, mutexes / concurrency) and enjoy ID matching, send Sourabh a note and we'll figure out what to do.

* If you are an experience iOS or Android person with a strong interest in Bluetooth Low Energy and want to do client-server bridge work, check out `mobile-app` and we will figure out what to do
