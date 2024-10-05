package controller

import (
	"bland/model"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
  "bytes"
	"github.com/gin-gonic/gin"
)

// SendCall godoc
// @Summary      Send call using Pathways
// @Description  Send call using Pathways by providing a phone number and pathway ID
// @Tags         SendCall
// @Accept       json
// @Produce      json
// @Param        request       body      model.SendCall  true  "Request body"
// @Success      200  {object}  model.CallResponse  "Success"
// @Failure      400  {object}  model.ErrorResponse  "Bad Request"
// @Failure      401  {object}  model.ErrorResponse  "Unauthorized - Bearer token required"
// @Failure      500  {object}  model.ErrorResponse  "Internal Server Error"
// @Security     bearerToken
// @Router       /call [post]
func SendCall(c *gin.Context) {

	// Step 1: Bind the JSON request body to the SendCall struct
	var requestData model.SendCall
	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Message: err.Error()})
		return
	}

	// Step 2: Extract the bearer token from the request header
	bearerToken := c.GetHeader("Authorization")
	if bearerToken == "" {
		c.JSON(http.StatusUnauthorized, model.ErrorResponse{Message: "Authorization token is required"})
		return
	}

	// Step 3: Prepare the payload for the external API
	payload := strings.NewReader(fmt.Sprintf(`{
		"phone_number": "%s",
		"pathway_id": "%s"
	}`, requestData.PhoneNumber, requestData.PathwayID))

	// Step 4: Create the POST request
	url := "https://api.bland.ai/v1/calls"
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		log.Printf("Error creating request: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Request creation failed"})
		return
	}

	// Step 5: Set headers, including the Authorization token
	req.Header.Add("Authorization", bearerToken) // Pass the extracted token
	req.Header.Add("Content-Type", "application/json")

	// Step 6: Send the request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Error making the request: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Request failed"})
		return
	}
	defer res.Body.Close()

	// Step 7: Read the response body
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("Error reading the response: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
		return
	}
	// Step 8: Unmarshal the response into the CallResponse struct
	var callResponse model.CallResponse
	if err := json.Unmarshal(body, &callResponse); err != nil {
		log.Printf("Error unmarshalling the response: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse response"})
		return
	}

	// Step 9: Return the external API's response in the expected format
	c.JSON(http.StatusOK, callResponse)
}

// AnalyzeCall godoc
// @Summary      Analyze a call with AI
// @Description  Analyze a call by providing the call ID, goal, and an array of questions
// @Tags         AnalyzeCall
// @Accept       json
// @Produce      json
// @Param        call_id      path      string              true   "Call ID"
// @Param        request      body      model.AnalyzeCallRequest   true   "Request body"
// @Success      200  {object}  model.AnalyzeCallResponse  "Success"
// @Failure      400  {object}  model.ErrorResponse        "Bad Request"
// @Failure      401  {object}  model.ErrorResponse        "Unauthorized - Bearer token required"
// @Failure      500  {object}  model.ErrorResponse        "Internal Server Error"
// @Security     bearerToken
// @Router       /call/{call_id}/analyze [post]
// AnalyzeCall analyzes a call using the provided call_id from the URL, goal, and questions from the request body
func AnalyzeCall(c *gin.Context) {
	// Step 1: Extract the call_id from the URL path dynamically based on user input
	callID := c.Param("call_id")

	// Log the callID to ensure it's being captured properly
	log.Printf("callID: %s", callID)

	// Step 2: Bind the request JSON to the AnalyzeCallRequest struct
	var requestBody model.AnalyzeCallRequest
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		log.Printf("Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Log the request body as JSON for debugging purposes
	requestBodyJSON, err := json.Marshal(requestBody)
	if err != nil {
		log.Printf("Error marshaling request body: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal request body"})
		return
	}
	log.Printf("Request Body (JSON): %s", string(requestBodyJSON))

	// Step 3: Extract the bearer token from the request header
	bearerToken := c.GetHeader("Authorization")
	if bearerToken == "" {
		log.Printf("Missing Authorization token")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token is required"})
		return
	}

	// Step 4: Send the request to the external API
	url := fmt.Sprintf("https://api.bland.ai/v1/calls/%s/analyze", callID)

	log.Printf("Calling URL: %s", url)
	req, err := http.NewRequest("POST", url, strings.NewReader(string(requestBodyJSON)))
	if err != nil {
		log.Printf("Error creating request: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Request creation failed"})
		return
	}

	req.Header.Add("Authorization", bearerToken)
	req.Header.Add("Content-Type", "application/json")

	// Step 5: Send the request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Error making the request: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Request failed"})
		return
	}
	defer res.Body.Close()

	// Log the response status code
	log.Printf("Response Status Code: %d", res.StatusCode)

	// Step 6: Read the response body
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("Error reading the response: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
		return
	}

	// Log the raw response body for debugging purposes
	log.Printf("Response Body: %s", string(body))

	// If the status code is not 200, log the response and return
	if res.StatusCode != http.StatusOK {
		log.Printf("Unexpected status code %d: %s", res.StatusCode, string(body))
		c.JSON(res.StatusCode, gin.H{"error": "Unexpected error occurred", "response": string(body)})
		return
	}

	// Step 7: Unmarshal the response into the AnalyzeCallResponse struct
	var analyzeResponse model.AnalyzeCallResponse
	if err := json.Unmarshal(body, &analyzeResponse); err != nil {
		log.Printf("Error unmarshalling the response: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse response"})
		return
	}

	// Step 8: Return the external API's response
	c.JSON(http.StatusOK, analyzeResponse)
}


