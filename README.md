# aspire-loan-app
A mini app that allows users to go through a loan application. Authenticated users can create a loan and submit weekly loan repayments

### System Requirements
- `docker`
- `docker-compose`

### Installation
- `make build`
- `docker-compose up`

The `make build` command will build the project and generate a docker image for the app

The `docker-compose up` command will start a mysql container and the aspire loan app container

Once the containers are up, the loan application can be accessed at `localhost:8080`

### About

The app exposes the following APIs
1. `POST /users`
 This endpoint is used to create users in the system. The API request looks like 
```
{
    "id": <user_id>
    "password": <password>
}
```
This endpoint does not require any auth

**Note**: The application has 1 admin user with the id `aspire` and the password `aspire` 


2. `GET /users/<user_id>/loans`
This endpoint is used to fetch all loans for a user. This endpoint requires Basic Auth

3. `POST /users/<user_id>/loans`
This endpoint is used to create a loan. The API request looks like
```
{
    "amount": <amount>
    "term": <term>
}
```
This endpoint requires Basic Auth

4. `GET /users/<user_id>/loans/<loan_id>`
This endpoint is used to fetch details of a loan using the loanID. This endpoint requires Basic Auth

5. `PUT /users/<user_id>/loans/<loan_id>`
This endpoint is used to update the status of a loan to `APPROVED`. Only the admin user (`aspire`) can access this endpoint. The API request looks like 
```
{
    "status": "APPROVED"
}
```
6. `POST /users/<user_id>/loans/<loan_id>/payments`
This endpoint is used to add a repayment corresponding to the scheduled repayment. This endpoint also updates the status of a scheduled repayment to `PAID`. It will also update the loan status to `PAID` if all the scheduled repayments have been moved to `PAID` status. The API request looks like
```
{
    "scheduled_repayment_id": <id>,
    "amount": <amount>
}
```