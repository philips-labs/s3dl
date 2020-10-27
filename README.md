# S3DL
S3DL generates pre-signed S3 Download URLS for any key in the 
attached S3 bucket. This is very convenient if you want to use
your S3 bucket as a public file store

# Usage
Deployment is optimized for Cloud foundry. Bind your S3 bucket to the app

| Path | Description |
|------|-------------|
| `/download?key=PATH` | Downloads the object at PATH |
| `/object/PATH` | Downloads the object at PATH |

# Config
| Environment | Description | Default |
|-------------|-------------|---------|
| `S3DL_EXPIRE` | Validity of pre-siging URL in minutes | 15  |


# Contact / Getting help
andy.lo-a-foe@philips.com

# license
License is MIT