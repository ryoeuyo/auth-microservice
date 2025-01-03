# Auth-Service

## Description
A **gRPC-based authentication service** for managing user authentication.

### Features
- **gRPC APIs** for login and registration
- **JWT-based authentication** system
- **Prometheus metrics** for monitoring
- **Database migrations** with goose


The service uses postgres as a main storage.<br>
Configuration is done with files in ./config directory

## Install and Run

1. Clone the repository:
    ```bash
   $ git clone https://github.com/ryoeuyo/auth-microservice.git
   ```
2. Navigate to the project directory:
   ```bash
   $ cd auth-microservice
   ```
3. Optionally, modify the `config/config-*.yml` files as needed.
4. Build container:
   ```bash 
   docker compose build
   ```
5. Start the service:
   ```bash 
   docker-compose up -d
   ```

## etc
You can also run the service using tasks defined in the `Taskfile.yml`. To do this, you need to:
1. Install the <a href="https://taskfile.dev/installation/">task util</a>
2. Read `Taskfile.yml` and run desired tasks

## Conclusion
If you encounter a bug, please create an issue with a detailed description.
