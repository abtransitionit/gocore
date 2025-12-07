
# Purpose

`OVH` is a cloud provider like `AWS`, `GCP`, `Azure`, .... This `ovh` package is built on top of the `apicli` package. It provides a (Go) OVH Client that covers:
- credential and token handling
- client creation
- caching
- OVH-specific endpoints

making it easy to operate on OVH objects like **VPS** without manually building, creating or defining HTTP requests.

| Concept              | apicli usage                                      |
| -------------------- | ------------------------------------------------- |
| HTTP client          | `apicli.NewClient(domain, logger)`                |
| HTTP verb & endpoint | `apicli.Endpoint{Verb, Path}`                     |
| Request execution    | `client.Do(ctx, &apicli.Request{...}, &response)` |
| Cached tokens        | `WithBearerToken(func() string)`                  |


# Termiology

|name|extension|comment|
|-|-|-|
|SA|**S**ervice **A**ccount|
|AT|**A**ccess **T**oken|


## Vps
- represents a remote VM in the OVH cloud
- example of `vpsId`: o1u, o2a, o3r, ...
- example of `vpsName`: vps-a7a8f7f6.vps.ovh.net
## Domains

In OVH, APIs are served under different domains:

| Constant     | Description                    |comment|
| ------------ | ------------------------------ |-|
| `DOMAIN_EU`  | `"eu.api.ovh.com"`|alias EU API    |
| `DOMAIN_STD` | `"www.ovh.com"`|alias standard API |

The **domain** defines the base URL for all requests. For example, with `DOMAIN_EU`:

```
https://eu.api.ovh.com/v1/me
```


## Endpoints (`apicli.Endpoint`)

Each API action is described as an `Endpoint`:

```go
var endpointReference = apicli.MapEndpoint{
    "MeGetInfo":  {Verb: "GET", Desc: "get Me:Info", Path: "/v1/me"},
    "VpsGetDetail": {Verb: "GET", Desc: "get Vps:Detail", Path: "/v1/vps/{id}"},
}
```

* `Verb` → HTTP method (`GET`, `POST`)
* `Path` → relative URL
* `Desc` → human-readable description
* `{id}` placeholders can be replaced with actual values using `BuildPath(map[string]string{"id": "123"})`.

# Service Account (SA)


A **Service Account (SA)** is a dedicated programmatic identity used to access OVH APIs without human credentials. In the context of this package the **SA** is used by Go code to automate tasks like VPS provisioning, SSH key management, and API token handling.

The service account 
- is used by a client to request the `OVH` API
- :Id is the client Id
- :secret is the client Secret

The **Service Account** acts as:

```
Go Code → Service Account (clientId + clientSecret) → Access Token → apicli.Client → OVH API → VPS / SSH / Resources
```

It enables:

* Automated API calls
* Secure token-based authentication
* Management of VPS instances and SSH keys programmatically


## 1️⃣ Load Credentials

Go code reads the credential file with functions like:

```go
func GetSaId() (string, error)
func GetSaSecret() (string, error)
func GetAccessTokenFromFile() (string, error)
```

These functions provide a **centralized way** to access SA info for API calls, caching tokens to avoid repeated file reads.

## 2️⃣ Create a New Access Token

Access tokens expire. To generate a new token:

```go
func CreateAccessTokenForServiceAccount(ctx context.Context, logger logx.Logger) (string, error)
```

* Reads `clientId` and `clientSecret`
* Calls OVH OAuth2 endpoint with `grant_type=client_credentials`
* Receives a new `access_token`
* Optionally updates `credential.json` with the new token


## 3️⃣ Use SA for VPS Operations

Example: reinstall a VPS:

```go
func VpsReinstallHelper(ctx context.Context, logger logx.Logger, vpsNameId string) (jsonx.Json, error)
```

**Workflow:**

1. Load `sshKeyId` from `credential.json`
2. Load VPS info from `vps.json`
3. Build request with OS image and SSH public key
4. Execute request via `apicli.Client` using the SA's access token

