# Description

- This application allows the user to process images collected from stores.
- Users can submit new jobs to the system, with each job potentially containing multiple images from various stores.
- Given that the number of images to be processed per job may be quite large, the actual job processing occurs asynchronously in the background.
- Users receive a job ID when a job is successfully created, which they can use to query the job status at any time.
- Additionally, users need to provide a list of valid Store IDs in the form of a CSV file.

# Endpoints

## 1. Submit New Job

- **Endpoint:** `/api/submit`
- **Method:** `POST`
- **Description:** Creates a background job to process the images collected from stores.

### Request Payload

```json
{
   "count":2,
   "visits":[
      {
         "store_id":"S00339218",
         "image_url":[
            "https://www.gstatic.com/webp/gallery/2.jpg",
            "https://www.gstatic.com/webp/gallery/3.jpg"
         ],
         "visit_time": "time of store visit"
      },
      {
         "store_id":"S01408764",
         "image_url":[
            "https://www.gstatic.com/webp/gallery/3.jpg"
         ],
         "visit_time": "time of store visit"
      }
   ]
}
```

### Success Response
- **Condition:** If everything is OK, and a job is created.
- **Status Code:** `201 CREATED`
- **Content:**
```json
{
  "job_id": "6763013d5e3caadda500e896"
}
```

### Error Response
- **Condition:** If can not decode JSON body or fields are missing OR count != len(visits) or basic data validation fails
- **Status Code:** `400 BAD REQUEST`
- **Content:**
```json
{
  "error": "image_url cannot be empty for store_id: RP00006"
}
```

## 2. Get Job Status Info
- **Endpoint:** `/api/status?jobid=6763013d5e3caadda500e896`
- **URL Parameters:** `jobid` Job ID received while creating the job.
- **Method:** `GET`
- **Description:** Fetches the current status of the job with the given Job ID.

### Success Response
- **Condition:** If everything is OK, and Job ID exists.
- **Status Code:** `200 OK`
- **Content:**

#### Job Status: completed/ongoing
```json
{
  "status": "completed",
  "job_id": "6738ddca9ed022cf4933f9d1"
}
```

#### Job Status: failed
If a `store_id` does not exist or an image download fails for any given URL.
```json
{
  "status": "failed",
  "job_id": "6738d31e1f67c7e7f5f70e2c",
  "error": {
    "store_id": "RP00006",
    "error": "failed to download image: Get \"https://www.gstatdic.com/webp/gallery/2.jpg\": dial tcp: lookup www.gstatdic.com on 192.168.1.1:53: no such host"
  }
}
```

### Error Response
- **Condition:** If Job ID is missing or does not exist in the system.
- **Status Code:** `400 BAD REQUEST`
- **Content:**
```json
{
  "error": "jobid does not exist"
}
```

# Assumptions
- The CSV containing the list of Store IDs has the first row as the header, and the Store IDs are located in the third column (1-based indexing).

# Installation Instructions

Unzip the given zip file and `cd` into it.
Make sure you have `docker` and `docker-compose` installed before proceeding.



- The default CSV file containing Store IDs is `StoreMasterAssignment.csv`, located in the store folder.I have given the filePath of StoreMasterAssignment.csv in .env file

**Note:** Skip to "Docker Compose" subsection for a single command install and run.

**Note**: This project relies on MongoDB as the database. The program expects the environment variable MONGODB_URL to be set to point to the URI of the MongoDB database (either the managed Atlas Cluster or a self-hosted server).

## Setting up MongoDB Docker Container
In this section, we will set up the MongoDB Docker image and run it to provide a database server. You can skip this section if you are using MongoDB Atlas; in that case, you will need to obtain the URI from the MongoDB Cluster information page.

