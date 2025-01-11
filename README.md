# MQTT Adaptor Microservices IoT App

This project implements an MQTT-based IoT adaptor with integrated time-series data storage using InfluxDB 1.8
(no authentication) and dashboard visualization using Grafana. The solution processes MQTT messages,
writes the data to InfluxDB, and visualizes the data through pre-configured dashboards.

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
- **Image:** `influxdb:1.8`
- **Purpose:** Stores time-series data from IoT devices, with measurements, tags, and fields corresponding to sensor data.
- **Configuration:**
    - Initialized with an admin user, organization, and bucket.
    - Environment variables are used for setup and authentication:
        - Authentication is disabled using the `INFLUXDB_HTTP_AUTH_ENABLED=false` environment variable.
        - Flux query support is enabled with `INFLUXDB_FLUX_ENABLED=true`
    - Data is persisted to the volume `influxdb1.8-data`.
- **Networks:**
    - `database_net` for communication with the MQTT adaptor.
    - `grafana_net` for integration with Grafana.
- **Deployment:**
    - Ensures high availability with `on-failure` restart policy.
- **Auth** - No authentication is enabled for InfluxDB.

### **3. Grafana**
- **Image:** `grafana/grafana`
- **Purpose:** Provides visualization dashboards for IoT measured data stored in InfluxDB.
- **Configuration:**
    - Pre-configured datasources and dashboards are provisioned via mounted directories.
    - Environment variables configure the admin credentials:
        - `GF_SECURITY_ADMIN_USER=asistent`
        - `GF_SECURITY_ADMIN_PASSWORD=grafanaSCD2024`
    - Data is persisted to the volume `grafana-data`.
- **Networks:**
    - `grafana_net` for integration with InfluxDB.

### **4. MQTT Adaptor**
- **Image:** Custom-built `mqtt-adaptor` image.
- **Purpose:** Serves as the bridge between IoT devices and InfluxDB, processing MQTT messages, validating payloads, and writing data to the database.
- **Configuration:**
    - Environment variables:
        - `MQTT_BROKER`: Connects to the Mosquitto broker.
        - `INFLUXDB_URL`, `INFLUXDB_DATABASE`: Configures InfluxDB connection and database.
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
    - Ensures integration between Grafana and InfluxDB.

---

### **Volumes**
- **`influxdb-data`:**
    - Persists InfluxDB time-series data and configuration.
- **`grafana-data`:**
    - Stores Grafana dashboards and configurations for persistence when doing restarts.
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
This structure ensures a consistent hierarchy.
- **Payload Format:**
    - The payload is a JSON dictionary with one level of nesting.
    - Keys represent sensor data fields (e.g., `BAT`, `HUMID`, `TMP`) and must be numeric.
    - A `timestamp` key indicates the measurement time; if missing or invalid, the current time is used.
    - Non-numeric values and invalid messages are ignored.

### **2. Persistence**
- InfluxDB and Grafana are used for telemetry data storage and visualization, with their data persisted to dedicated volumes.

### **3. Modular Design**
- The system is containerized with the following services:
    - **MQTT Broker:** Mosquitto broker for IoT device communication.
    - **MQTT Adaptor:** Connects IoT devices to InfluxDB and synchronizes broker configuration.
    - **Database:** Stores data for IoT sensors measurements.
    - **Grafana:** Provides visualization dashboards for the measurements.

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

