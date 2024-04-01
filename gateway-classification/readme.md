# Panduan Setup MQTT

Berikut adalah langkah-langkah untuk melakukan setup pada MQTT:
1. Buatlah folder mqtt yang didalamnya terdapat 3 folder lain yaitu config, data, dan log.
2. Lalu konfigurasi mqtt sesuai kebutuhan pada folder mqtt -> config -> mosquitto.conf

# Panduan Menjalankan Dockerfile

Berikut adalah langkah-langkah untuk menjalankan Dockerfile dan membangun image Docker:

## Langkah 1: Build Docker Image

```bash
docker build -t myapp .
```

## Langkah 2: Kirim Image ke Registry Cloud
```bash
docker login <registry-url>
docker tag <image-name>:<tag> <registry-url>/<username>/<image-name>:<tag>
```

# Panduan Testing

Berikut adalah langkah-langkah untuk melakukan testing menggunakan JMeter:
1. Tambahkan plugin MQTT yang terdapat di folder infrastructure pada JMeter.
2. Lalu open file dengan tipe .jmx pada folder infrastructure tersebut.
