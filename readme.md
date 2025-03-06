
## Project setup

```bash
$ go mod init <project_name>
```

## Compile and run the project

```bash
# run
$ go run cmd/main.go
```

## Run tests

```bash
# unit tests

```

## Deployment

When you're ready to deploy your NestJS application to production, there are some key steps you can take to ensure it runs as efficiently as possible. Check out the [deployment documentation](https://docs.nestjs.com/deployment) for more information.

If you are looking for a cloud-based platform to deploy your NestJS application, check out [Mau](https://mau.nestjs.com), our official platform for deploying NestJS applications on AWS. Mau makes deployment straightforward and fast, requiring just a few simple steps:

```bash
$ npm install -g mau
$ mau deploy
```

## Project Description 
### FlowCx billing microservice
* This is the service in charge of creating new tiers or plans for the users.
* It also handles subscriptions, payments, invoices, and other billing-related operations.


### FlowCx billing design

### 1. Plan (Tiers)
- The basic endpoint to create a new plan, which in turns create a tier or plan on Paystack.
- Considering whether there is a need to have this plans stored in our database or have it paystack and then cache it when it is called

### 2. Subscription
- This endpoint is used to create a new subscription for a tenant.
- It also handles the renewal of the subscription when it expires.
- It also handles the cancellation of the subscription.
- It also handles the upgrade or downgrade of the subscription (and also trigger the migration of data).

### 3. Transaction
- This endpoint is used to make a payment for a subscription.
- Listens to the event when a transaction has been made:
  The design to subscribe to the external paystack webhook using the `AZURE EVENT GRID`, this will listen to every transaction events sent via the webhook and it will trigger the `AZURE CLOUD FUNCTION` which will then make an update based on the tenant ID, hence it will be a `SERVERLESS ARCHITECTURE`.

### 4. Customers
- This endpoint is used to will show all the tenants using FLOWCX and their subscription status.
- It will also show the transaction history of the tenant.
- It also allows us to know and contact tenants when they leave or make stop subscribing and allow our CX team to meet and discuss with them, `EVERY CUSTOMER MATTERS`.

### 5. Discounts
- This endpoint is used to create a discount for a tenant or a lists of tenants.

### RUN LOCALLY
```bash
    npm run dev
```

### RUN DOCKER FILE  WITH THE ENVIRONMENT VARIABLES
```bash
    docker run --env-file .env myapp
```
### RUN DOCKER COMPOSE
```bash
    docker-compose up
```

### RUN DOCKER COMPOSE <SPECIFIC SERVICE>
```bash
    docker-compose up <service-name>
```

### RUN DOCKER COMPOSE <SPECIFIC SERVICE> IN DETACHED MODE
```bash
    docker-compose up -d <service-name>
```

### STOP DOCKER COMPOSE
```bash
    docker-compose down
```

### STOP DOCKER COMPOSE <SPECIFIC SERVICE>
```bash
    docker-compose down <service-name> || docker-compose stop <service-name>
```

