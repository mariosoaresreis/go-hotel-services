# Introduction
This is an exercise to get a look at your coding abilities and knowledge of fields you will be working with if you decided to work with us. We will be looking at several factors including:

1. Code Quality
2. Use of Interfaces
3. Knowledge of Rules Engines
4. The ability to properly connect to an external API
5. Coding Patterns
6. Testability of Code
7. Usefulness of Error Messages Returned

Please add your work to a public repository and send up a link to it once you are ready for us to review it.

# Authentication

Use the following API-Key information to authenticate with the open-api server:

-Authentication URL:  [https://test.optii.io/oauth/authorize](https://test.optii.io/oauth/authorize)  
-Body (x-www-form-urlencoded)  
-client_secret:  oGDXMy5jsD4qEEjfWr1AOW1fuzKKZvQ9y1RcmhbYPNiesfi0haNYO7e9sL6FtlGr  
-client_id:  IXhXcAAfhY4WUIck8Sp1CgavF1zMVjVoihMoIvSuePQN6dpr  
-grant_type: client_credentials  
-scope: openapi  

Once you authenticate, you will receive an access token in response.  Use token on subsequent requests to the API by adding token as Authorization header.  Authorization Bearer **Access Token**

# Exercise

Using the swagger documentation found [here](https://app.swaggerhub.com/apis/luke.kennedy/optii-open_api/0.8.1#), We want you to to build a Microservice that can act as a mini rules engine to create certain jobs based on certain data being sent in the payload to your soon to be API.
Your api should accept a payload with the following json:

```
{
    "department": "<name of department>",
    "job_item": "<name of job_item>",
    "locations": [
        "<name of location 1>", "<name of location 2>"
    ]
}
```

`Note: If you want you can use the ids rather than the names for each of these fields, just leave a note on which path you took.`

The base URL to hit these endpoints will be [https://test.optii.io](https://test.optii.io).

The workflow engine should manage the following rules:
1. If the given department is not null and does not exist in the property then return a bad request
2. If the given job_item is not null and does not exist in the property then return a bad request
3. If the given locations are not null and do not exist in the property then return a bad request
4. If the given department is ‘Housekeeping’, the given job_item is one of ['Blanket','Sheets','Mattress'], and a location with a location type of 'Room' is given then create a job to clean the bed(s) in the given room
    * If a location with a location type of 'Floor' is given then create a job to clean the beds in all rooms with a location type of ‘Room’ on that floor.
    * If the given location is not of location type ‘Room’ or ‘Floor’ then return a bad request.
5. If the given department is ‘Engineering’ and a job_item and at least one location is given then create a job to repair the given job item at the given location(s).
    * If only one location is given and is of location type ‘Floor’ then create a job to repair the given job item in all locations on that floor.
6. If the department is ‘Room Service’ and a job item and at least 1 location is given then create a job to deliver that job item to the given locations.
    * If only one location is given and is of location type ‘Floor’ then create a job to deliver the given job item in all locations with a location type of 'Room' on that floor.
7. If none of the above criteria is met return a bad request.

`Note: Some of the above rules may require processing of data returned from one of the api calls rather than adding a filter in the api call itself, so please feel free to examine the structure of the data returned`