// GetCallDetails godoc
// @Summary      Get call details
// @Description  Retrieve detailed information, metadata, and transcripts for a call
// @Tags         CallDetails
// @Accept       json
// @Produce      json
// @Param        call_id  path  string  true  "Call ID"
// @Success      200  {object}  model.CallDetail  "Call details retrieved successfully"
// @Failure      400  {object}  model.ErrorResponse  "Invalid input"
// @Failure      500  {object}  model.ErrorResponse  "Internal server error"
// @Security     bearerToken
// @Router       /calls/{call_id} [get]
// GetCallDetails retrieves detailed information about a specific call
func GetCallDetails(c *gin.Context) {
	// Step 1: Extract the call_id from the path
	callID := c.Param("call_id")

	// Step 2: Extract the bearer token from the request header
	bearerToken := c.GetHeader("Authorization")
	if bearerToken == "" {
		log.Printf("Missing Authorization token")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token is required"})
		return
	}
	// Step 3: Prepare the URL with the dynamic call_id
	url := fmt.Sprintf("https://api.bland.ai/v1/calls/%s", callID)

	// Step 4: Create a new GET request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("Error creating request: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}

	// Step 5: Add the Authorization header (replace <API_KEY> with actual API key)
	req.Header.Add("Authorization", bearerToken)

	// Step 6: Execute the request
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Printf("Error making the request: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to make request"})
		return
	}
	defer res.Body.Close()

	// Step 7: Read the response body
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response body"})
		return
	}

	// Step 8: Check if the status code is not 200
	if res.StatusCode != http.StatusOK {
		log.Printf("Received non-200 status code: %d", res.StatusCode)
		c.JSON(res.StatusCode, gin.H{"error": "Received non-200 response", "response": string(body)})
		return
	}

	// Step 9: Unmarshal the response into the CallDetail struct
	var callDetail model.CallDetail
	if err := json.Unmarshal(body, &callDetail); err != nil {
		log.Printf("Error unmarshalling response: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse response"})
		return
	}

	// Step 10: Return the call details as a JSON response
	c.JSON(http.StatusOK, callDetail)
}



