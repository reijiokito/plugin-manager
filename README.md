# Sigma Plugin Manager

Sample:
1. Build plugin: 
- cd /core
- go build -buildmode plugin -o /usr/local/sigma/go-plugins/plugin_a ./plugins/plugin_a/main.go
- go build -buildmode plugin -o /usr/local/sigma/go-plugins/plugin_b ./plugins/plugin_b/main.go

2. Run:
- cd /core
- go run main.go

**------------------------------------------------------------------


# APIs

1. http://localhost:8000/plugin/get-all
2. http://localhost:8000/plugin/instance/get-all
3. http://localhost:8000/plugin/init
4. http://localhost:8000/plugin/dump/:pluginName
5. http://localhost:8000/plugin/dump-all
6. http://localhost:8000/plugin/instance/service/close/:id
7. http://localhost:8000/plugin/instance/built-in/close/:id

1. Get all plugin information:
* Request: http://localhost:8000/plugin/get-all
   
* Response: 
```json
{
    "data": [
        {
            "Name": "nats",
            "ModTime": "0001-01-01T00:00:00Z",
            "LoadTime": "0001-01-01T00:00:00Z",
            "Phases": [
                "access"
            ],
            "Version": "1.0.0",
            "Priority": 1,
            "Schema": {
                "Name": "nats",
                "fields": [
                    {
                        "Config": {
                            "fields": [
                                {
                                    "natsurl": {
                                        "type": "string"
                                    }
                                },
                                {
                                    "natsusername": {
                                        "type": "string"
                                    }
                                },
                                {
                                    "natspassword": {
                                        "type": "string"
                                    }
                                }
                            ],
                            "type": "record"
                        }
                    }
                ]
            }
        },
        {
            "Name": "plugin_a",
            "ModTime": "0001-01-01T00:00:00Z",
            "LoadTime": "0001-01-01T00:00:00Z",
            "Phases": [
                "access"
            ],
            "Version": "",
            "Priority": 0,
            "Schema": {
                "Name": "plugin_a",
                "fields": [
                    {
                        "Config": {
                            "fields": [
                                {
                                    "name": {
                                        "type": "string"
                                    }
                                }
                            ],
                            "type": "record"
                        }
                    }
                ]
            }
        },
        {
            "Name": "plugin_b",
            "ModTime": "0001-01-01T00:00:00Z",
            "LoadTime": "0001-01-01T00:00:00Z",
            "Phases": [
                "access"
            ],
            "Version": "",
            "Priority": 0,
            "Schema": {
                "Name": "plugin_b",
                "fields": [
                    {
                        "Config": {
                            "fields": [
                                {
                                    "name": {
                                        "type": "string"
                                    }
                                }
                            ],
                            "type": "record"
                        }
                    }
                ]
            }
        },
        {
            "Name": "plugin_c",
            "ModTime": "0001-01-01T00:00:00Z",
            "LoadTime": "0001-01-01T00:00:00Z",
            "Phases": [
                "access"
            ],
            "Version": "1.0.1",
            "Priority": 2,
            "Schema": {
                "Name": "plugin_c",
                "fields": [
                    {
                        "Config": {
                            "fields": [
                                {
                                    "name": {
                                        "type": "string"
                                    }
                                },
                                {
                                    "subject": {
                                        "type": "string"
                                    }
                                }
                            ],
                            "type": "record"
                        }
                    }
                ]
            }
        },
        {
            "Name": "plugin_d",
            "ModTime": "0001-01-01T00:00:00Z",
            "LoadTime": "0001-01-01T00:00:00Z",
            "Phases": [
                "access"
            ],
            "Version": "",
            "Priority": 0,
            "Schema": {
                "Name": "plugin_d",
                "fields": [
                    {
                        "Config": {
                            "fields": [
                                {
                                    "name": {
                                        "type": "string"
                                    }
                                },
                                {
                                    "subject": {
                                        "type": "string"
                                    }
                                }
                            ],
                            "type": "record"
                        }
                    }
                ]
            }
        }
    ]
}
```

