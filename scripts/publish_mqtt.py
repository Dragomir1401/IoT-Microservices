import paho.mqtt.client as mqtt
import time
import json
import random

BROKER_HOST = "localhost"
BROKER_PORT = 1883
TOPICS = ["UPB/RPi_1", "UPB/RPi_2", "Dorinel/Zeus"]

# Function to generate random payload
def generate_payload():
    return {
        "BAT": random.randint(0, 100),
        "HUMID": round(random.uniform(30.0, 70.0), 1),
        "TMP": round(random.uniform(20.0, 35.0), 1),
        "timestamp": time.strftime("%Y-%m-%dT%H:%M:%S%z"),
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

            time.sleep(5)

    except KeyboardInterrupt:
        print("\nStopped publishing messages.")
    except Exception as e:
        print(f"Error: {e}")
    finally:
        client.disconnect()

if __name__ == "__main__":
    main()
