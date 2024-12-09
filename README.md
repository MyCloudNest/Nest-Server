# Documentation for MyCloudNest Server

## Overview
MyCloudNest Server is a file management server built using the Fiber framework in Go. It provides features like file upload, download, temporary link generation, and file statistics with robust security and performance configurations.

---

## Features
1. **File Management**:
   - Upload, retrieve, and delete files.
   - Organize files into specific directories based on file type.

2. **Temporary Links**:
   - Generate temporary download links for files with expiration.

3. **Statistics**:
   - Track file downloads and last access timestamps.

4. **Rate Limiting**:
   - Control the number of requests using Redis as a storage backend.

5. **Caching**:
   - Cache frequently accessed data with support for memory-based storage.

6. **Access Control**:
   - Restrict access to specific IP addresses using a whitelist middleware.

---

## Configuration
The server's configuration is loaded from `~/.cloudnest/config.toml`. Below are the configurable sections:

### Server Configuration
```toml
[server]
host = "0.0.0.0"
port = 3000
```

### Rate Limiting
```toml
[rate_limit]
enabled = true
limit_body = 104857600 # 100MB
max_requests = 100
expire_time = 60
```

### Whitelist
```toml
[whitelist]
enabled = true
whitelisted_ips = ["127.0.0.1"]
```

### Performance
```toml
[perfomance]
perfork = false
concurrency = 65536
```

### Cache
```toml
[cache]
enabled = true
```

---

## API Endpoints
### **1. File Upload**
- **Endpoint**: `/api/v1/files`
- **Method**: `POST`
- **Description**: Uploads a file to the server and organizes it based on its type.

### **2. Retrieve Files**
- **Endpoint**: `/api/v1/files`
- **Method**: `GET`
- **Description**: Retrieves a list of all files with their metadata.

### **3. Get File**
- **Endpoint**: `/api/v1/files/:id`
- **Method**: `GET`
- **Description**: Downloads or displays a file.

### **4. Delete File**
- **Endpoint**: `/api/v1/files/:id`
- **Method**: `DELETE`
- **Description**: Deletes a file from the server and its metadata.

### **5. Generate Temporary Link**
- **Endpoint**: `/api/v1/files/:id/temp-link`
- **Method**: `POST`
- **Description**: Generates a temporary download link for a file.

### **6. Validate Temporary Link**
- **Endpoint**: `/api/v1/files/download`
- **Method**: `GET`
- **Description**: Validates and serves a file through a temporary link.

### **7. File Statistics**
- **Endpoint**: `/api/v1/files/:id/stats`
- **Method**: `GET`
- **Description**: Retrieves download statistics for a specific file.

---

## Middleware
### **Whitelist Middleware**
- Restricts access to the application based on IP addresses.

---

## Database
The server uses SQLite for persistent storage. The database is located at `~/.cloudnest/db.sqlite`.

### Tables
- **`files`**: Stores metadata about files.
- **`temporary_links`**: Stores information about temporary links.
- **`file_stats`**: Tracks file download statistics.

---

## How to Run
1. Ensure the configuration file exists at `~/.cloudnest/config.toml`.
2. Start the server:
   ```bash
   go run main.go
   ```
3. Access the API at `http://<host>:<port>`.

---

## Contributors
- [Rio](https://github.com/pageton)

## License
This project is licensed under the MIT License.
