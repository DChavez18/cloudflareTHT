## Weather Data Retrieval with Redis Cache

This Go program retrieves weather data from the National Weather Service API, stores it in a Redis cache, and then retrieves the data from the cache. It uses Docker to run a Redis server locally.

### Prerequisites

- Go installed on your system.
- Docker Desktop installed and running.

### Installation and Setup

1. **Clone the Repository**: Clone the repository containing the Go code from the provided source.

2. **Build the Docker Container**: Open a terminal and run the following command to build the Redis Docker container:


3. **Run the Go Program**: Navigate to the directory containing the Go file (`weather.go`) and run the following command to execute the program:



### Usage

- The Go program will fetch weather data from the National Weather Service API for the Austin, Texas area and store it in the Redis cache.
- After storing the data, it retrieves the data from the Redis cache and starts an HTTP server locally on localhost:8080 to display the weather data retrieved from the Redis cache in a web browser.

### Additional Notes

- The weather data is fetched from the National Weather Service API using the endpoint `http://api.weather.gov/gridpoints/EWX/31,80/forecast/hourly`, where `EWX` corresponds to the Austin/San Antonio, Texas forecast office.
- The Redis server is accessed locally on port `6379`.
- Ensure that Docker is running before executing the Docker command to start the Redis container.
- The program uses the Go Redis client library to interact with the Redis server.
- You can customize the weather data retrieval by modifying the API endpoint or the grid coordinates in the `fetchWeatherData` function.
