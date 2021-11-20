# PUF API README

## 

## Configure
The server can be configured with a secret file yaml file. The default file name is `secret.yaml` and should have the following structure:
```yaml
env: "STATE OF ENV"                 | String
server_addr: "ADDRESS OF SERVER"    | String
server_port: PORT OF SERVER         | Intger
db_addr: "ADDRESS OF DATABASE"      | String
db_port: PORT OF DATABASE           | Integer
db_username: "DATABASE USERNAME"    | String
db_password: "DATABASE PASSWORD"    | String
```
Another filename can be used by passing the name to the binary on execution with the `-path` flag, i.e. `./api -path very-secret.yaml`

## End Points

## Running



