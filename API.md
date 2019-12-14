# cURL

### Install
```$ brew install jq > curl ... | jq```

##### GET SignIn
```
$ curl -u name:passwd "localhost:8080/sign-in
```

##### GET Verify
```
$ curl -H "Authorization: Bearer REPLACE_WITH_TOKEN" \
 -X GET "localhost:8080/verify"
```