# About
This is a program that calls retrieves stored data previously polled from [N2YO.com's REST API](https://www.n2yo.com/api/),
which holds position data on the International Space Station relative to the coordinates of Rice Creek Field Station
at SUNY Oswego.  They are pulled from the DynamoDB AWS database they were stored in via the corresponding 
[agent](https://github.com/Victor-Lockwood/CSC482-agent-vlockwoo) component.

## How To Deploy

1. Kick off the multistage build with `docker build -t vlockwoo-server .`
2. To test locally, run `docker run -p 40000:8080 -d vlockwoo-server` then in a 
browser, hit `localhost:40000/vlockwoo/status`.  You should get JSON with a 
200 response back.  Feel free to kill the local container.
3. In the same directory as your pem for authentication 
(mainly for convenience’s sake), run `docker save --output vlockwoo-server.tar server`
4. Run `scp -i <pem file name> vlockwoo-server.tar <ec2 username>@<ec2 IP>:`
5. SSH into the EC2 instance
6. Run `docker load --input vlockwoo-server.tar`
7. Run `docker run -e AWS_SECRET_ACCESS_KEY=<key> -e AWS_ACCESS_KEY_ID=<key> -e LOGGLY_TOKEN=<key> -p 40000:8080 -d vlockwoo-server`
8. Hit `http://<ip>:40000/vlockwoo/status` and if all is well you should have gotten a 200 response with JSON.
9. Clean up and remove your `.tar` file on the EC2 instance with `rm vlockwoo-server.tar`.
10. You can now exit the SSH session.


## Endpoints

### vlockwoo/status
Returns a JSON response with the table name and current record count.

#### URL Parameters
None

#### Example Response
```json
{
    "table": "vlockwoo-satellites",
    "recordCount": 84
}
```

### vlockwoo/all
Returns a JSON list of all entries in the table.

#### URL Parameters
None

#### Example Response
```json
[
    {
        "azimuth": {
            "B": null,
            "BOOL": null,
            "BS": null,
            "L": null,
            "M": null,
            "N": "68.04",
            "NS": null,
            "NULL": null,
            "S": null,
            "SS": null
        },
        "dec": {
            "B": null,
            "BOOL": null,
            "BS": null,
            "L": null,
            "M": null,
            "N": "12.32551325",
            "NS": null,
            "NULL": null,
            "S": null,
            "SS": null
        },
        "eclipsed": {
            "B": null,
            "BOOL": true,
            "BS": null,
            "L": null,
            "M": null,
            "N": null,
            "NS": null,
            "NULL": null,
            "S": null,
            "SS": null
        },
        "elevation": {
            "B": null,
            "BOOL": null,
            "BS": null,
            "L": null,
            "M": null,
            "N": "-5.58",
            "NS": null,
            "NULL": null,
            "S": null,
            "SS": null
        },
        "ra": {
            "B": null,
            "BOOL": null,
            "BS": null,
            "L": null,
            "M": null,
            "N": "107.54645057",
            "NS": null,
            "NULL": null,
            "S": null,
            "SS": null
        },
        "sataltitude": {
            "B": null,
            "BOOL": null,
            "BS": null,
            "L": null,
            "M": null,
            "N": "420.09",
            "NS": null,
            "NULL": null,
            "S": null,
            "SS": null
        },
        "satid": {
            "B": null,
            "BOOL": null,
            "BS": null,
            "L": null,
            "M": null,
            "N": "25544",
            "NS": null,
            "NULL": null,
            "S": null,
            "SS": null
        },
        "satlatitude": {
            "B": null,
            "BOOL": null,
            "BS": null,
            "L": null,
            "M": null,
            "N": "46.05890763",
            "NS": null,
            "NULL": null,
            "S": null,
            "SS": null
        },
        "satlongitude": {
            "B": null,
            "BOOL": null,
            "BS": null,
            "L": null,
            "M": null,
            "N": "-39.35875431",
            "NS": null,
            "NULL": null,
            "S": null,
            "SS": null
        },
        "satname": {
            "B": null,
            "BOOL": null,
            "BS": null,
            "L": null,
            "M": null,
            "N": null,
            "NS": null,
            "NULL": null,
            "S": "SPACE STATION",
            "SS": null
        },
        "timestamp": {
            "B": null,
            "BOOL": null,
            "BS": null,
            "L": null,
            "M": null,
            "N": "1699322119",
            "NS": null,
            "NULL": null,
            "S": null,
            "SS": null
        }
    }
]
```


### vlockwoo/search
Search by Timestamp or the Eclipsed flag.

#### URL Parameters
One or the other.
- `eclipsed`  - Either `true` or `false` 
- `timestamp` - A positive integer representing the timestamp, eg. `1699322119`

#### Example Response
```json
[
    {
        "azimuth": {
            "B": null,
            "BOOL": null,
            "BS": null,
            "L": null,
            "M": null,
            "N": "68.04",
            "NS": null,
            "NULL": null,
            "S": null,
            "SS": null
        },
        "dec": {
            "B": null,
            "BOOL": null,
            "BS": null,
            "L": null,
            "M": null,
            "N": "12.32551325",
            "NS": null,
            "NULL": null,
            "S": null,
            "SS": null
        },
        "eclipsed": {
            "B": null,
            "BOOL": true,
            "BS": null,
            "L": null,
            "M": null,
            "N": null,
            "NS": null,
            "NULL": null,
            "S": null,
            "SS": null
        },
        "elevation": {
            "B": null,
            "BOOL": null,
            "BS": null,
            "L": null,
            "M": null,
            "N": "-5.58",
            "NS": null,
            "NULL": null,
            "S": null,
            "SS": null
        },
        "ra": {
            "B": null,
            "BOOL": null,
            "BS": null,
            "L": null,
            "M": null,
            "N": "107.54645057",
            "NS": null,
            "NULL": null,
            "S": null,
            "SS": null
        },
        "sataltitude": {
            "B": null,
            "BOOL": null,
            "BS": null,
            "L": null,
            "M": null,
            "N": "420.09",
            "NS": null,
            "NULL": null,
            "S": null,
            "SS": null
        },
        "satid": {
            "B": null,
            "BOOL": null,
            "BS": null,
            "L": null,
            "M": null,
            "N": "25544",
            "NS": null,
            "NULL": null,
            "S": null,
            "SS": null
        },
        "satlatitude": {
            "B": null,
            "BOOL": null,
            "BS": null,
            "L": null,
            "M": null,
            "N": "46.05890763",
            "NS": null,
            "NULL": null,
            "S": null,
            "SS": null
        },
        "satlongitude": {
            "B": null,
            "BOOL": null,
            "BS": null,
            "L": null,
            "M": null,
            "N": "-39.35875431",
            "NS": null,
            "NULL": null,
            "S": null,
            "SS": null
        },
        "satname": {
            "B": null,
            "BOOL": null,
            "BS": null,
            "L": null,
            "M": null,
            "N": null,
            "NS": null,
            "NULL": null,
            "S": "SPACE STATION",
            "SS": null
        },
        "timestamp": {
            "B": null,
            "BOOL": null,
            "BS": null,
            "L": null,
            "M": null,
            "N": "1699322119",
            "NS": null,
            "NULL": null,
            "S": null,
            "SS": null
        }
    }
]
```