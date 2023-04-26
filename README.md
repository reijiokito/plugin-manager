** Sigma Plugin Manager

Sample:
1. Build plugin: 
- cd /core
- go build -buildmode plugin -o /usr/local/sigma/go-plugins/plugin_a ./plugins/plugin_a/main.go
- go build -buildmode plugin -o /usr/local/sigma/go-plugins/plugin_b ./plugins/plugin_b/main.go

2. Run:
- cd /core
- go run main.go