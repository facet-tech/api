## Signup

AWS Congito is used for authentication. Clients utilize Cognito directly to be authenticated. 

The generated tokens (ID Token,  Access Token, Refresh Token, JWT) are used to authenticate users into the API endpoints. 

The JWT is verified utilizing the common AWS guidelines, [stated here](https://docs.aws.amazon.com/cognito/latest/developerguide/amazon-cognito-user-pools-using-tokens-with-identity-providers.html). The logic is implemented through the [JWTMiddleware](../src/facet/api/middleware/JWTMiddleware.go).

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
    $CO-->>$API: TRIGGER UpdateDB λ function \nTODO https://github.com/facets-io/api/issues/12
    $API-->>DB: UpdateDB {emain,workspaceId}
    $CO-->>$C: Verification Email {verification code}
    Note right of $CE: If either service is down, applied changes\n are undone in all the affiliated systems
    
    $C-->>$Admin: setup PW {verification code}
    $C-->>$Admin: Login {email, pw}
@endtuml
```