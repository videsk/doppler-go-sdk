# Doppler Go SDK

To start using this SDK, just:

## 1. Install the package

```bash
go get github.com/videsk/doppler-go-sdk
```

## 2. Import the package

Using the method `GetSecrets` you will retrieve all secrets in once.

```go
package main

import (
	"log"
	
	"github.com/videsk/doppler-go-sdk"
)

func main() {
	
	dopplerClient, err := doppler.NewDoppler("projectName", "secretKey", "environment")
    if err != nil {
        log.Fatalf("Error with Doppler instance: %v", err)
    }
	
    dopplerClient.SetHTTPClient(http.DefaultClient)
	
    secrets, err := dopplerClient.GetSecrets()
    if err != nil {
        log.Fatalf("Error getting the secrets: %v", err)
    }

    for key, value := range secrets {
        log.Printf("Error getting: %s, Value: %s\n", key, value)
    }
	
	// The rest of the code ...
}
```

## 3. Get specific secret by name

Using the method `GetOne` you will retrieve one secret in once, but calling `getSecrets`. Which is not efficient if you need get more than one at the same time.

```go
package main

import (
	"log"
	
	"github.com/videsk/doppler-go-sdk"
)

func main() {
	
	dopplerClient, err := doppler.NewDoppler("projectName", "secretKey", "environment")
    if err != nil {
        log.Fatalf("Error with Doppler instance: %v", err)
    }
	
    dopplerClient.SetHTTPClient(http.DefaultClient)
	
	secretName := "my_secret_name"
    secretValue, err := dopplerClient.GetOne(secretName)
    if err != nil {
        log.Fatalf("Error getting the secret %s: %v", secretName, err)
    }

    log.Printf("The secret value is '%s': %s\n", secretName, secretValue)
	
	// The rest of the code ...
}
```
## Efficient single value

In case you want retrieve multiples values you can just:

```go
secrets, err := dopplerClient.GetSecrets()

if err != nil {
    log.Fatalf("Error getting the secrets: %v", err)
}

mySecret1 := secrets["my_secret_name"]
mySecret2 := secrets["my_super_secret"]
```