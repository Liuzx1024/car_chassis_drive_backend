# Design of Serial-Bus
&emsp;&emsp;Use just one serial port,serve serval device.<br>
&emsp;&emsp;Contain one master and several slaves.
# Hardware

## Master
0. Set CE pin to LOW
1. Send data in buf to Arduino
2. Send FINISH signal to Arduino
3. Copy data to buf from Arduino
4. Waiting for FINISH signal
5. Set CE pin to HIGH

## Slave
0. Enter handler function if interrupt by CE
1. Copy data to buf from raspberryPi
2. Waiting for FINISH signal from raspberryPi
3. Send data in buf to raspberryPi
4. Send FINISH signal to raspberryPi
5. Exit interrupt function

# Software

## Protocal

### request (Master to Slave)
```json
{
    "id": "8aa74234-e004-4295-b753-81ba8514de3d",
    "task":
        {
            "device": "Ultrasonic sensor",
            "operation": "read",
            "parameter": ""
        }
}
```
* id:a UUID string,to identify a task
* task:the content of a task
* device:name of a device
* operation:name of a operation to do to a device
* parameter:parameter of a operation

### response (Slave to Master)
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

### ***FINISH*** signal
```json
{}
```
This is an empty json package

## Data Flow
JUST LIKE THIS
```
{/*a json package*/...},{/*a json package*/...},{/*a json package*/...},...,{/*FINISH signal*/}
```