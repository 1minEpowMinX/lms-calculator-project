# Distributed calculator of arithmetic expressions
A user wants to count arithmetic expressions. He enters the string 2 + 2 * 2 and wants to get 6 as the answer. But our operations of addition and multiplication (also division and subtraction) take a "very, very" long time. That's why the variant when the user makes an http-request and gets the result as a response is impossible. Moreover: calculation of each such operation in our "alternative reality" takes "gigantic" computing power. Accordingly, we must be able to perform each operation separately and scale this system by adding computing power to our system in the form of new "machines". Therefore, the user, sending an expression, receives an expression identifier in response and can check with the server at regular intervals to see if the expression has been calculated. If the expression is finally calculated - he will get the result. Remember that some parts of an arithmetic expression can be calculated in parallel.

# Running Instructions for the Application

This guide will help you run your application, consisting of backend and frontend servers, using the Go programming language.

## Step 1: Get the Source Code

First, you need to clone the repository of your project to your local machine. To do this, execute the following command in the terminal:
```sh
git clone https://github.com/1minEpowMinX/lms-calculator-project
```
## Step 2: Prepare Your Environment

Make sure you have Go installed on your computer. If Go is not installed, you can download and install it from the [official Go website](https://golang.org/dl/).

## Step 3: Running the Application

1. Open a terminal and navigate to the root directory of your project.
2. Execute the following command to run the application:
```sh
make run-project
```
This command will start your application. It will compile and run the backend and frontend servers.

## Step 4: Checking the Application

After completing the previous step, your servers will be up and running. Open a web browser and navigate to `http://localhost:<port>`, where `<port>` is the port of your frontend (8080) server. You should see the application running.

## Notes

- If you need to change ports or any other settings, you can do so by editing the corresponding configuration project files.
- Make sure all project dependencies are installed. If you encounter any errors during the startup process, ensure all necessary dependencies are installed using `go get`.

## How it works?

![Untitled-2024-04-22-1149](https://github.com/1minEpowMinX/lms-calculator-project/assets/129176682/bd57ae1c-8c08-49f4-aa1d-e11bf7cfe726)
  
## Endpoint Descriptions

- ```POST /auth/signup/``` - this endpoint is for registering new users. It accepts a JSON object with a username and password. If a user with the specified name already exists, it returns an error. Otherwise, it creates a new account and returns success status.
- ```POST /auth/login/``` - this endpoint is designed for user authentication. It accepts a JSON object with username and password. After successful authentication, it sets a cookie with the authentication token and redirects to the main page. In case of authentication error, it returns an appropriate error.
- ```POST /expression/``` - this endpoint is designed to store a math expression in the database. It accepts a POST request with a JSON object containing the expression. Access to this endpoint requires authorization by a token, which is passed as a cookie. The expression is stored along with the user ID and date and status metadata. In case of successful saving, it returns a success status.
- ```GET /expression/``` - this endpoint is designed to retrieve a list of all saved expressions from the database for a particular user. Access to this endpoint also requires token authorization. Returns a JSON list of all saved expressions of the user, including their IDs, expressions, responses, dates and statuses.
- ```DELETE /expression/{id}/``` - this endpoint is designed to delete a stored expression from the database by its identifier. It accepts a DELETE request with the parameter of the expression identifier in the URL. In case of successful deletion, it returns the confirmation status.

That should be all you need to run your application. Happy coding!
