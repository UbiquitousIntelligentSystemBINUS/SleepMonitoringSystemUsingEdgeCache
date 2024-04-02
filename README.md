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

Berikut adalah langkah-langkah untuk melakukan testing:
1. Build image dari servis pada folder gateway-application dan quantification
2. Masukkan image ke dalam docker
3. Buat container dari image-image tersebut
4. Pull image-image lain yang diperlukan (redis, database, mqtt)
5. Jalankan semua container yang diperlukan
6. Buka aplikasi JMeter
7. Set up testing sesuai yang diperlukan
8. Lakukan testing