> SA allows **full automation** of VPS operations without user credentials.


# Clients (`apicli.Client`)

The client is responsible for sending requests.

```go
client := apicli.NewClient(DOMAIN_EU, logger)
```

You can attach a **bearer token** (OAuth) for authenticated requests:

```go
client.WithBearerToken(GetAccessTokenFromFileCached)
```

This ensures that all requests made with this client include the proper authentication header.


## 1️⃣ Create Cached Clients

To avoid repeatedly creating clients (which is costly), the package caches them with `sync.Once`.

```go
var onceOvhClient sync.Once
var OvhClientCached *apicli.Client

func GetOvhClientCached(logger logx.Logger) *apicli.Client {
    onceOvhClient.Do(func() {
        OvhClientCached = apicli.NewClient(DOMAIN_EU, logger).
            WithBearerToken(GetAccessTokenFromFileCached)
    })
    return OvhClientCached
}
```

* **`sync.Once`** ensures the client is created only once per program run.
* Cached clients are thread-safe and can be reused across goroutines.

Similarly, `GetOvhClientTokenCached` creates a non-bearer client for token generation:

```go
OvhClientTokenCached = apicli.NewClient(DOMAIN_STD, logger)
```


## 2️⃣ Handle Access Tokens

**Read Token from File**

Tokens are stored in a JSON credential file:

```go
func GetAccessTokenFromFile() (string, error)
```

* `ServiceAccount.AccessToken` is extracted.
* Cached with `sync.Once` for performance:

```go
var cachedToken string
var tokenOnce sync.Once

func GetAccessTokenFromFileCached() (string, error) {
    var err error
    tokenOnce.Do(func() {
        cachedToken, err = GetAccessTokenFromFile()
    })
    return cachedToken, err
}
```


**Refresh Token**

If the token is missing or expired, generate a new one:

```go
func CreateAccessTokenForServiceAccount(ctx context.Context, logger logx.Logger) (string, error)
```

* Reads `ClientID` and `ClientSecret` from credentials
* Calls `POST /auth/oauth2/token` via a raw `apicli.Client`
* Returns a new access token for the service account

The refreshed token is saved back to the credential file using:

```go
func updateToken(newToken string) error
```


## 3️⃣ Manage SSH Keys

OVH uses SSH keys for VPS authentication. Functions in the package allow the following operations:

* **List keys:**

```go
keys, err := SshKeyGetList(ctx, logger)
```

* **Get a key detail:**

```go
detail, err := SshKeyGetDetail(ctx, logger, sshKeyId)
```

* **Get a public key string:**

```go
pubKey, err := SshKeyGetPublicKeyCached(ctx, logger, sshKeyId)
```

> These use the cached client `GetOvhClientCached` internally.


## 4️⃣ VPS Operations


**List VPS**

```go
vpsList, err := VpsGetList(ctx, logger)
```

* Calls `GET /v1/vps`
* Returns an array of VPS identifiers

**Get VPS Details**

```go
vpsDetail, err := VpsGetDetail(ctx, logger, vpsNameId)
```

* Calls `GET /v1/vps/{id}`
* Returns full VPS info in JSON

**Reinstall a VPS**

```go
vpsInfo, err := VpsReinstallHelper(ctx, logger, vpsNameId)
```

Steps involved:

1. Retrieve SSH key ID from file (cached)
2. Retrieve public SSH key (cached)
3. Resolve VPS OS image ID
4. Call `POST /v1/vps/{id}/rebuild` with:

```go
type VpsReinstallParam struct {
    DoNotSendPassword bool
    ImageId           string
    PublicSshKey      string
}
```

---

## **5️⃣ Cached VPS Lists and Images**

The package provides **cached access to static VPS lists and OS image maps**:

```go
listVps, err := GetListVpsFromFileCached()
osImageMap := GetListOsImageFromStruct()
```

* Avoids reading JSON from disk multiple times.
* Helps quickly map VPS distro to an OS image ID.