// CreateFolder godoc
// @Summary      Create a new folder
// @Description  Creates a new folder for the authenticated user
// @Tags         Folder
// @Accept       json
// @Produce      json
// @Param        request body model.CreateFolderRequest true "Request body for creating folder"
// @Success      200  {object}  model.CreateFolderResponse  "Folder created successfully"
// @Failure      400  {object} model.ErrorResponse  "Invalid input"
// @Failure      500  {object}  model.ErrorResponse  "Internal server error"
// @Security     bearerToken
// @Router       /folders [post]
// CreateFolder creates a new folder for the authenticated user
func CreateFolder(c *gin.Context) {
	// Step 1: Bind the request JSON to the CreateFolderRequest struct
	var requestBody model.CreateFolderRequest
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		log.Printf("Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Log the request body as JSON for debugging purposes
	requestBodyJSON, err := json.Marshal(requestBody)
	if err != nil {
		log.Printf("Error marshaling request body: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal request body"})
		return
	}
	log.Printf("Request Body (JSON): %s", string(requestBodyJSON))

	// Step 2: Extract the bearer token from the request header
	bearerToken := c.GetHeader("Authorization")
	if bearerToken == "" {
		log.Printf("Missing Authorization token")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token is required"})
		return
	}

	// Step 3: Prepare the URL for the folder creation endpoint
	url := "https://us.api.bland.ai/v1/pathway/folders"
	log.Printf("Calling URL: %s", url)

	// Step 4: Create a new POST request
	req, err := http.NewRequest("POST", url, strings.NewReader(string(requestBodyJSON)))
	if err != nil {
		log.Printf("Error creating request: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Request creation failed"})
		return
	}

	// Step 5: Set headers for the request
	req.Header.Add("Authorization", bearerToken)
	req.Header.Add("Content-Type", "application/json")

	// Step 6: Send the request to the external API
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Error making the request: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Request failed"})
		return
	}
	defer res.Body.Close()

	// Log the response status code
	log.Printf("Response Status Code: %d", res.StatusCode)

	// Step 7: Read the response body
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("Error reading the response: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
		return
	}

	// Log the raw response body for debugging purposes
	log.Printf("Response Body: %s", string(body))

	// If the status code is not 200, log the response and return
	if res.StatusCode != http.StatusOK {
		log.Printf("Unexpected status code %d: %s", res.StatusCode, string(body))
		c.JSON(res.StatusCode, gin.H{"error": "Unexpected error occurred", "response": string(body)})
		return
	}

	// Step 8: Unmarshal the response into the CreateFolderResponse struct
	var folderResponse model.CreateFolderResponse
	if err := json.Unmarshal(body, &folderResponse); err != nil {
		log.Printf("Error unmarshalling the response: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse response"})
		return
	}

	// Step 9: Return the folder details as a JSON response (return only the "data" part)
	c.JSON(http.StatusOK, folderResponse.Data)
}




