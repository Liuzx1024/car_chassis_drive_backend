# Protocal

## request (Master to Slave)
```json
{
    "id": "8aa74234-e004-4295-b753-81ba8514de3d",
    "device": "Ultrasonic sensor",
    "operation": "read",
    "parameter": ""
}
```
* id:a UUID string,to identify a task
* device:name of a device
* operation:name of a operation to do to a device
* parameter:parameter of a operation

## response (Slave to Master)
```json
{
    "id": "8aa74234-e004-4295-b753-81ba8514de3d",
    "status": "success",
    "result": "10"
}
```
* id:a UUID string,to identify a task
* status:represent the status of a task
* result:result of the operation
