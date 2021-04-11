#Portfolio

An extremely basic tracker for your assets. This is able to track prices of most stocks, bonds & cryptos.

###Project Structure
The project consists of a golang api and vue.js UI. The golang API embeds and serves the UI


####Usage
1. Clone the repo
2. Create a directory in your home folder called Portfolio. 
```shell
mkdir ~/Porfolio
```
3. Use the makefile to install the project. This will compile the application and copy it to your
Portfolio home directory.
```shell
make deploy
```
4. Create the portfolio database file
```shell
cd ~/Porfolio 
touch portfolio.db
```
5. Run the portfolio migration
```shell
cd ~/Porfolio 
./portfolio -runMigration=true
```
6. Run the portfolio app
```shell
cd ~/Porfolio 
./portfolio -runMigration=true
```
7. Navigate to http://127.0.0.1:5000 in a browser


####Arguments
The portfolio binary supports the following arguments
1. `db` - This allows you to set the database path. Example usage `porfolio -db=/home/user/dbfile.db`
2. `runHistory` - The portfolio app runs a cron to get the current value daily by 17:00. This argument forces the app to get the value of the portfolio now.
3. `runMigration` - This runs the database migration. See the migration function in the cmd/main.go file.
