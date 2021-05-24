# Analytics
Analytics is a Napptive library to monitoring operational data

This lib allows you to send data to a provider for monitoring tasks

For now, only BigQuery provider is implemented, but it is too easy to add providers. Just implementing the interface:

```
// Provider with an interface that defines the monitoring provider methods
type Provider interface {
	// SendLoginData puts a login in the database
	SendLoginData(data entities.LoginData) error
	// SendOperationData puts an operation data in the database
	SendOperationData(data entities.OperationData) error
}
```

## BigQuery provider

1.Create a client
```
func NewBigQueryProvider(projectID string, credentialsPath string, loopTime time.Duration) (Provider, error) {
```
where 
- `projectID` is the GKE Project identifier
- `credentialsPath` is the path of the credentials file. The service account credentials
- `loopTime` is the waitting time to make inserts in the database. 
  The data is stored in a cache and sent every so often

---

## Getting credentials
Read [google documentation](https://cloud.google.com/docs/authentication/production#cloud-console) to create the *ServiceAccount* and, a *key*
Once you have created them, you can share the tables with the service account:

1. Go to table schema, and click on `SHARE` option
   ![share](images/share.png)
2. Add a member
   ![Addmember](images/add_member.png)
3. Fill the member with the Service Account created and BigQuery>BigQuert Data Editor Role
   ![member](images/member.png)
   
## License

 Copyright 2020 Napptive

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

      https://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
