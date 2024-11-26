This is a very basic version of a leaderboard for GameOff-2024 game jam. It is not "secure" but will store user data in an aws table. 

There are two internal packages aside form main, auth and handler. 
- auth handles JWT tokens for the client and service to ensure you cant just send a request from postman withought knowing the keys.
- handler handles the various dynamo requests for the api calls. 
This can be improved apon but for now it will work for the game jam.