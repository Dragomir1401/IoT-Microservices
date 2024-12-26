import paho.mqtt.client as mqtt
import time
import json
import random
from datetime import datetime, timezone, timedelta

BROKER_HOST = "localhost"
BROKER_PORT = 1883
TOPICS = ["UPB/RPi_1", "UPB/RPi_2", "Dorinel/Zeus"]

# Function to generate random payload
def generate_payload():
    # Generate a timestamp with RFC3339 format including timezone offset
    current_time = datetime.now(timezone(timedelta(hours=2)))
    timestamp = current_time.isoformat()

    return {
        "BAT": random.randint(80, 100),
        "HUMID": round(random.uniform(50.0, 70.0), 1),
        "TMP": round(random.uniform(20.0, 25.0), 1),
        "timestamp": timestamp,
    }

def main():
    client = mqtt.Client()

    try:
        client.connect(BROKER_HOST, BROKER_PORT, 60)
        print(f"Connected to MQTT broker at {BROKER_HOST}:{BROKER_PORT}")

        while True:
            for topic in TOPICS:
                payload = generate_payload()
                message = json.dumps(payload)
                client.publish(topic, message)
                print(f"Published to {topic}: {message}")

            time.sleep(0.5)

    except KeyboardInterrupt:
        print("\nStopped publishing messages.")
    except Exception as e:
        print(f"Error: {e}")
    finally:
        client.disconnect()

if __name__ == "__main__":
    main()
