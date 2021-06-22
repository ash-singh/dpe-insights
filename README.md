# dpe-insights

## Config Parameters
### Add Git Oauth2 Token

https://github.com/settings/tokens

create a copy of .env.example as .env and 
set own values for env variables


```bash
cp .env.example .env

vim .env 

PLUGIN_GITHUB_OAUTH2_TOKEN=<you git oauth2 token>
```

## Project Initialization

```shell script
    $ docker-compose build
    $ docker-compose up -d 
``` 

## App

## Build commands executables
```shell script
    $ make build
``` 

### Import and transform daily GitHub data
```bash
    bash ./cmd/daily_import_and_transform.sh
```
    

### DPE cli tool documentation
* [dpe-cli](docs/cli/.md)

### Import daily pull requests count data for last n days
```shell script
    $ ./build/dpe-cli github import pull-request-counts --days=30
``` 
    

### Import all pull requests data for a range startDate and endDate or number of days
```shell script
    $ ./build/dpe-cli github import pull-requests --days=30
    $ ./build/dpe-cli github import pull-requests --startDate=2020-11-01 --endDate=2020-11-07
``` 
    
### Import Teams related data
```shell script
    $ ./build/dpe-cli github import teams
```
    
### Transform pull requests data
```shell script
    $ ./build/dpe-cli github transform pull-requests
``` 

### Transform requests count data for last n days
```shell script
    $ ./build/dpe-cli github transform pull-request-counts --days=30
``` 
     
### Transform team - pull request count data for last n weeks
```shell script
    $ ./build/dpe-cli github transform team-pull-request-counts --weeks=10
``` 
    
### Sync imported pull request data
```shell script
    $ ./build/dpe-cli github sync  
```    

## Provisioning Grafana Dashboard
You can manage dashboards in Grafana by adding one or more YAML config files in the: 
```
.docker/grafana/var/lib/grafana/provisioning/dashboards/
```
directory. By default, these dashboards cannot be updated from Grafana UI, and an app restart
is required to reflect any changes in JSON files. [Docs](https://grafana.com/docs/grafana/latest/administration/provisioning/#dashboards) 


### API details
Check api/routing.go file
    
### Grafana Dashboard
Grafana is accessible under the following URL: http://localhost:3002 using the following credentials:
- User: _admin_
- Password: _admin123_

### For Development
#### Live reload of app

https://github.com/cosmtrek/air

#### Lint and Fix
```shell script
 <project_root> make lint
```
```shell script
 <project_root> make lint-fix
```

#### Tests
```shell script
 <project_root> make test
```