// CreateAndMovePathway godoc
// @Summary      Create and move pathway
// @Description  Creates a new conversational pathway and moves it to a folder
// @Tags         Pathway
// @Accept       json
// @Produce      json
// @Param        request body model.CreatePathwayRequest true "Request body for creating pathway"
// @Param        folder_id query string false "Folder ID to move the pathway into"
// @Success      200  {object}  model.CombinedResponse  "Combined response of creating and moving pathway"
// @Failure      400  {object}   model.ErrorResponse  "Invalid input"
// @Failure      500  {object}   model.ErrorResponse  "Internal server error"
// @Security     bearerToken
// @Router       /pathways/create-and-move [post]
// CreateAndMovePathway creates a new conversational pathway and moves it to a folder
// CreateAndMovePathway creates a new conversational pathway and moves it to a folder
// CreateAndMovePathway creates a new conversational pathway and moves it to a folder
func CreateAndMovePathway(c *gin.Context) {
	// Step 1: Bind the request JSON for creating a pathway
	var createRequest model.CreatePathwayRequest
	if err := c.ShouldBindJSON(&createRequest); err != nil {
		log.Printf("Error binding JSON for CreatePathwayRequest: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Step 2: Extract the bearer token from the request header
	bearerToken := c.GetHeader("Authorization")
	if bearerToken == "" {
		log.Printf("Missing Authorization token")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token is required"})
		return
	}

	// Step 3: Create the pathway (first API call)
	createPathwayURL := "https://api.bland.ai/v1/convo_pathway/create"
	createPathwayBody, _ := json.Marshal(createRequest)

	req, err := http.NewRequest("POST", createPathwayURL, bytes.NewBuffer(createPathwayBody))
	if err != nil {
		log.Printf("Error creating request: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}

	req.Header.Add("Authorization", bearerToken)
	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Error making the request: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Request failed"})
		return
	}
	defer res.Body.Close()

	// Read and log the response from the first API call
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("Error reading the response: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
		return
	}

	var createPathwayResponse model.CreatePathwayResponse
	if err := json.Unmarshal(body, &createPathwayResponse); err != nil {
		log.Printf("Error unmarshalling the response for CreatePathway: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse CreatePathway response"})
		return
	}

	// Log the response of creating pathway
	log.Printf("CreatePathwayResponse: Status=%s, PathwayID=%s", createPathwayResponse.Status, createPathwayResponse.PathwayID)

	// Check if pathway creation was successful
	if createPathwayResponse.Status != "success" {
		log.Printf("Pathway creation failed: %v", string(body))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Pathway creation failed", "response": string(body)})
		return
	}

	// Step 4: Move the pathway (second API call)
	var moveRequest model.MovePathwayRequest
	moveRequest.PathwayID = createPathwayResponse.PathwayID // Use the pathway ID from the first response

	// Optional: Add folder ID if provided in the request
	folderID := c.Query("folder_id") // assuming folder_id is passed as a query param
	if folderID != "" {
		moveRequest.FolderID = folderID
	}

	// Log the move request before sending
	log.Printf("MovePathwayRequest: PathwayID=%s, FolderID=%s", moveRequest.PathwayID, moveRequest.FolderID)

	// Prepare the API request to move the pathway
	movePathwayURL := "https://us.api.bland.ai/v1/pathway/folders/move"
	movePathwayBody, _ := json.Marshal(moveRequest)

	req, err = http.NewRequest("POST", movePathwayURL, bytes.NewBuffer(movePathwayBody))
	if err != nil {
		log.Printf("Error creating request for moving pathway: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create move request"})
		return
	}

	req.Header.Add("Authorization", bearerToken)
	req.Header.Add("Content-Type", "application/json")

	res, err = http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Error making the move request: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Request failed"})
		return
	}
	defer res.Body.Close()

	// Read and log the response from the second API call
	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("Error reading the move response: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read move response"})
		return
	}

	var movePathwayResponse model.MovePathwayResponse
	if err := json.Unmarshal(body, &movePathwayResponse); err != nil {
		log.Printf("Error unmarshalling the response for MovePathway: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse MovePathway response"})
		return
	}

	// Log the response of moving pathway
	oldFolderID := "nil"
	newFolderID := "nil"
	if movePathwayResponse.Data.OldFolderID != nil {
		oldFolderID = *movePathwayResponse.Data.OldFolderID
	}
	if movePathwayResponse.Data.NewFolderID != nil {
		newFolderID = *movePathwayResponse.Data.NewFolderID
	}
	log.Printf("MovePathwayResponse: PathwayID=%s, OldFolderID=%s, NewFolderID=%s", 
	    movePathwayResponse.Data.PathwayID, oldFolderID, newFolderID)

	// Step 5: Combine the responses and return
	combinedResponse := model.CombinedResponse{
		CreatePathwayResponse: createPathwayResponse,
		MovePathwayData:       movePathwayResponse.Data,  // Use MovePathwayData from Data field
	}

	// Log the combined response
	combinedResponseJSON, _ := json.Marshal(combinedResponse)
	log.Printf("CombinedResponse (JSON): %s", string(combinedResponseJSON))

	c.JSON(http.StatusOK, combinedResponse)
}

// CreateChat godoc
// @Summary      Create a pathway chat
// @Description  Creates a chat instance for a pathway, which can be used to send and receive messages.
// @Tags         Chat
// @Accept       json
// @Produce      json
// @Param        request body model.CreateChatRequest true "Request body for creating chat"
// @Success      200  {object}  model.CreateChatResponse  "Chat instance created successfully"
// @Failure      400  {object}   model.ErrorResponse  "Invalid input"
// @Failure      500  {object}   model.ErrorResponse  "Internal server error"
// @Security     bearerToken
// @Router       /pathways/chat/create [post]
func CreateChat(c *gin.Context) {
	// Step 1: Bind the request body to CreateChatRequest struct
	var createChatRequest model.CreateChatRequest
	if err := c.ShouldBindJSON(&createChatRequest); err != nil {
		log.Printf("Error binding JSON for CreateChatRequest: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Log the values of the request model
	log.Printf("CreateChatRequest: PathwayID=%s, StartNodeID=%s", createChatRequest.PathwayID, createChatRequest.StartNodeID)

	// Step 2: Extract the bearer token from the request header
	bearerToken := c.GetHeader("Authorization")
	if bearerToken == "" {
		log.Printf("Missing Authorization token")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token is required"})
		return
	}

	// Step 3: Marshal the request body
	requestBodyJSON, err := json.Marshal(createChatRequest)
	if err != nil {
		log.Printf("Error marshaling request body: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal request body"})
		return
	}

	// Log the actual JSON request being sent
	log.Printf("Request Body (JSON): %s", string(requestBodyJSON))

	// Step 4: Prepare the API request to create a chat
	url := "https://us.api.bland.ai/v1/pathway/chat/create"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBodyJSON))
	if err != nil {
		log.Printf("Error creating request: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}

	// Set headers
	req.Header.Add("Authorization", bearerToken)
	req.Header.Add("Content-Type", "application/json")

	// Step 5: Execute the request
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Printf("Error making the request: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Request failed"})
		return
	}
	defer res.Body.Close()

	// Step 6: Read and log the response body
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
		return
	}

	// Log the actual response from the API
	log.Printf("API Response Body: %s", string(body))

	// Step 7: Unmarshal the response into CreateChatResponse struct
	var createChatResponse model.CreateChatResponse
	if err := json.Unmarshal(body, &createChatResponse); err != nil {
		log.Printf("Error unmarshalling the response: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse response"})
		return
	}

	// Log the unmarshalled response values (wrapped in the 'data' field)
	log.Printf("CreateChatResponse: ChatID=%s, Message=%s", createChatResponse.Data.ChatID, createChatResponse.Data.Message)

	// Step 8: Return the chat creation response
	c.JSON(http.StatusOK, createChatResponse)
}

