## Signup

AWS Congito is used for authentication. Clients utilize Cognito directly to be authenticated. 

The generated tokens (ID Token,  Access Token, Refresh Token, JWT) are used to authenticate users into the API endpoints. 

The JWT is verified utilizing the common AWS guidelines, [stated here](https://docs.aws.amazon.com/cognito/latest/developerguide/amazon-cognito-user-pools-using-tokens-with-identity-providers.html). The logic is implemented through the [JWTMiddleware](../src/facet.ninja/api/middleware/JWTMiddleware.go).

Sequence diagram Signup:

```plantuml
   @startuml
    !$CE = "client_email"
    !$C = "User"
    !$Admin = "FN_Application"
    !$API = "FN_API"
    !$DB = "DynamoDB"
    !$CO = "Cognito"

    $C-->>$CO: Signup {email,workspaceId}
    $CO-->>$API: TRIGGER UpdateDB λ function
    $API-->>DB: UpdateDB {emain,workspaceId}
    $CO-->>$C: Verification Email!
    Note right of $CE: If either service is down, applied \nchanges are undone (cascade)
    $CO-->>$CE: temp PW
    $C<-->$CE: Retreive PW
    $C-->>$Admin: setup PW
    $C-->>$Admin: Login {email, pw}
@endtuml
```

## Signup a test account

// TODO script