2. Get all plugin instances
* Request: http://localhost:8000/plugin/instance/get-all
* Response:
```json
{
    "data": [
        {
            "Id": 1,
            "Name": "plugin_a",
            "Modtime": "2023-05-08T11:48:46.928932175+07:00",
            "Loadtime": "2023-05-08T14:48:17.302103563+07:00",
            "LastStartInstance": "2023-05-08T14:48:18.320837846+07:00",
            "LastCloseInstance": "0001-01-01T00:00:00Z"
        },
        {
            "Id": 2,
            "Name": "plugin_b",
            "Modtime": "2023-05-08T11:59:26.318194458+07:00",
            "Loadtime": "2023-05-08T14:48:17.307845269+07:00",
            "LastStartInstance": "2023-05-08T14:48:18.320933951+07:00",
            "LastCloseInstance": "0001-01-01T00:00:00Z"
        },
        {
            "Id": 3,
            "Name": "plugin_c",
            "Modtime": "2023-05-08T13:47:49.442587685+07:00",
            "Loadtime": "2023-05-08T14:48:17.313707996+07:00",
            "LastStartInstance": "2023-05-08T14:48:18.321043431+07:00",
            "LastCloseInstance": "0001-01-01T00:00:00Z"
        },
        {
            "Id": 0,
            "Name": "nats",
            "Modtime": "2023-05-08T11:58:26.982656106+07:00",
            "Loadtime": "2023-05-08T14:48:17.296685294+07:00",
            "LastStartInstance": "2023-05-08T14:48:17.319476184+07:00",
            "LastCloseInstance": "0001-01-01T00:00:00Z"
        }
    ]
}
```
3. Initialize a plugin
After build a plugin executable file and move to plugin manager directory.
Server will create yaml file config for the plugin and save in repository.
Then, server will start plugin instance and run the plugin's main handler function (Access).
* Request: http://localhost:8000/plugin/init
```json
{
    "name": "plugin_d",
    "config": {
        "name":"PLUGIN_D",
        "subject": "hihi"
    }
}
```
* response:
```json
    {
  "data": "Success"
}
```
4. Dump plugin executable file
When you want to dump new version of plugin executable file and ready to release new version of plugin.
Load plugin information.

* Request: http://localhost:8000/plugin/dump/:pluginName
* Response:
```json
{
    "data": {
        "Name": "plugin_c",
        "ModTime": "0001-01-01T00:00:00Z",
        "LoadTime": "0001-01-01T00:00:00Z",
        "Phases": [
            "access"
        ],
        "Version": "1.0.0",
        "Priority": 1,
        "Schema": {
            "Name": "plugin_c",
            "fields": [
                {
                    "Config": {
                        "fields": [
                            {
                                "name": {
                                    "type": "string"
                                }
                            },
                            {
                                "subject": {
                                    "type": "string"
                                }
                            }
                        ],
                        "type": "record"
                    }
                }
            ]
        }
    }
}
```
5. Dump all plugin executable files which already exited in plugin manager directory
When you want to dump all new version of plugin executable files and ready to release new version of plugins.
Load all plugins information 
* Request: http://localhost:8000/plugin/dump-all
* Response: 
```json
{
    "data": [
        {
            "Name": "nats",
            "ModTime": "0001-01-01T00:00:00Z",
            "LoadTime": "0001-01-01T00:00:00Z",
            "Phases": [
                "access"
            ],
            "Version": "1.0.0",
            "Priority": 1,
            "Schema": {
                "Name": "nats",
                "fields": [
                    {
                        "Config": {
                            "fields": [
                                {
                                    "natsurl": {
                                        "type": "string"
                                    }
                                },
                                {
                                    "natsusername": {
                                        "type": "string"
                                    }
                                },
                                {
                                    "natspassword": {
                                        "type": "string"
                                    }
                                }
                            ],
                            "type": "record"
                        }
                    }
                ]
            }
        },
        {
            "Name": "plugin_a",
            "ModTime": "0001-01-01T00:00:00Z",
            "LoadTime": "0001-01-01T00:00:00Z",
            "Phases": [
                "access"
            ],
            "Version": "",
            "Priority": 0,
            "Schema": {
                "Name": "plugin_a",
                "fields": [
                    {
                        "Config": {
                            "fields": [
                                {
                                    "name": {
                                        "type": "string"
                                    }
                                }
                            ],
                            "type": "record"
                        }
                    }
                ]
            }
        },
        {
            "Name": "plugin_b",
            "ModTime": "0001-01-01T00:00:00Z",
            "LoadTime": "0001-01-01T00:00:00Z",
            "Phases": [
                "access"
            ],
            "Version": "",
            "Priority": 0,
            "Schema": {
                "Name": "plugin_b",
                "fields": [
                    {
                        "Config": {
                            "fields": [
                                {
                                    "name": {
                                        "type": "string"
                                    }
                                }
                            ],
                            "type": "record"
                        }
                    }
                ]
            }
        }
    ]
}
```
6. Close plugin instance which is running
* Request: http://localhost:8000/plugin/instance/service/close/:id
* Response: 
```json
{
    "data": "Success"
}
```

7. Close built-in plugin instance which is running
* Request: http://localhost:8000/plugin/instance/built-in/close/:id
* response:
```json
{
    "data": "Success"
}
```

