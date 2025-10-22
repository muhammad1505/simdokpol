# Sistem Informasi Manajemen Dokumen Kepolisian (SIMDOKPOL)

![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)
![Go Version](https://img.shields.io/badge/Go-1.22%2B-blue.svg)

**SIMDOKPOL** adalah aplikasi desktop _standalone_ yang dirancang untuk membantu unit kepolisian dalam manajemen dan penerbitan surat keterangan secara efisien, cepat, dan aman. Dibangun untuk berjalan 100% offline, aplikasi ini menyediakan alur kerja modern dalam format yang mudah didistribusikan dan diinstal.

---

## ✨ Fitur Utama

-   **Aplikasi Desktop Standalone**: Berjalan sebagai aplikasi mandiri dengan ikon di _system tray_ dan notifikasi _native_, memberikan pengalaman pengguna yang terintegrasi.
-   **Installer Profesional**: Didistribusikan sebagai file installer (`.exe`) untuk Windows dan AppImage untuk Linux, membuat proses instalasi menjadi mudah bagi pengguna akhir.
-   **Alur Setup Awal Terpandu**: Konfigurasi pertama kali yang mudah untuk mengatur detail instansi (KOP surat, nama kantor) dan membuat akun Super Admin.
-   **Manajemen Dokumen Lengkap (CRUD)**: Sistem penuh untuk Membuat, Membaca, Memperbarui, dan Menghapus surat keterangan, termasuk fitur **Buat Ulang (Duplikat)** untuk efisiensi.
-   **Manajemen Pengguna Berbasis Peran**: Dua tingkat hak akses (Super Admin & Operator) dengan fitur untuk menonaktifkan dan mengaktifkan kembali akun pengguna.
-   **Dasbor Analitik Real-Time**: Tampilan ringkasan data dengan kartu statistik dan grafik interaktif untuk memonitor aktivitas operasional.
-   **Formulir Cerdas & Dinamis**: Input tanggal yang konsisten, data barang hilang yang interaktif, dan sistem rekomendasi petugas otomatis berdasarkan regu.
-   **Fitur Backup & Restore**: Super Admin dapat dengan mudah mencadangkan dan memulihkan seluruh database aplikasi.
-   **Modul Audit Log Komprehensif**: Setiap aksi penting (pembuatan/pembaruan/penghapusan data) dicatat secara otomatis untuk akuntabilitas.
-   **Pratinjau Cetak Presisi Tinggi**: Halaman pratinjau cetak yang dirancang agar 100% cocok dengan format fisik surat resmi.
-   **100% Offline**: Semua aset (font, CSS, JavaScript) dan fungsionalitas dirancang untuk berjalan tanpa koneksi internet.

---

## 🛠️ Tumpukan Teknologi (Technology Stack)

-   **Backend**: Go (Golang)
-   **Web Framework**: Gin
-   **ORM & Migrasi**: GORM & golang-migrate
-   **Database**: SQLite
-   **Desktop UI**: Go HTML Templates, CSS, JavaScript (jQuery, Bootstrap, Chart.js)
-   **Integrasi Desktop**:
    -   System Tray: `getlantern/systray`
    -   Notifikasi: `gen2brain/beeep`
-   **Build & Packaging**:
    -   Installer Windows: NSIS
    -   Paket Linux: AppImageKit

---

## 🚀 Memulai (Getting Started)

### Untuk Pengguna Akhir

1.  Unduh versi terbaru dari halaman **[Releases](https://github.com/USERNAME/REPO/releases)**.
2.  **Di Windows**:
    -   Jalankan file `SIMDOKPOL-Setup-vX.X.X.exe`.
    -   Ikuti wizard instalasi.
    -   Jalankan aplikasi dari shortcut di Desktop atau Start Menu.
3.  **Di Linux**:
    -   Unduh file `.AppImage`.
    -   Buat file tersebut dapat dieksekusi: `chmod +x SIMDOKPOL-x86_64.AppImage`.
    -   Jalankan dengan klik dua kali atau dari terminal: `./SIMDOKPOL-x86_64.AppImage`.

### Untuk Pengembang

#### Prasyarat

-   Go (versi 1.22+)
-   Git
-   C Compiler (TDM-GCC/Mingw-w64 untuk Windows, `build-essential` untuk Debian/Ubuntu, `base-devel` untuk Manjaro/Arch)
-   Untuk _Cross-Compile_ ke Windows dari Linux: `mingw-w64-gcc`
-   Untuk membuat installer: NSIS (bisa diinstal via Wine atau native di Linux)

#### Instalasi & Menjalankan

1.  **Kloning Repositori**:
    ```bash
    git clone [URL_REPOSITORI_ANDA] simdokpol
    cd simdokpol
    ```
2.  **Konfigurasi Environment**:
    ```bash
    cp .env.example .env
    ```
3.  **Instalasi Dependensi**:
    ```bash
    go mod tidy
    ```
4.  **Menjalankan di Mode Pengembangan (dengan Live Reload)**:
    Untuk pengembangan, disarankan menggunakan [Air](https://github.com/cosmtrek/air).
    ```bash
    air
    ```
    Aplikasi akan berjalan di `http://localhost:8080`.

#### Membangun Aplikasi (Building)

-   **Build untuk Linux**:
    ```bash
    go build -ldflags="-s -w" -o simdokpol ./cmd/main.go
    ```
-   **Cross-Compile untuk Windows (dari Linux)**:
    ```bash
    CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc go build -ldflags="-s -w -H windowsgui -extldflags \"-static\"" -o simdokpol.exe ./cmd/main.go
    ```

---

## 📂 Struktur Proyek

simdokpol/ ├── cmd/main.go # Entrypoint aplikasi (logika systray & server) ├── dist/ # Folder untuk menampung file siap distribusi ├── internal/ # Logika inti aplikasi (controllers, services, etc.) ├── migrations/ # File migrasi skema database ├── web/ # Aset frontend (HTML, CSS, JS, images) ├── .air.toml # Konfigurasi untuk live-reload Air ├── .env.example # Contoh file environment ├── go.mod # Manajemen dependensi Go ├── setup.nsi # Skrip installer NSIS untuk Windows └── vhost-manager.sh/bat # Skrip helper untuk konfigurasi virtual host

---

## 🛣️ Rencana Pengembangan (Roadmap)

-   [x] **Arsitektur Backend & Frontend yang Bersih**
-   [x] **Alur Kerja Surat Keterangan Hilang (CRUD Lengkap)**
-   [x] **Otentikasi & Otorisasi Berbasis Peran**
-   [x] **Manajemen Pengguna (CRUD, Aktivasi/Deaktivasi)**
-   [x] **Dasbor Analitik dengan Grafik Dinamis**
-   [x] **Alur Setup Awal Terpandu**
-   [x] **Modul Audit Log & Fitur Backup/Restore**
-   [x] **Aplikasi Desktop Standalone (via Systray & Beeep)**
-   [x] **Installer Profesional (NSIS untuk Windows, AppImage untuk Linux)**
-   [x] **Penambahan Unit Test & Integration Test (Blueprint sudah ada)**
-   [x] **Dokumentasi API yang lebih mendalam**
-   [x] **Peningkatan Fitur Pencetakan (kustomisasi template, dll)**

---

## 📄 Lisensi

Proyek ini dilisensikan di bawah **MIT License**.

![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)
![Go Version](https://img.shields.io/badge/Go-1.22%2B-blue.svg)
