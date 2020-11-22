## Signup

Signup should not go through the API. It should directly go through cognito.

Current state:

```plantuml
   @startuml
    !$CE = "client_email"
    !$C = "User"
    !$Admin = "facet.ninja"
    !$API = "FN_API"
    !$DB = "DynamoDB"
    !$CO = "Cognito"

    $C-->>$CO: Signup {email,workspaceId}
    $CO-->>$API: TRIGGER UpdateDB {email,workspaceId}
    $API-->>DB: UpdateDB {emain,workspaceId}
    $CO-->>$C: Check your email!
    Note right of $CE: If either service is down,\n applied changes are undone
    $CO-->>$CE: temp PW
    $C<-->$CE: Retreive PW
    $C-->>$Admin: setup PW
    $C-->>$Admin: Login {email, pw}
@endtuml
```

## Signup a test account

// TODO script