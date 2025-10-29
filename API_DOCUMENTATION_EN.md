# Electric Circuit Web API Documentation (English)

## Overview
Electric Circuit Web API is a RESTful service built with Go and Firebase for electronic circuit design applications. The API follows Clean Architecture principles with Controller → Handler → Service → Repository → pkg layers.

## Base URL
- Development: `http://localhost:8080/api`
- Production: `https://api.electric-circuit-web.com/api`

## Authentication

This API uses Firebase Authentication for user management.

### Recommended Authentication Flow (Firebase SDK Method)

The recommended and common approach is client-side authentication using Firebase SDK:

1. **Login directly with Firebase SDK on client**
   ```javascript
   import { signInWithEmailAndPassword, getAuth } from 'firebase/auth';
   
   const auth = getAuth();
   const userCredential = await signInWithEmailAndPassword(auth, email, password);
   const idToken = await userCredential.user.getIdToken();
   ```

2. **Send ID token to server for verification**
   ```javascript
   const response = await fetch('/api/auth/verify', {
     method: 'POST',
     headers: {
       'Content-Type': 'application/json',
     },
     body: JSON.stringify({ token: idToken })
   });
   ```

3. **Include token in subsequent API requests**
   ```
   Authorization: Bearer <firebase_id_token>
   ```

### APIs requiring authentication
Most API endpoints require Firebase ID Token authentication (except `/health`, `/auth/verify`, `/auth/create-user`)

### Token refresh
Firebase ID tokens expire after 1 hour, so the client should automatically refresh tokens:
```javascript
// Firebase automatically refreshes tokens
const freshToken = await user.getIdToken(true);
```

## API Endpoints

### System Endpoints

#### Health Check
```http
GET /health
```
Check server status.

**Response:**
```json
{
  "status": "healthy"
}
```

### Authentication Endpoints

#### Verify Token (Recommended Method)
```http
POST /auth/verify
```
**Verify Firebase ID token obtained from client-side Firebase SDK authentication.**

**Use case**: After Firebase SDK login, send ID token to server for verification

**Usage scenario:**
1. Client logs in with Firebase SDK
2. Send obtained ID token to this endpoint
3. Server verifies token validity and returns user information

**Request Body:**
```json
{
  "token": "firebase_id_token"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Token verified successfully",
  "user": {
    "uid": "firebase_user_id",
    "email": "user@example.com",
    "displayName": "John Doe",
    "photoURL": "https://example.com/photo.jpg",
    "emailVerified": true
  }
}
```

#### Create User (Server-side Creation)
```http
POST /auth/create-user
```
**Create user account directly on server (for admin use or special cases)**

