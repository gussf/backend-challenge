# backend-challenge

## This is an implementation of an ecommerce checkout service which has an integration with a discount service (via gRPC)

<br>

# Table of Contents
1. [Prerequisites](#prerequisites)
2. [Setting Up](#setting-up)
3. [Sending requests](#sending-requests)
4. [Changing Behavior](#changing-behavior)

<br>

# Prerequisites
### In order to run this application, you will need: 
* docker
* docker-compose

<br>
<br>

# Setting Up
### The ecommerce service uses some envvars, which are consumed from the .env file by docker-compose
### Envvars are explained in more detail in the Changing Behavior section
### Getting everything setup is pretty straightforward:

<br>

```shell
# In the project's root directory, run:
docker-compose up
# Or
make compose
```

### You should see a message such as:

<br> 

```shell
$ docker-compose up
[+] Running 2/0
 ⠿ Container backend-challenge-discount-1   Created            0.0s
 ⠿ Container backend-challenge-ecommerce-1  Created            0.0s
Attaching to backend-challenge-discount-1, backend-challenge-ecommerce-1
backend-challenge-discount-1   | 2021/11/09 20:01:57 Starting discount server
backend-challenge-ecommerce-1  | 2021/11/09 20:01:58 Starting ecommerce server on 0.0.0.0:3000
backend-challenge-ecommerce-1  | 2021/11/09 20:01:58 Black friday: November 9
``` 

<br>
<br>


# Sending Requests

## The ecommerce service expects HTTP connections on localhost:3000 by default, only allowing <b>POST method</b> on /checkout endpoint
## If you wish to change the listen address, see 'Changing Behaviour' section

<br>

## Request example:
Endpoint: <b>localhost:3000/checkout</b> <br>
HTTP Method: <b>POST</b>

```json
{
    "products": [
        {
            "id": 1,
            "quantity": 1
        }
    ]
}
```

<br>

## Response example:

```json
{
    "total_amount": 15157,
    "total_amount_with_discount": 15157,
    "total_discount": 0,
    "products": [
        {
            "id": 1,
            "quantity": 1,
            "unit_amount": 15157,
            "total_amount": 15157,
            "discount": 0,
            "is_gift": false
        }
    ]
}
```

<br> 
<br> 

# Changing Behavior
### Through the use of envvars, the application allows you to change its behavior by modifying the <b>.env file</b>
### docker-compose should read the .env file automatically

<br> 

## <b><u>Black Friday</b></u>
BLACK_FRIDAY_DATE_MMDD - Black friday date, in MMDD format
```shell
# Example: If you want BlackFriday on December 2nd
export BLACK_FRIDAY_DATE_MMDD=1202
```

<br>

## <b><u>Endpoints</b></u>
ECOMMERCE_LISTEN_ADDRESS - "IP:port" that the ecommerce service will listen on
```shell
# Example
export ECOMMERCE_LISTEN_ADDRESS="0.0.0.0:1234"
```
<br>

DISCOUNT_GRPC_ADDRESS - "IP:port" on which the ecommerce will reach the discount gRPC server
```shell
# Example: 'discount' refers to the container name used on docker-compose.yaml
export DISCOUNT_GRPC_ADDRESS="discount:50051"
```

<br>

## <b><u>gRPC Deadline</b></u>
GRPC_DEADLINE_MS - Amount of milliseconds the service will wait for a response from discount service
```shell
# Example: Will wait 50ms for each gRPC call 
export GRPC_DEADLINE_MS=50
```
