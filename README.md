# Gator project

This is a guided poroject from [Boot.dev](www.boot.dev)

### Requirements

To run this project it is necesary to have [Go](https://go.dev/doc/install) installed

```Bash

# Remove any previous GO installation if exist

rm -rf /usr/local/go && tar -C /usr/local -xzf go1.26.1.linux-amd64.tar.gz

# Add /usr/local/go/bin to the PATH environment variable. 
export PATH=$PATH:/usr/local/go/bin

#Verify that you've installed Go
go version
```

You will also need Postgres installed 

```bash
# Run the following commands
sudo apt update
sudo apt install postgresql postgresql-contrib

#Make sure the installation worked
psql --version
```

### Instalation

To install gator you can do:

```bash
# Install from github
go install github.com/Dxnax-RS/gator@latest
```

### Usage

For this program to run correctly a file located in `~/` is needed using the following structure:

```json
{
        "Db_url": "postgres://postgres:postgres@localhost:5432/gator?sslmode=disable",
        "Current_user_name": "Username"
}
```

The file must contain your DB url of the DB running on postgres

To use this program you must have feeds url at hand and use two terminals

The first terminal is an infinite loop that searches for updates in the feeds provided:

```bash
gator agg
```

The second terminal is how you interact with the program and you can use one of the following arguments

```bash
gator agg # Infinite loop that searches for updates of the registered feeds
gator register # Registers a new user and logs it in
gator login <username> # Logs you in so the program identifies your followed feeds, the user must exist in the DB
gator reset # Resets the users table and all it's relations
gator users # Prints a list of all the registered users
gator addfeed <Feed name> <Feed url> # Registers a new feed and mark it for follow
gator feeds # Prints all the feeds registered
gator follow <Feed url> # Marks the desired feed to follow
gator following # Prints all the feeds you are following 
gator unfollow <Feed url> # Unfollows the feed
gator browse <number>(Defaults to 2) # Shows the specified number of the newer posts of the feeds you follow
```