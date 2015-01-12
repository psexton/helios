helios
======

Backup and restore a couchdb npm registry to s3.

Setup
-----

You'll need a JSON config file that looks something like this:
```
{
    "AWS": {
        "AccessKeyID": "XXXXXXXXXXXXXXXXXXXX",
        "SecretAccessKey": "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
        "S3BucketName": "your-bucket-name"
    },
    "Couch": {
        "Username": "jrandom",
        "Password": "secret",
        "URL": "http://localhost:5984/"
    },
    "DaemonPause": "60s"
}
```

Where the AWS section corresponds to a keypair and S3 bucket, and the Couch section corresponds to an admin user. Most reasonable durations can be used for `DaemonPause` (e.g. 15s, 5m, 1h).

Running
-------

`helios --conf ~/helios.json --sunrise`

`helios --conf ~/helios.json --sunset`

`helios --conf ~/helios.json --daemon`
