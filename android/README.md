# coepi-android-app

Very simple client android application for this client-server flow:

1. Contact uuid and information is generated (Button "new contact")
2. User records symptoms in their CoEpi app, (Button "New Symptom"), resulting in a POST to `/exposureandsymptoms` endpoint with their `Symptoms` and UUIDs/Dates of `Contacts`
3. Client poll for symptoms of contacted with POST to `/exposurecheck` to see if the server has received any symptoms of the above matching those that the device has seen.


## Endpoints

Test Endpoint: https://coepi.wolk.com:8081
