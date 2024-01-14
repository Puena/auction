# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [proto/account.proto](#proto_account-proto)
    - [Account](#auction-Account)
    - [AccountActivateMessage](#auction-AccountActivateMessage)
    - [AccountActivatedMessage](#auction-AccountActivatedMessage)
    - [AccountBidReserveMessage](#auction-AccountBidReserveMessage)
    - [AccountCreateMessage](#auction-AccountCreateMessage)
    - [AccountCreatedMessage](#auction-AccountCreatedMessage)
    - [AccountDepositMessage](#auction-AccountDepositMessage)
    - [AccountDepositedMessage](#auction-AccountDepositedMessage)
    - [AccountGetData](#auction-AccountGetData)
    - [AccountGetDataResponse](#auction-AccountGetDataResponse)
    - [AccountReserveMessage](#auction-AccountReserveMessage)
    - [AccountReservedMessage](#auction-AccountReservedMessage)
    - [AccountSuspendMessage](#auction-AccountSuspendMessage)
    - [AccountSuspendedMessage](#auction-AccountSuspendedMessage)
    - [AccountWithdrawMessage](#auction-AccountWithdrawMessage)
    - [AccountWithdrawnMessage](#auction-AccountWithdrawnMessage)
  
    - [Account.Status](#auction-Account-Status)
  
- [proto/error.proto](#proto_error-proto)
    - [EventError](#auction-EventError)
    - [EventErrorOccurred](#auction-EventErrorOccurred)
  
- [proto/product.proto](#proto_product-proto)
    - [CommandCreateProduct](#auction-CommandCreateProduct)
    - [CommandDeleteProduct](#auction-CommandDeleteProduct)
    - [CommandUpdateProduct](#auction-CommandUpdateProduct)
    - [CreateProduct](#auction-CreateProduct)
    - [DeleteProduct](#auction-DeleteProduct)
    - [EventProductCreated](#auction-EventProductCreated)
    - [EventProductDeleted](#auction-EventProductDeleted)
    - [EventProductFound](#auction-EventProductFound)
    - [EventProductUpdated](#auction-EventProductUpdated)
    - [EventProductsFound](#auction-EventProductsFound)
    - [FindProduct](#auction-FindProduct)
    - [FindProducts](#auction-FindProducts)
    - [Product](#auction-Product)
    - [QueryFindProduct](#auction-QueryFindProduct)
    - [QueryFindProducts](#auction-QueryFindProducts)
    - [UpdateProduct](#auction-UpdateProduct)
  
- [proto/user.proto](#proto_user-proto)
    - [EventUserData](#auction-EventUserData)
    - [EventUserError](#auction-EventUserError)
    - [EventUserRevoked](#auction-EventUserRevoked)
    - [EventUserUpdated](#auction-EventUserUpdated)
    - [EventUserVerified](#auction-EventUserVerified)
    - [User](#auction-User)
    - [UserId](#auction-UserId)
    - [UserMessageError](#auction-UserMessageError)
  
- [Scalar Value Types](#scalar-value-types)



<a name="proto_account-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## proto/account.proto



<a name="auction-Account"></a>

### Account
Account represent event message for user account.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |
| user_id | [string](#string) |  |  |
| reserved_funds | [string](#string) |  |  |
| available_funds | [string](#string) |  |  |
| status | [Account.Status](#auction-Account-Status) |  |  |
| created_at | [string](#string) |  |  |






<a name="auction-AccountActivateMessage"></a>

### AccountActivateMessage



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |
| user_id | [string](#string) |  |  |






<a name="auction-AccountActivatedMessage"></a>

### AccountActivatedMessage



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| account | [Account](#auction-Account) |  |  |






<a name="auction-AccountBidReserveMessage"></a>

### AccountBidReserveMessage



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| account_id | [string](#string) |  |  |
| user_id | [string](#string) |  |  |
| bid_id | [string](#string) |  |  |
| amount | [string](#string) |  |  |






<a name="auction-AccountCreateMessage"></a>

### AccountCreateMessage
AccountCreateMessage represent event message for user account creation.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user_id | [string](#string) |  |  |






<a name="auction-AccountCreatedMessage"></a>

### AccountCreatedMessage



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| account | [Account](#auction-Account) |  |  |






<a name="auction-AccountDepositMessage"></a>

### AccountDepositMessage



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |
| user_id | [string](#string) |  |  |
| amount | [string](#string) |  |  |






<a name="auction-AccountDepositedMessage"></a>

### AccountDepositedMessage



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| account | [Account](#auction-Account) |  |  |






<a name="auction-AccountGetData"></a>

### AccountGetData



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user_id | [string](#string) |  |  |






<a name="auction-AccountGetDataResponse"></a>

### AccountGetDataResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| account | [Account](#auction-Account) |  |  |






<a name="auction-AccountReserveMessage"></a>

### AccountReserveMessage



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |
| user_id | [string](#string) |  |  |
| amount | [string](#string) |  |  |






<a name="auction-AccountReservedMessage"></a>

### AccountReservedMessage



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| account | [Account](#auction-Account) |  |  |






<a name="auction-AccountSuspendMessage"></a>

### AccountSuspendMessage



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |
| user_id | [string](#string) |  |  |






<a name="auction-AccountSuspendedMessage"></a>

### AccountSuspendedMessage



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| account | [Account](#auction-Account) |  |  |






<a name="auction-AccountWithdrawMessage"></a>

### AccountWithdrawMessage



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |
| user_id | [string](#string) |  |  |
| amount | [string](#string) |  |  |






<a name="auction-AccountWithdrawnMessage"></a>

### AccountWithdrawnMessage



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| account | [Account](#auction-Account) |  |  |





 


<a name="auction-Account-Status"></a>

### Account.Status


| Name | Number | Description |
| ---- | ------ | ----------- |
| UNKNOWN | 0 |  |
| ACTIVE | 1 |  |
| INACTIVE | 2 |  |


 

 

 



<a name="proto_error-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## proto/error.proto



<a name="auction-EventError"></a>

### EventError



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| stream_name | [string](#string) |  |  |
| consumer_name | [string](#string) |  |  |
| subject | [string](#string) |  |  |
| reference_event_key | [string](#string) |  |  |
| message | [string](#string) |  |  |
| code | [int32](#int32) |  |  |
| data | [bytes](#bytes) |  |  |
| headers | [string](#string) |  |  |
| time | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  |  |






<a name="auction-EventErrorOccurred"></a>

### EventErrorOccurred



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [EventError](#auction-EventError) |  |  |





 

 

 

 



<a name="proto_product-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## proto/product.proto



<a name="auction-CommandCreateProduct"></a>

### CommandCreateProduct



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [CreateProduct](#auction-CreateProduct) |  |  |






<a name="auction-CommandDeleteProduct"></a>

### CommandDeleteProduct



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [DeleteProduct](#auction-DeleteProduct) |  |  |






<a name="auction-CommandUpdateProduct"></a>

### CommandUpdateProduct



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [UpdateProduct](#auction-UpdateProduct) |  |  |






<a name="auction-CreateProduct"></a>

### CreateProduct



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |
| media | [string](#string) | repeated | it is a list of media urls. |
| description | [string](#string) |  |  |
| created_by | [string](#string) |  | it is user id. |






<a name="auction-DeleteProduct"></a>

### DeleteProduct



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |
| created_by | [string](#string) |  | it is user id. |






<a name="auction-EventProductCreated"></a>

### EventProductCreated



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [Product](#auction-Product) |  |  |






<a name="auction-EventProductDeleted"></a>

### EventProductDeleted



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [Product](#auction-Product) |  |  |






<a name="auction-EventProductFound"></a>

### EventProductFound



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [Product](#auction-Product) |  |  |






<a name="auction-EventProductUpdated"></a>

### EventProductUpdated



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [Product](#auction-Product) |  |  |






<a name="auction-EventProductsFound"></a>

### EventProductsFound



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [Product](#auction-Product) | repeated |  |






<a name="auction-FindProduct"></a>

### FindProduct



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |






<a name="auction-FindProducts"></a>

### FindProducts







<a name="auction-Product"></a>

### Product
Product message represent a product.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |
| name | [string](#string) |  |  |
| media | [string](#string) | repeated | it is a list of media urls. |
| description | [string](#string) |  |  |
| created_by | [string](#string) |  | it is user id. |
| created_at | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  |  |
| updated_at | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  |  |






<a name="auction-QueryFindProduct"></a>

### QueryFindProduct



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [FindProduct](#auction-FindProduct) |  |  |






<a name="auction-QueryFindProducts"></a>

### QueryFindProducts



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [FindProducts](#auction-FindProducts) |  |  |






<a name="auction-UpdateProduct"></a>

### UpdateProduct



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |
| name | [string](#string) |  |  |
| media | [string](#string) | repeated | it is a list of media urls. |
| description | [string](#string) |  |  |
| created_by | [string](#string) |  | it is user id. |





 

 

 

 



<a name="proto_user-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## proto/user.proto



<a name="auction-EventUserData"></a>

### EventUserData
EventUserData message represent event after user data request.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| data | [User](#auction-User) |  |  |






<a name="auction-EventUserError"></a>

### EventUserError
EventUserError message represent event after user message error.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| data | [UserMessageError](#auction-UserMessageError) |  |  |






<a name="auction-EventUserRevoked"></a>

### EventUserRevoked
EventUserRevoked message represent event after user revocation.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| data | [UserId](#auction-UserId) |  |  |






<a name="auction-EventUserUpdated"></a>

### EventUserUpdated
EventUserUpdated message represent event after user update.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| data | [User](#auction-User) |  |  |






<a name="auction-EventUserVerified"></a>

### EventUserVerified
EventUserVerified message represent event after user verification.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| data | [UserId](#auction-UserId) |  |  |






<a name="auction-User"></a>

### User
User messages represent user data.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |
| name | [string](#string) |  |  |
| email | [string](#string) |  |  |
| phone | [string](#string) |  |  |
| photo | [string](#string) |  |  |
| createdAt | [int64](#int64) |  |  |
| lastLogin | [int64](#int64) |  |  |
| disabled | [bool](#bool) |  |  |






<a name="auction-UserId"></a>

### UserId
UserID message represent user id data.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |






<a name="auction-UserMessageError"></a>

### UserMessageError
UserMessageError message represent data of user message error.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| msgId | [string](#string) |  |  |
| user | [User](#auction-User) |  |  |
| userId | [UserId](#auction-UserId) |  |  |
| msgError | [string](#string) |  |  |
| msgCreatedAt | [string](#string) |  |  |





 

 

 

 



## Scalar Value Types

| .proto Type | Notes | C++ | Java | Python | Go | C# | PHP | Ruby |
| ----------- | ----- | --- | ---- | ------ | -- | -- | --- | ---- |
| <a name="double" /> double |  | double | double | float | float64 | double | float | Float |
| <a name="float" /> float |  | float | float | float | float32 | float | float | Float |
| <a name="int32" /> int32 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint32 instead. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="int64" /> int64 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint64 instead. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="uint32" /> uint32 | Uses variable-length encoding. | uint32 | int | int/long | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="uint64" /> uint64 | Uses variable-length encoding. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum or Fixnum (as required) |
| <a name="sint32" /> sint32 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sint64" /> sint64 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="fixed32" /> fixed32 | Always four bytes. More efficient than uint32 if values are often greater than 2^28. | uint32 | int | int | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="fixed64" /> fixed64 | Always eight bytes. More efficient than uint64 if values are often greater than 2^56. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum |
| <a name="sfixed32" /> sfixed32 | Always four bytes. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sfixed64" /> sfixed64 | Always eight bytes. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="bool" /> bool |  | bool | boolean | boolean | bool | bool | boolean | TrueClass/FalseClass |
| <a name="string" /> string | A string must always contain UTF-8 encoded or 7-bit ASCII text. | string | String | str/unicode | string | string | string | String (UTF-8) |
| <a name="bytes" /> bytes | May contain any arbitrary sequence of bytes. | string | ByteString | str | []byte | ByteString | string | String (ASCII-8BIT) |

