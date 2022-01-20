# Shopify Backend Developer Intern Challenge - Summer 2022

## Project Brief
Build an inventory tracking web application for a logistics company. Requirements include basic CRUD Functionality. You should be able to:
- Create inventory items;
- Edit them;
- Delete them; and
- View a list of them.

You can tackle this challenge using any technology you want. This is an open-ended task, but we want to focus on high quality back-end code. Custom UIs and interactive frameworks like React are not necessary.

You can use frameworks, libraries and external dependencies to help you get to the parts you are interested in building, if this helps you; or start from scratch.

## Prequisites
* Ensure that [Docker](https://docs.docker.com/get-docker/) is installed on your computer
* Download an API platform such as [Postman](https://www.postman.com/) to allow you to send HTTP requests and interact with this app

## Running the App
1. Clone the repo to your computer
2. In a terminal, navigate to the root folder of the repo and then run `docker-compose up -d`
3. Open your web browser and go to `localhost:8000/api/items` to see the app
4. Explore! Please refer to the server's [API documentation](./server/API.md) to understand how to interact with the app using Postman.

## Closing and Restarting the App
* Run `docker-compose stop` to stop the app, `docker-compose start` to restart it.
* Run `docker-compose down -v` to kill the app, wipe all database data, and remove the containers.

## Future Features
- Soft deletion and un-deletion of items
- Permanent item deletion after 30 days