// GetPathwayInfo godoc
// @Summary      Get pathway information
// @Description  Returns detailed information about a specific pathway, including nodes and edges.
// @Tags         Pathway
// @Accept       json
// @Produce      json
// @Param        pathway_id  path      string  true  "The pathway ID"
// @Success      200  {object}  model.GetPathwayResponse  "Pathway information retrieved successfully"
// @Failure      400  {object} model.ErrorResponse  "Invalid input"
// @Failure      500  {object}  model.ErrorResponse  "Internal server error"
// @Security     bearerToken
// @Router       /convo_pathway/{pathway_id} [get]
func GetPathwayInfo(c *gin.Context) {
	// Step 1: Get the pathway_id from the URL path
	pathwayID := c.Param("pathway_id")

	// Step 2: Extract the bearer token from the request header
	bearerToken := c.GetHeader("Authorization")
	if bearerToken == "" {
		log.Printf("Missing Authorization token")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token is required"})
		return
	}

	// Step 3: Create the API request to get pathway information
	url := fmt.Sprintf("https://api.bland.ai/v1/convo_pathway/%s", pathwayID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("Error creating request: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}

	// Set the Authorization header
	req.Header.Add("Authorization", bearerToken)

	// Step 4: Execute the request
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Printf("Error making the request: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Request failed"})
		return
	}
	defer res.Body.Close()

	// Step 5: Read the response body
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("Error reading the response body: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
		return
	}

	// Log the response from the API
	log.Printf("API Response: %s", string(body))

	// Step 6: Unmarshal the response into GetPathwayResponse struct
	var pathwayResponse model.GetPathwayResponse
	if err := json.Unmarshal(body, &pathwayResponse); err != nil {
		log.Printf("Error unmarshalling the response: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse response"})
		return
	}

	// Step 7: Return the pathway information as JSON
	c.JSON(http.StatusOK, pathwayResponse)
}

