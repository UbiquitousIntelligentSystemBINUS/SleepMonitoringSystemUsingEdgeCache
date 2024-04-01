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

# Panduan Instalasi Kubernetes Cluster menggunakan Terraform

Berikut adalah langkah-langkah untuk menggunakan Terraform untuk menerapkan kluster Kubernetes di lingkungan Anda.

## Persyaratan Sebelum Instalasi:

1. Pastikan Anda memiliki akun layanan cloud (seperti AWS, GCP, atau Azure) dan telah mengkonfigurasi kredensial akses yang sesuai.
2. Instalasi Terraform di komputer lokal Anda. Anda dapat mengunduh Terraform dari situs web resmi mereka. Direkomendasikan menggunakan Google Cloud Console dimana sudah terdapat instalasi terraform nya.

## Langkah-langkah:

```bash
$ terraform init
$ terraform plan
$ terraform apply
```

# Panduan Testing

Berikut adalah langkah-langkah untuk melakukan testing menggunakan JMeter:
1. Tambahkan plugin MQTT yang terdapat di folder infrastructure pada JMeter.
2. Lalu open file dengan tipe .jmx pada folder infrastructure tersebut.
