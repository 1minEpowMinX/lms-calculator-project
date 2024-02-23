# lms-calculator-project
A user wants to count arithmetic expressions. He enters the string 2 + 2 * 2 and wants to get 6 as the answer. But our addition and multiplication operations (also division and subtraction) take a "very, very" long time.  
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
go run main.go
```
This command will start your application. It will compile and run the backend and frontend servers.

## Step 4: Checking the Application

After completing the previous step, your servers will be up and running. Open a web browser and navigate to `http://localhost:<port>`, where `<port>` is the port of your frontend (8081) and backend (8080) servers. You should see the application running.

## Notes

- If you need to change ports or any other settings, you can do so by editing the corresponding configuration project files.
- Make sure all project dependencies are installed. If you encounter any errors during the startup process, ensure all necessary dependencies are installed using `go get`.

That should be all you need to run your application. Happy coding!