// UpdatePathway godoc
// @Summary      Update conversational pathway
// @Description  Updates a conversational pathway’s fields including name, description, nodes, and edges
// @Tags         Pathway
// @Accept       json
// @Produce      json
// @Param        pathway_id path string true "Pathway ID to update"
// @Param        request body model.UpdatePathwayRequest true "Request body for updating the pathway"
// @Success      200  {object}  model.PathwayData  "Pathway updated successfully"
// @Failure      400  {object}  model.ErrorResponse  "Invalid input"
// @Failure      500  {object}  model.ErrorResponse  "Internal server error"
// @Security     bearerToken
// @Router       /pathway/update/{pathway_id} [post]
func UpdatePathway(c *gin.Context) {
    pathwayID := c.Param("pathway_id")
    log.Printf("Received request to update pathway. Pathway ID: %s", pathwayID)

    var updateRequest model.UpdatePathwayRequest
    if err := c.ShouldBindJSON(&updateRequest); err != nil {
        log.Printf("Error binding JSON: %v", err)
        c.JSON(http.StatusBadRequest, model.ErrorResponse{Message: err.Error()})
        return
    }

    // Log the model data
    log.Printf("Update Request Model: %+v", updateRequest)

    bearerToken := c.GetHeader("Authorization")
    if bearerToken == "" {
        log.Printf("Authorization token missing")
        c.JSON(http.StatusUnauthorized, model.ErrorResponse{Message: "Authorization token is required"})
        return
    }
    log.Printf("Authorization token received")

    requestBodyJSON, err := json.Marshal(updateRequest)
    if err != nil {
        log.Printf("Error marshaling request body: %v", err)
        c.JSON(http.StatusInternalServerError, model.ErrorResponse{Message: "Failed to marshal request body"})
        return
    }

    apiURL := fmt.Sprintf("https://api.bland.ai/v1/convo_pathway/%s", pathwayID)
    log.Printf("API URL: %s", apiURL)

    req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(requestBodyJSON))
    if err != nil {
        log.Printf("Error creating request: %v", err)
        c.JSON(http.StatusInternalServerError, model.ErrorResponse{Message: "Failed to create request"})
        return
    }

    // Add headers
    req.Header.Add("Authorization", bearerToken)
    req.Header.Add("Content-Type", "application/json")
    log.Printf("Request Headers: %v", req.Header)

    // Log request body
    log.Printf("Request Body: %s", string(requestBodyJSON))

    client := &http.Client{}
    res, err := client.Do(req)
    if err != nil {
        log.Printf("Error sending request: %v", err)
        c.JSON(http.StatusInternalServerError, model.ErrorResponse{Message: "Failed to update pathway"})
        return
    }
    defer res.Body.Close()

    // Log response status code and headers
    log.Printf("Response Status Code: %d", res.StatusCode)
    log.Printf("Response Headers: %v", res.Header)

    responseBody, err := ioutil.ReadAll(res.Body)
    if err != nil {
        log.Printf("Error reading response: %v", err)
        c.JSON(http.StatusInternalServerError, model.ErrorResponse{Message: "Failed to read response"})
        return
    }

    // Log response body
    log.Printf("Response Body: %s", string(responseBody))

    var apiResponse model.UpdatePathwayResponse
    if err := json.Unmarshal(responseBody, &apiResponse); err != nil {
        log.Printf("Error unmarshaling response: %v", err)
        c.JSON(http.StatusInternalServerError, model.ErrorResponse{Message: "Failed to parse response"})
        return
    }

    if apiResponse.Status != "success" {
        log.Printf("API returned error status: %s, message: %s", apiResponse.Status, apiResponse.Message)
        c.JSON(http.StatusInternalServerError, model.ErrorResponse{Message: apiResponse.Message})
        return
    }

    log.Printf("Pathway updated successfully.")
    c.JSON(http.StatusOK, apiResponse.PathwayData)
}

// DeletePathway godoc
// @Summary      Delete a conversational pathway
// @Description  Deletes a specific conversational pathway by its ID
// @Tags         Pathway
// @Accept       json
// @Produce      json
// @Param        pathway_id path string true "Pathway ID to delete"
// @Success      200  {object}  model.DeletePathwayResponse  "Pathway deleted successfully"
// @Failure      400  {object}  model.ErrorResponse  "Invalid input"
// @Failure      500  {object}  model.ErrorResponse  "Internal server error"
// @Security     bearerToken
// @Router       /delete/convo_pathway/{pathway_id} [delete]
func DeletePathway(c *gin.Context) {
    // Step 1: Extract pathway_id from the URL path
    pathwayID := c.Param("pathway_id")
    if pathwayID == "" {
        log.Printf("Error: Pathway ID is missing")
        c.JSON(http.StatusBadRequest, gin.H{"error": "Pathway ID is required"})
        return
    }

    // Step 2: Extract the authorization bearer token
    bearerToken := c.GetHeader("Authorization")
    if bearerToken == "" {
        log.Printf("Authorization token missing")
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token is required"})
        return
    }

    // Step 3: Prepare the external API request to delete the pathway
    apiURL := fmt.Sprintf("https://api.bland.ai/v1/convo_pathway/%s", pathwayID)
    req, err := http.NewRequest("DELETE", apiURL, nil)
    if err != nil {
        log.Printf("Error creating external API request: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create external API request"})
        return
    }

    // Add necessary headers for the external API call
    req.Header.Add("Authorization", bearerToken)

    // Step 4: Send the request to the external API
    client := &http.Client{}
    res, err := client.Do(req)
    if err != nil {
        log.Printf("Error sending external API request: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete pathway"})
        return
    }
    defer res.Body.Close()

    // Step 5: Read and parse the response from the external API
    responseBody, err := ioutil.ReadAll(res.Body)
    if err != nil {
        log.Printf("Error reading external API response: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
        return
    }

    var apiResponse model.DeletePathwayResponse
    if err := json.Unmarshal(responseBody, &apiResponse); err != nil {
        log.Printf("Error unmarshaling external API response: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse response"})
        return
    }

    // Step 6: Log and return the successful response to the client
    log.Printf("Pathway deleted successfully. Pathway ID: %s", apiResponse.PathwayID)
    c.JSON(http.StatusOK, apiResponse)
}



