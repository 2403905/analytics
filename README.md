# Analytics

The Analytics is a CLI application that has the outputs:

[1] Top 10 active users sorted by amount of PRs created and commits pushed  
[2] Top 10 repositories sorted by amount of commits pushed  
[3] Top 10 repositories sorted by amount of watch events

### Installation
* Clone a repository
* Call the `make` command to set up the required dependencies and run an application.  
Notice: The `make` command call the wget and tar tools to set up the required dependencies. If you don't have such tools please install them or download the required dependencies manually.
Download and unzip to the root project directory a data.tar.gz file https://github.com/adjust/analytics-software-engineer-assignment/raw/master/data.tar.gz
* Call the `make build` command to build a binary file.
* Call the `make start` command to run an application.
* Call the `make test` command to run tests and benchmarks.
