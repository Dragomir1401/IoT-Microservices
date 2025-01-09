# MQTT Adaptor Microservices IoT App

This project implements an MQTT-based IoT adaptor with integrated authentication and dashboard visualization using Grafana.
The solution manages MQTT topics, authenticates IoT devices, and authorizes their access to publish/subscribe
to specific topics. It uses docker to manage the services for the adapter, mqtt broker, influxdb database and grafana dashboard.

---

## **Stack Overview**

The system is designed as a modular stack using Docker Compose, with each service fulfilling a specific
role in the MQTT-based IoT platform. Below is an overview of the stack configuration:

### **1. Mosquitto (MQTT Broker)**
- **Image:** `eclipse-mosquitto`
- **Purpose:** Acts as the MQTT message broker to handle IoT device communication.
- **Configuration:**
    - Port `1883` is exposed for MQTT communication.
    - A custom configuration file (`mosquitto.conf`) is mounted to allow for authentication and authorization setups.
- **Networks:**
    - `mqtt_net` for communication with MQTT clients and the adaptor.

### **2. InfluxDB**
- **Image:** `influxdb:2.0`
- **Purpose:** Stores time-series data from IoT devices, with measurements, tags, and fields corresponding to sensor data.
- **Configuration:**
    - Initialized with an admin user, organization, and bucket.
    - Environment variables are used for setup and authentication:
        - `DOCKER_INFLUXDB_INIT_*` for initial setup.
        - `DOCKER_INFLUXDB_HTTP_AUTH_ENABLED=true` for enabling HTTP authentication.
    - Data is persisted to the volume `influxdb-data`.
- **Networks:**
    - `database_net` for communication with the MQTT adaptor.
    - `grafana_net` for integration with Grafana.
- **Deployment:**
    - Ensures high availability with `on-failure` restart policy.

### **3. Grafana**
- **Image:** `grafana/grafana`
- **Purpose:** Provides visualization dashboards for IoT telemetry data stored in InfluxDB.
- **Configuration:**
    - Pre-configured datasources and dashboards are provisioned via mounted directories.
    - Environment variables configure the admin credentials:
        - `GF_SECURITY_ADMIN_USER=asistent`
        - `GF_SECURITY_ADMIN_PASSWORD=grafanaSCD2024`
    - InfluxDB token is passed via the environment for secure integration.
    - Data is persisted to the volume `grafana-data`.
- **Networks:**
    - `grafana_net` for integration with InfluxDB.

### **4. MQTT Adaptor**
- **Image:** Custom-built `mqtt-adaptor` image.
- **Purpose:** Serves as the bridge between IoT devices and InfluxDB, processing MQTT messages, validating payloads, and writing data to the database.
- **Configuration:**
    - Environment variables:
        - `MQTT_BROKER`: Connects to the Mosquitto broker.
        - `INFLUXDB_URL`, `INFLUXDB_ORG`, `INFLUXDB_BUCKET`, and `INFLUXDB_TOKEN`: Configures InfluxDB connection and authentication.
        - `DEBUG_DATA_FLOW`: Enables debug logging.
- **Networks:**
    - `mqtt_net` for MQTT communication.
    - `database_net` for InfluxDB integration.

---

### **Networks**
- **`mqtt_net`:**
    - Handles communication between the MQTT broker and clients, as well as the MQTT adaptor.
- **`database_net`:**
    - Facilitates secure data transfer between the MQTT adaptor and InfluxDB.
- **`grafana_net`:**
    - Ensures seamless integration between Grafana and InfluxDB.

---

### **Volumes**
- **`influxdb-data`:**
    - Persists InfluxDB time-series data and configuration.
- **`grafana-data`:**
    - Stores Grafana dashboards and configurations for persistence across restarts.
- **Custom Mounts:**
    - `mosquitto.conf`: Custom Mosquitto configuration.
    - Grafana provisioning directories for datasources and dashboards.

---

### **Deployment**
1. Ensure Docker and Docker Compose are installed.
2. Start the stack with:
   ```bash
   ./scripts/run.sh
    ```

## **Implementation Decisions**

### **1. MQTT Topic and Payload Structure**
- **Topic Format:** MQTT topics follow the format `<location>/<station>` (e.g., `UPB/RPi_1`). 
This structure ensures clarity and a consistent hierarchy.
- **Payload Format:**
    - The payload is a JSON dictionary with one level of nesting.
    - Keys represent sensor data fields (e.g., `BAT`, `HUMID`, `TMP`) and must be numeric.
    - A `timestamp` key indicates the measurement time; if missing or invalid, the current time is used.
    - Non-numeric values and invalid messages are ignored.

### **2. Authentication and Authorization**
- **Authentication:**
    - Implemented using the Mosquitto MQTT broker with authentication enabled via a `password_file`.
    - Devices must provide valid credentials to connect.

- **Authorization:**
    - Permissions (read/write access) are managed through an `acl_file`.
    - Permissions are dynamically updated based on user settings via a web interface.

### **3. Web Interface for User Management**
- A RESTful API (`auth-api`) is provided to:
    - Create and manage user accounts.
    - Define publish/subscribe permissions for topics.
    - Retrieve and modify existing permissions.

### **4. Persistence**
- Configuration files (`password_file` and `acl_file`) and user permissions are stored persistently using Docker volumes.
- InfluxDB and Grafana are used for telemetry data storage and visualization, with their data persisted to dedicated volumes.

### **5. Modular Design**
- The system is containerized with the following services:
    - **MQTT Broker:** Mosquitto configured with authentication and ACLs.
    - **Auth API:** A web service to manage users and permissions.
    - **MQTT Adaptor:** Connects IoT devices to InfluxDB and synchronizes broker configuration with the Auth API.
    - **Database:** Stores user credentials and permissions.

---

## **How to Launch**

### **Prerequisites**
1. Install Docker and Docker Compose.
2. Ensure port `8080` (Auth API) and `1883` (MQTT Broker) are available.

### **Steps to Launch**

1. Clone the repository:
```bash
git clone <repository_url>
cd <repository_directory>
```
2. Initialize the required directories:
```bash
mkdir -p ./volumes/influxdb-data ./volumes/grafana-data ./volumes/mosquitto-config
```
3. Start the system:
```bash
./scripts/run.sh
```
4. Verify the services:
```bash
Access the Auth API at http://localhost:8080.
Connect an MQTT client to the broker at mqtt://localhost:1883.
```

### **Additional Notes**
#### Persistence:

All configuration files and databases are stored persistently in the ./volumes directory.
Scalability:

The modular design allows for easy scaling of individual components (e.g., multiple MQTT adaptors).

#### Extensibility:

Additional features, such as advanced permission management or encryption for sensitive data,
can be added with minimal changes to the existing architecture.
