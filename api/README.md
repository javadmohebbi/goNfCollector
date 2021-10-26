# Netflow Collector API Guide
Netflow Collector API Guide helps developers integrate their own apps with netflow collector application. All of the outputs are JSON

This api is listening based on configuration file `api_server.yml` provided by `-confPath /path/to/directory/of/configuration/files/` (must be ended with `/`). The default value is `0.0.0.0:9999`

To check if the API server is running you can open `http://0.0.0.0:9999/v1/api/` (or `https://0.0.0.0:9999/v1/api/` if TLS is enabled)

## API Version
Current version is 1 and all the sub routes are started with:
- `/v1/api`


# Devices
All the routes related to Netflow Exporter Devices are started with `/v1/api/device` and here are the sub routes:

## Get all netflow exporter devices:
  - **Route**: `/v1/api/device/get/all`
  - **Description**: Returns all of netflow exporter devices which has at least one record in the database
  - **Methods**: `GET`, `OPTIONS`
  - **Input**: N/A
  - **Output**: an array of devices
    ```
    [
        {
            "ID":1,
            "CreatedAt":"2021-06-15T12:46:06.004006+04:30",
            "UpdatedAt":"2021-06-15T12:46:06.004006+04:30",
            "DeletedAt":null,
            "Device":"127.0.0.1",
            "Name":"",
            "Info":""
        }
    ]
    ```
## Device: Get summary of all device based on interval
  - **Route**: `/v1/api/device/get/summary/interval/{interval}`
  - **Description**: Returns information about the last summary of devices base on an `interval` like 1m, 2h ...
    - `interval` is the combinations of an integer with the letter at the end of it which repesent specific amount of time:
      - **s**: Seconds
      - **m**: Minutes
      - **h**: Houres
      - **d**: Days
      - **w**: Weeks
      - **M**: Month
      - **y**: Years
    - If provided `interval` value was invalid it will replace it with **15m** (15 minutes)
  - **Methods**: `GET`, `OPTIONS`
  - **Input**: N/A
  - **Output**: An array of all devices with their summary
    ```
    [
        {
            "device_id": 1,
            "device": "127.0.0.1",
            "flow_count": 36758,
            "device_name": "",
            "device_info": "",
            "total_bytes": 8848500,
            "total_packets": 89610
        }
    ]
    ```
## Device: Get grouped summary based on Interval
  - **Route**: `/v1/api/device/get/summary/group/{interval}/by/{DeviceID}`
  - **Description**: returns an array of packets, bytes .... group them by the interval that a user provides. It will automatically checks the provided interval and group them with the most effective time it can. For example if user provid `1h` it will group the series to `every minute` but if `24h` it will group them by `day` and if user provide `3h` it will group them by `hour`
    - `interval` is like previous method
    - `deviceID` id the ID of netflow exporter device
  - **Methods**: `GET`, `OPTIONS`
  - **Input**: N/A
  - **Output**: An array of device summary metrics-series group by the provided `interval`
    - In this example this request has been sent to the API server `/v1/api/device/get/summary/group/15m/by/1` and the below **JSON** array is the response. As you can see interval is **15m** and results grouped into **minutes**
        ```
            [
                }
                  "device":"127.0.0.1",
                    "flow_count":5706,
                    "device_name":"",
                    "device_info":"",
                    "total_bytes":2107179,
                    "total_packets":15175
                },
                {
                    "_time":"2021-06-15T16:12:00+04:30",
                    "device_id":1,
                    "device":"127.0.0.1",
                    "flow_count":3355,
                    "device_name":"",
                    "device_info":"",
                    "total_bytes":701587,
                    "total_packets":7465
                },
                {
                    "_time":"2021-06-15T16:13:00+04:30",
                    "device_id":1,
                    "device":"127.0.0.1",
                    "flow_count":2051,
                    "device_name":"",
                    "device_info":"",
                    "total_bytes":751719,
                    "total_packets":5820
                },
                {
                    "_time":"2021-06-15T16:14:00+04:30",
                    "device_id":1,
                    "device":"127.0.0.1",
                    "flow_count":1795,
                    "device_name":"",
                    "device_info":"",
                    "total_bytes":605859,
                    "total_packets":4890
                },
                {
                    "_time":"2021-06-15T16:15:00+04:30",
                    "device_id":1,
                    "device":"127.0.0.1",
                    "flow_count":1120,
                    "device_name":"",
                    "device_info":"",
                    "total_bytes":683647,
                    "total_packets":3270
                },
                {
                    "_time":"2021-06-15T16:16:00+04:30",
                    "device_id":1,
                    "device":"127.0.0.1",
                    "flow_count":354,
                    "device_name":"",
                    "device_info":"",
                    "total_bytes":150363,
                    "total_packets":805
                }
            ]
        ```