# Putting it Together: an example workflow

**purpose**: reinstall a VPS

```go
ctx := context.Background()
logger := logx.logger

vpsName := "vps-9c33782a.vps.ovh.net"

// 1️⃣ Reinstall VPS
vpsInfo, err := VpsReinstallHelper(ctx, logger, vpsName)
if err != nil {
    logger.Errorf("VPS reinstall failed: %v", err)
}

// 2️⃣ Check if VPS is ready
ready, err := CheckVpsIsReady(ctx, logger, vpsName)
if err != nil {
    logger.Errorf("Failed to check VPS state: %v", err)
}
logger.Infof("VPS %s ready? %v", vpsName, ready)
```

* `VpsReinstallHelper` internally calls **cached clients, SSH key retrieval, and OS image mapping**
* Uses `apicli.Client.Do` to perform all API requests



# Credential Files

the file Located at `~/wkspc/.config/ovh/credential.json`, this file stores:

* **Service Account info** (`clientId`, `clientSecret`, `access_token`)
* **SSH key ID** for VPS operations
* Optional policies and user credentials

Example:

```json
{
  "myPolicy": [
    {
      "id": "ae5b786b-fe87-57b3-8b8e-0fa1cf1666a5"
    }
  ],
  "nic": {
    "id": "34eb8aa-ovh",
    "pwd": "enckshdè§te"
  },
  "serviceAccount": {
    "access_token": "xxxxxxx",
    "clientId": "EU.4332dff0faf4e0cc",
    "clientSecret": "4332dff0faf4e0cc4332dff0faf4e0cc",
    "description": "service account to be used by GO code to request OVH API",
    "flow": "CLIENT_CREDENTIALS",
    "identity": "urn:v1:eu:identity:credential:34eb8aa-ovh/oauth2-EU.4332dff0faf4e0cc",
    "name": "my_service_name"
  },
  "sshKeyId": "MySshKeyId"
}
```

**Key Fields:**

| Field          | Description                                                   |
| -------------- | ------------------------------------------------------------- |
| `clientId`     | Unique identifier for the service account                     |
| `clientSecret` | Secret key for generating access tokens                       |
| `access_token` | Temporary token for API requests                              |
| `sshKeyId`     | ID of the SSH key used for VPS reinstallation or provisioning |
| `identity`     | Full OAuth2 identity string                                   |
| `flow`         | Authentication flow (`CLIENT_CREDENTIALS` in this case)       |
| `name`         | Human-readable service account name                           |



## `vps.json`

Located at `~/wkspc/.config/ovh/vps.json`, this file defines your VPS instances:

```json
{
	"vps01": { "DisplayName": "vm.ovh.01", "NameId": "vps-4873ab3.vps.ovh.net", "Distro": "ubuntu2504" },
	"vps02": { "DisplayName": "vm.ovh.02", "NameId": "vps-aa12987.vps.ovh.net", "Distro": "alma9" },
	"vps03": { "DisplayName": "vm.ovh.03", "NameId": "vps-987af5e.vps.ovh.net", "Distro": "rocky9" },
	"vps04": { "DisplayName": "vm.ovh.04", "NameId": "vps-cab2764.vps.ovh.ca", "Distro": "Fedora42" },
	"vps05": { "DisplayName": "vm.ovh.05", "NameId": "vps-92afb534.vps.ovh.net", "Distro": "debian12" }
}
```

**Key Fields:**:

* `DisplayName`: name for display purposes and also used by SSH client conf (`~/.ssh/config`)
* `NameId`: Actual OVH VPS identifier used in API calls
* `Distro`: OS type for VPS operations (like reinstalling a specific OS image)



# Security & Best Practices

1. **Store credentials securely** (`~/.config/ovh/credential.json` with limited permissions)
2. **Use access tokens** instead of `clientSecret` for routine API calls
3. **Rotate tokens** periodically
4. **Limit scope** of SA policies (if applicable) for least privilege