// SendMessageToChat godoc
// @Summary      Send a message to a pathway chat
// @Description  Sends a message to a specific pathway chat and receives a response
// @Tags         Chat
// @Accept       json
// @Produce      json
// @Param        chat_id path string true "Chat ID to send message to"
// @Param        request body model.SendMessageRequest true "Request body for sending a message"
// @Success      200  {object}  model.SendMessageResponse  "Message sent successfully"
// @Failure      400  {object}  model.ErrorResponse        "Invalid input"
// @Failure      401  {object}  model.ErrorResponse        "Unauthorized - Bearer token required"
// @Failure      500  {object}  model.ErrorResponse        "Internal server error"
// @Security     bearerToken
// @Router       /pathways/chat/{chat_id}/send [post]
func SendMessageToChat(c *gin.Context) {
	// Step 1: Extract chat_id from the URL path
	chatID := c.Param("chat_id")
	if chatID == "" {
		log.Printf("Error: Chat ID is missing")
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Message: "Chat ID is required"})
		return
	}

	// Step 2: Bind the request body to SendMessageRequest struct
	var messageRequest model.SendMessageRequest
	if err := c.ShouldBindJSON(&messageRequest); err != nil {
		log.Printf("Error binding JSON for SendMessageRequest: %v", err)
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Message: "Invalid request body"})
		return
	}

	// Step 3: Extract the authorization bearer token
	bearerToken := c.GetHeader("Authorization")
	if bearerToken == "" {
		log.Printf("Authorization token missing")
		c.JSON(http.StatusUnauthorized, model.ErrorResponse{Message: "Authorization token is required"})
		return
	}

	// Step 4: Marshal the messageRequest struct into JSON format
	requestBodyJSON, err := json.Marshal(messageRequest)
	if err != nil {
		log.Printf("Error marshaling request body: %v", err)
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Message: "Failed to process request"})
		return
	}

	// Step 5: Prepare the external API request to send the message
	apiURL := fmt.Sprintf("https://api.bland.ai/v1/pathway/chat/%s", chatID)
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(requestBodyJSON))
	if err != nil {
		log.Printf("Error creating external API request: %v", err)
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Message: "Failed to create external API request"})
		return
	}

	// Add necessary headers for the external API call
	req.Header.Add("Authorization", bearerToken)
	req.Header.Add("Content-Type", "application/json")

	// Step 6: Send the request to the external API
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending external API request: %v", err)
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Message: "Failed to send message"})
		return
	}
	defer res.Body.Close()

	// Step 7: Read and parse the response from the external API
	responseBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("Error reading external API response: %v", err)
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Message: "Failed to read response"})
		return
	}

	// Log the response for debugging
	log.Printf("External API Response: %s", string(responseBody))

	// Step 8: Unmarshal the response into SendMessageResponse struct
	var apiResponse model.SendMessageResponse
	if err := json.Unmarshal(responseBody, &apiResponse); err != nil {
		log.Printf("Error unmarshaling external API response: %v", err)
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Message: "Failed to parse response"})
		return
	}

	// Step 9: Check for errors in the API response
	if apiResponse.Errors != nil {
		log.Printf("API responded with errors: %v", *apiResponse.Errors)
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Message: "External API returned errors"})
		return
	}

	// Step 10: Return the successful response to the client
	c.JSON(http.StatusOK, apiResponse)
}