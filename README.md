# Iptablelb4
**iptablelb4** is a Layer 4 TCP load balancer built using Linux's `iptables` utilities. The project provides a straightforward API to manage load balancing rules for TCP traffic between backend servers. Users can easily configure backend server IP addresses, application ports, and select a load balancing algorithm (either Round Robin or Random). The API offers a user-friendly interface that simplifies the configuration of TCP traffic distribution across multiple servers, eliminating the need to manually work with complex `iptables` rules.

## Features

- **API Endpoints for Backend Management**: Allows you to add, list, update, or delete backend servers via HTTP.
- **TCP Load Balancing Algorithms**:
  - **Round Robin**: Distributes requests evenly across all available backend servers.
  - **Random**: Randomly selects a backend server for each incoming TCP connection.

## API Endpoints

#### `/api/v1/iptables/health`

- **Method**: `GET`
- **Response**:
```json
{
    "ping": "pong"
}
```

#### `/api/v1/iptables/add`

- **Method**: `POST`
- **Request Body**:

```json
{
    "Data": {
        "upstreams": [
            {
                "ipaddress": "192.168.100.110",
                "port": "8080"
            },
            {
                "ipaddress": "192.168.100.111",
                "port": "8080"
            }
        ],
        "algorithm": "round-robin",
        "server-farm": "web-server"
    }
}
```

- **Response**: 

```json
{
	"data": null,
	"message": "Loadbalancing rule for web-server configured successfully",
	"status": "success"
}
```

#### `/api/v1/iptables/list`

- **Method**: `GET`
- **Response**:

```json
{
	"data": [
		"web-server"
	],
	"message": "Listed all the backend server farms",
	"status": "success"
}
```

#### `/api/v1/iptables/list/<farm>`

- **Method**: `GET`
- **Response**:

```json
{
	"data": {
		"upstreams": [
			{
			"ipaddress": "172.16.251.231",
			"port": "8000"
			}
		],
		"algorithm": "random",
		"server-farm": "web-server"
	},
	"message": "Listed all the backend servers",
	"status": "success"
}
```



#### `/api/v1/iptables/update`

- **Method**: `POST`
- **Request Body**:

```json
{
    "Data": {
        "upstreams": [
            {
                "ipaddress": "192.168.100.110",
                "port": "8080"
            },
            {
                "ipaddress": "192.168.100.111",
                "port": "8080"
            }
        ],
        "algorithm": "round-robin",
        "server-farm": "web-server"
    }
}
```

- **Response**: 

```json
{
	"data": null,
	"message": "Loadbalancing rule for web-server updated successfully",
	"status": "success"
}
```
#### `/api/v1/iptables/delete/<farm>`

- **Method**: `POST`
- **Response**:

```json
{
	"data": null,
	"message": "Loadbalancing rule for web-server deleted successfully",
	"status": "success"
}
```

