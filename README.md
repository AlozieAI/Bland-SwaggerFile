# Bland-SwaggerFile


**Swagger Implementation for a Go API using Gin**

This README provides a comprehensive overview of setting up a Go-based API with Swagger documentation, utilizing the Gin framework. The API manages various functionalities, including sending calls, analyzing calls, managing folders and pathways, and integrating chat features. This implementation uses the swaggo library to generate and serve the Swagger documentation.

**Table of Contents**

Prerequisites

Project Structure

Setup and Installation

Running the API

Swagger Documentation

API Endpoints

Models

Notes

**Prerequisites**

Before running the API, ensure you have the following installed:

Go (version 1.16 or later)

Gin framework

Swaggo for Swagger integration

**Project Structure**

Here's an overview of the key files in this project:

main.go: Initializes the Gin server, sets up API routes, and serves the Swagger documentation.
controller: Contains the functions that handle API requests.

model: Defines the request and response data structures.

docs: Contains the Swagger documentation files.

**Setup and Installation**

Clone the Repository

bash
```
git clone <repository-url>
cd <repository-folder>
```

Install Dependencies

bash
```
go mod tidy
```

Install Swaggo for Swagger Documentation

bash
```
go install github.com/swaggo/swag/cmd/swag@latest
```

Generate Swagger Documentation

bash
```
swag init
```

This command generates Swagger documentation files in the docs folder.

**Running the API**

To start the API server, run:

bash
```
go run main.go
```

The API will be available at http://0.0.0.0:8080.

**Swagger Documentation**

The Swagger UI can be accessed at:

arduino

```
http://0.0.0.0:8080/swagger/index.html
```
Swagger Annotations

@title: The title of the API.
@version: Version of the API.
@description: A brief description of the API.
@termsOfService: Terms of service URL.
@contact: Contact information for API support.
@license: License information.
@host: The host address for the API.
@BasePath: The base path for the API endpoints.
@securityDefinitions: Defines the security scheme for API requests.

**API Endpoints**

Call Management

Send a Call

POST /api/v1/call

Sends a call using pathways by providing a phone number and pathway ID.


Analyze a Call

POST /api/v1/call/:call_id/analyze

Analyzes a call with AI by providing a call ID, goal, and questions.


Get Call Details

GET /api/v1/calls/:call_id

Retrieves detailed information, metadata, and transcripts for a call.


Folder and Pathway Management

Create Folder

POST /api/v1/folders

Creates a new folder for the authenticated user.


Create and Move Pathway

POST /api/v1/pathways/create-and-move

Creates a new conversational pathway and moves it to a folder.


Get Pathway Information

GET /api/v1/convo_pathway/:pathway_id

Returns detailed information about a specific pathway.


Update Pathway

POST /api/v1/pathway/update/:pathway_id
Updates a conversational pathway’s fields.


Delete Pathway

DELETE /api/v1/delete/convo_pathway/:pathway_id
Deletes a specific conversational pathway.


Chat Management
Create Chat

POST /api/v1/pathways/chat/create
Creates a chat instance for a pathway.
Send Message to Chat

POST /api/v1/pathways/chat/:chat_id/send
Sends a message to a specific pathway chat and receives a response.


**Models**

The model package defines the data structures used for API requests and responses. Some key models include:

ErrorResponse: Defines the structure for error responses.

SendCall: Request structure for sending a call.

CallResponse: Response structure for call status.

AnalyzeCallRequest: Request structure for analyzing a call.

AnalyzeCallResponse: Response structure for analysis results.

CallDetail: Structure for call details and transcripts.

CreateFolderRequest/Response: Structures for folder creation.

CreatePathwayRequest/Response: Structures for pathway management.

SendMessageRequest/Response: Structures for chat messages.

**Notes**

Security: The API uses bearer token authentication. Ensure you include the Authorization header with your requests.

External API Calls: The API interacts with external services, so proper error handling and logging are essential.

Swagger Integration: The Swagger documentation helps visualize and interact with the API easily. Make sure to update the Swagger comments if any changes are made to the API.