**Note**: Generally, it's recommended to create accounts using Firebase SDK on the client side

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "password123",
  "display_name": "John Doe",
  "photo_url": "https://example.com/photo.jpg"
}
```

#### Get User
```http
GET /auth/get-user?uid=user_id
```
Retrieve user information.

#### Update User
```http
PUT /auth/update-user
```
Update user information.

**Request Body:**
```json
{
  "display_name": "Updated Name",
  "photo_url": "https://example.com/new-photo.jpg"
}
```

#### Delete User
```http
DELETE /auth/delete-user?uid=user_id
```
Delete user account.

### Project Endpoints

#### List Projects
```http
GET /projects
```
Get all projects for the authenticated user.

**Response:**
```json
{
  "success": true,
  "message": "Projects retrieved successfully",
  "projects": [
    {
      "id": "project_123",
      "name": "LED Circuit Design",
      "description": "Basic LED circuit project",
      "user_id": "user_123",
      "status": "active",
      "settings": {
        "grid_size": 10,
        "snap_to_grid": true
      },
      "created_at": "2024-01-15T10:30:00Z",
      "updated_at": "2024-01-16T14:20:00Z"
    }
  ]
}
```

#### Create Project
```http
POST /projects/create
```
Create a new project.

**Request Body:**
```json
{
  "name": "LED Circuit Design",
  "description": "Basic LED circuit project"
}
```

#### Get Project
```http
GET /projects/get?projectId=project_id
```
Get specific project details.

#### Update Project
```http
PUT /projects/update
```
Update project information.

**Request Body:**
```json
{
  "project_id": "project_123",
  "name": "Updated Project Name",
  "description": "Updated description"
}
```

#### Delete Project
```http
DELETE /projects/delete?projectId=project_id
```
Delete a project.

#### Duplicate Project
```http
POST /projects/duplicate
```
Create a copy of an existing project.

**Request Body:**
```json
{
  "project_id": "original_project_id",
  "name": "Duplicated Project"
}
```

### Circuit Endpoints

#### List Circuits
```http
GET /circuits?projectId=project_id
```
Get all circuits for a specific project.

**Response:**
```json
{
  "success": true,
  "message": "Circuits retrieved successfully",
  "circuits": [
    {
      "id": "circuit_123",
      "name": "Main Circuit",
      "description": "Primary circuit of the project",
      "project_id": "project_123",
      "user_id": "user_123",
      "data": {
        "elements": [
          {
            "id": "R1",
            "type": "resistor",
            "value": "1k",
            "position": { "x": 100, "y": 100 }
          }
        ],
        "connections": [
          {
            "from": "R1.pin2",
            "to": "LED1.pin1"
          }
        ]
      },
      "version": 1,
      "is_template": false,
      "tags": ["led", "resistor"],
      "created_at": "2024-01-15T11:00:00Z",
      "updated_at": "2024-01-16T15:30:00Z"
    }
  ]
}
```

#### Create Circuit
```http
POST /circuits/create
```
Create a new circuit.

**Request Body:**
```json
{
  "project_id": "project_123",
  "name": "Main Circuit",
  "description": "Primary circuit description",
  "data": {
    "elements": [...],
    "connections": [...]
  }
}
```

#### Get Circuit
```http
GET /circuits/get?circuitId=circuit_id
```
Get specific circuit details.

#### Update Circuit
```http
PUT /circuits/update
```
Update circuit design and information.

**Request Body:**
```json
{
  "circuit_id": "circuit_123",
  "name": "Updated Circuit",
  "description": "Updated description",
  "data": {
    "elements": [...],
    "connections": [...]
  }
}
```

#### Delete Circuit
```http
DELETE /circuits/delete?circuitId=circuit_id
```
Delete a circuit.

#### Get Templates
```http
GET /circuits/templates
```
Get available circuit templates (not yet implemented).

### Storage Endpoints

#### Upload File
```http
POST /storage/upload
Content-Type: multipart/form-data
```
Upload files to cloud storage.

**Form Data:**
- `file`: Binary file data
- `folder`: Target folder name (optional)

**Response:**
```json
{
  "success": true,
  "message": "File uploaded successfully",
  "download_url": "https://storage.example.com/files/file_id",
  "file_path": "/uploads/user_id/filename.ext",
  "file_name": "filename.ext",
  "size": 1024
}
```

#### Get File URL
```http
GET /storage/url?filePath=file_path
```
Get download URL for stored file.

#### Delete File
```http
DELETE /storage/delete?filePath=file_path
```
Delete stored file.

## Error Responses

All error responses follow this format:
```json
{
  "success": false,
  "error": "Error message",
  "message": "Additional details"
}
```

### Common HTTP Status Codes
- `200`: Success
- `201`: Created
- `400`: Bad Request
- `401`: Unauthorized
- `404`: Not Found
- `500`: Internal Server Error

## Data Models

### User
```json
{
  "uid": "string",
  "email": "string",
  "displayName": "string",
  "photoURL": "string",
  "emailVerified": "boolean"
}
```

### Project
```json
{
  "id": "string",
  "name": "string",
  "description": "string",
  "user_id": "string",
  "status": "active|archived|deleted",
  "settings": "object",
  "created_at": "datetime",
  "updated_at": "datetime"
}
```

### Circuit
```json
{
  "id": "string",
  "name": "string",
  "description": "string",
  "project_id": "string",
  "user_id": "string",
  "data": {
    "elements": "array",
    "connections": "array"
  },
  "version": "integer",
  "is_template": "boolean",
  "tags": "array",
  "created_at": "datetime",
  "updated_at": "datetime"
}
```

## Rate Limiting
API requests are subject to rate limiting. Contact support if you need higher limits.

## Support
For API support, contact: support@example.com