[Install Docker Engine](https://docs.docker.com/engine/install/) before proceeding.

Pull the MongoDB Docker Image \
`docker pull mongodb/mongodb-community-server:latest`

Create a Docker Network
`docker network create kirana_club`

Run the Image as a Container (attaching to `kirana_club` network) \
`docker run --name mongodb -p 27017:27017 -d --network kirana_club mongodb/mongodb-community-server:latest`

You can check if the MongoDB container is running by executing: \
`docker ps`

**Note:** If you ever get the following error:
> docker: Error response from daemon: Conflict. The container name "/mongodb" is already in use by container "< HASH >". You have to remove (or rename) that container to be able to reuse that name.

Then run the following commands to remove the `mongodb` container:

`docker stop mongodb` \
`docker rm mongodb`

## Without Docker

**Note:** This project was built and tested with Go version 1.23.1. Please ensure that this version is installed on your machine before proceeding further. Please refer to [Manage Go Installations](https://go.dev/doc/manage-install)

Navigate to the project root.

### Setting Env Variable

Before running any of the following methods, set the `MONGODB_URI` env variable using the following command (Linux):  
`export MONGODB_URI=<value>` \
Replace the `<value>` with the URI in your specific deployment. 

If you followed the MongoDB Docker setup above then use `export MONGODB_URI="mongodb://localhost:27017"              `

### Building the binary
Run `go build ./cmd/main.go` to build the binary, optionally specifying the generated binary's name using the `-o` flag like `go build -o kirana_club ./cmd/main.go`

The binary will be generated in the current directory.

Run the binary using `./main` or `./kirana_club` if you used a custom name while building.

### Running without building

If you want to compile the file, run it and then remove the binary after the program ends using a single command then use `go run ./cmd/main.go`



## With Docker

A Dockerfile is present in the root directory, which can be used to create the Docker image.

Navigate to the root directory of the project.

Use the following command to build the docker image with the tag `kirana_club`

`docker build -t kirana_club .`

Run the container with the following command, replacing `<YOUR_URI>` with the appropriate URI:

`docker run -p 8080:8080 --network kirana_club -e MONGODB_URI=<YOUR_URI> -v $(pwd)/docker_mounts/files:/app/files -v $(pwd)/docker_mounts/logs:/app/logs kirana_club`

If you followed the MongoDB Docker setup above, then use:

`docker run -p 8080:8080 --network kirana_club -e MONGODB_URI="mongodb://mongodb:27017" -v $(pwd)/docker_mounts/files:/app/files -v $(pwd)/docker_mounts/logs:/app/logs  kirana_club`

All the images downloaded by the application will be saved in `./docker_mounts/files/`.

## Docker Compose

Navigate to the project root. It contains a `docker-compose.yml` file which can be used to setup, build and run both MongoDB and image-job-processor containers using only a single command.

If logging to `stdout`, then use the following command to run the application and receive the logs on the terminal:

`docker-compose up`

To stop the application, press `Ctrl+C` and wait for graceful shutdown.


# Development Environment

### Text Editor/IDE
- **Text Editor/IDE:** VSCode

### Libraries and Frameworks
- **Programming Language:** Go v1.23.0
- **Libraries/Frameworks Used:**
    - github.com/gin-gonic/gin v1.10.0
    - github.com/google/uuid v1.6.0
    - github.com/joho/godotenv v1.5.1
    - github.com/google/uuid v1.6.0
    - go.mongodb.org/mongo-driver v1.17.1
    - github.com/klauspost/compress v1.13.6
    - github.com/montanaflynn/stats v0.7.1
    - github.com/xdg-go/pbkdf2 v1.0.0
    - github.com/xdg-go/scram v1.1.2
    - github.com/xdg-go/stringprep v1.0.4
    - github.com/youmark/pkcs8 v0.0.0-20240726163527-a2c0da244d78
    - github.com/yuin/goldmark v1.4.13
    - go.mongodb.org/mongo-driver v1.17.1
    - golang.org/x/crypto v0.26.0
    - golang.org/x/mod v0.17.0
    - golang.org/x/net v0.21.0
    - golang.org/x/sync v0.8.0
    - golang.org/x/sys v0.23.0
    - golang.org/x/term v0.23.0
    - golang.org/x/text v0.17.0
    - golang.org/x/tools v0.21.1-0.20240508182429-e35e4ccd0d2d
    - golang.org/x/xerrors v0.0.0-20190717185122-a985d3407aa7

### Additional Tools
- **Version Control:** Git v2.46.1
- **Containerization:** Docker v27.2.1

# Future Improvement Scope
- The `ProcessJob` function downloads and processes only those images that have not yet been processed. This feature can be utilized to resume jobs that were still "ongoing" during a previous run if the application crashed. This resumption capability can be implemented at the start of the program.
- Currently, we have only a single instance of the server running, which could easily become overwhelmed in the event of very high loads. Our system should be able to dynamically scale the number of server instances to better manage the workload.
- At present, a single mistyped Store ID or image URL causes the entire job to be marked as failed. We should provide the user with more specific feedback and allow them to make corrections. In that case, the system should process only the corrected fields. 